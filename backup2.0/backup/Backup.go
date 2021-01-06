package backup

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"

	//"github.com/pkg/profile"

	"github.com/michaeldcanady/Project01/backup2.0/dispatcher"
	"github.com/michaeldcanady/Project01/backup2.0/file"
	"github.com/michaeldcanady/Project01/backup2.0/struct"
	"github.com/michaeldcanady/Project01/backup2.0/worker"
	"github.com/michaeldcanady/Project01/backup2.0/zip"
	//"github.com/michaeldcanady/Test/test/2/conversion"
)

var (
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func check(e error) {
	if e != nil {
		ErrorLogger.Println(e)
	}
}

func createdst(dst string, ext string) *os.File {
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		head, _ := filepath.Split(dst)
		os.MkdirAll(head, 0700)
	}
	destination, err := os.OpenFile(dst+ext, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	check(err)
	return destination
}

func Backup(users []string, dst, backuptype, name string, conf structure.Config, backup bool) (int64, int64) {

	file1 := createdst(filepath.Join(dst, "logs", "errorLog"), ".txt")

	WarningLogger = log.New(file1, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file1, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	//defer profile.Start().Stop()
	//This is for tracking the time it takes to backup
	defer timeTrack(time.Now(), "main loop")
	var wg sync.WaitGroup
	dd := dispatcher.New(runtime.NumCPU()).Start()
	output := make(chan *file.File)
	barlist := make(map[string]*mpb.Bar)
	bars := true

	go func() {
		defer close(output)
		for _, user := range users {
			User := getUser(user)
			l, _ := filepath.Glob(User.HomeDir + "/**")
			if bars {
				loadBars(wg, l, &barlist)
			}
			gather(User.HomeDir, output, conf)
		}
		//close(output)
	}()
	switch backuptype {
	case "InLine Copy":
		return InLineCopy(backup, dd, barlist, bars, dst, output)
	case "Zip":
		return ZipCopy(backup, barlist, bars, dst, output)
	}
	return 0, 0
}

func loadBars(wg sync.WaitGroup, list []string, barlist *map[string]*mpb.Bar) {
	p := mpb.New(mpb.WithWaitGroup(&wg))
	for _, barname := range list {
		dir, _ := IsDirectory(barname)
		if !dir {
			continue
		}
		total, _ := DirSize(barname)
		bar := p.AddBar(int64(total),
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(barname),
				// decor.DSyncWidth bit enables column width synchronization
				decor.Percentage(decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				// replace ETA decorator with "done" message, OnComplete event
				decor.OnComplete(
					// ETA decorator with ewma age of 60
					decor.EwmaETA(decor.ET_STYLE_GO, 60), "done",
				),
			),
		)
		(*barlist)[barname] = bar
	}
}

// ZipCopy Copies all users and thier files in to a single zip file
func ZipCopy(backup bool, barlist map[string]*mpb.Bar, bars bool, dst string, output chan *file.File) (int64, int64) {

	newZipFile, err := ioutil.TempFile("", "Users*.zip")
	if err != nil {
		panic(err)
	}

	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()
	var i, s int64
	// Add files to zip
ziploop:
	for {
		select {
		case file, ok := <-output:
			if ok {
				//first parameter zip.Writer, full filepath, elements being removed
				if err = zipfile.AddFileToZip(zipWriter, file.PathStr(), file.PathVol()+"\\Users"); err != nil {
					panic(err)
				}
				i++
				s += (*file.File).Size()
				b := barlist[filepath.Join((*file).PathVol(), (*file).PathUserf(), (*file).PathHead())]
				if b == nil {
					continue
				}
				if bars {
					b.IncrInt64(file.RawSize())
				}
			} else {
				fmt.Println("Channel closed")
				break ziploop
			}
		default:
			continue
		}
	}
	_, name := filepath.Split(newZipFile.Name())
	fmt.Println(filepath.Join(dst, name))
	err = os.Rename(newZipFile.Name(), filepath.Join(dst, name))
	if err != nil {
		os.Remove(newZipFile.Name())
	}
	return i, s
}

// InLineCopy copies all files gathered in Gatherer and sends them directly to thier new location
func InLineCopy(backup bool, dd *dispatcher.Dispatcher, barlist map[string]*mpb.Bar, bars bool, dst string, output chan *file.File) (int64, int64) {

	var i, s int64
copyloop:
	for {
		select {
		case file, ok := (<-output):
			if ok {
				dd.Submit(worker.NewJob(i, file, dst, backup))
				i++
				s += (*file.File).Size()
				b := barlist[filepath.Join((*file).PathVol(), (*file).PathUserf(), (*file).PathHead())]
				if b == nil {
					continue
				}
				if bars {
					b.IncrInt64(file.RawSize())
				}
			} else {
				fmt.Println("Channel closed!")
				break copyloop
			}
		default:
			continue
		}
	}
	return i, s
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}
