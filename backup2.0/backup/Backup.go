package backup

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"

	//"github.com/pkg/profile"

	"github.com/michaeldcanady/LookBack/backup2.0/dispatcher"
	"github.com/michaeldcanady/LookBack/backup2.0/file"
	"github.com/michaeldcanady/LookBack/backup2.0/struct"
	//"github.com/michaeldcanady/Test/test/2/conversion"
)

var (
	TraceLogger     *log.Logger
	DebugLogger     *log.Logger
	InfoLogger      *log.Logger
	WarnLogger      *log.Logger
	ErrorLogger     *log.Logger
	FatalLogger     *log.Logger
	HashErrorLogger *log.Logger
)

// Mega checking function for writing errors according to type
func check(err error, errType string) bool {
	if err != nil {
		switch strings.Title(errType) {
		case "Trace":
			TraceLogger.Println(err)
		case "Debug":
			DebugLogger.Println(err)
		case "Info":
			InfoLogger.Println(err)
		case "Warn":
			WarnLogger.Println(err)
		case "Error":
			ErrorLogger.Println(err)
		case "Fatal":
			FatalLogger.Println(err)
		default:
			log.Fatalf("%s is an invalid type", errType)
		}
		return false
	}
	return true
}

//createdst checks if file does not exist
//if it doesn't, the file is split
func createdst(dst string, ext string) *os.File {
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		// splits the head (C:/Users/Username/.../) and tail file.ext
		head, _ := filepath.Split(dst)
		// creates all missing directories up to but not including the file
		os.MkdirAll(head, 0700)
	}
	//opens the file, with ability to append and create
	destination, err := os.OpenFile(dst+ext, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	//checks for errors
	check(err, "error")
	// return pointer of os.File (https://golang.org/pkg/os/#File)
	return destination
}

func createAllLogFiles(traceLog, debugLog, infoLog, warnLog, hErrorLog, errorLog, fatalLog *os.File) {
	TraceLogger = log.New(traceLog, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(debugLog, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(infoLog, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(warnLog, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	HashErrorLogger = log.New(hErrorLog, "HASH: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorLog, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	FatalLogger = log.New(fatalLog, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Backup(users []structure.User, dst, backuptype, name string, conf structure.Config, backup bool) (int64, int64) {

	traceLog := createdst(filepath.Join(dst, "logs", "trace"), ".log")
	debugLog := createdst(filepath.Join(dst, "logs", "debug"), ".log")
	infoLog := createdst(filepath.Join(dst, "logs", "info"), ".log")
	warnLog := createdst(filepath.Join(dst, "logs", "warn"), ".log")
	hErrorLog := createdst(filepath.Join(dst, "logs", "hashError"), ".log")
	errorLog := createdst(filepath.Join(dst, "logs", "error"), ".log")
	fatalLog := createdst(filepath.Join(dst, "logs", "fatal"), ".log")

	createAllLogFiles(traceLog, debugLog, infoLog, warnLog, hErrorLog, errorLog, fatalLog)

	//defer profile.Start().Stop()
	//This is for tracking the time it takes to backup
	defer timeTrack(time.Now(), "main loop")
	var wg sync.WaitGroup
	dd := dispatcher.New(runtime.NumCPU()).Start()
	output := make(chan *file.File)
	barlist := make(map[string]*mpb.Bar)
	bars := true

	go func() {
		//Closes output channel ones the goroutine finishes
		defer close(output)
		for _, user := range users {
			//Checks if progress bars are enabled
			if bars {
				//Loads progress bars
				loadBars(wg, &user.RootDirs, &barlist, 0)
			}
			for k, v := range user.Files{
				output <- file.New(k, &v)
			}
		}
	}()
	// switch statement used to decide what method is used to backup
	// Looking to change this to a cleaner method
	switch backuptype {
	case "InLine Copy":
		// returns file count and size
		return InLineCopy(backup, dd, barlist, bars, dst, output)
	case "Zip":
		// returns file count and size
		return ZipCopy(backup, barlist, bars, dst, output)
	}
	fmt.Println(" ")
	// If error arise both values will return 0
	return 0, 0
}

// include wait group, list of files for bars, bar map for incrementing the desired bar
// total will be utilized when replacement function is created
func loadBars(wg sync.WaitGroup, list *map[string]int64, barlist *map[string]*mpb.Bar, total int64) {
	// Creates new multibar struct utlizing waitgroups
	p := mpb.New(mpb.WithWaitGroup(&wg))
	// iterates through dir paths
	for barname, totalsize := range *list {
		// verifies it is a Is a Directory
		// currently if file continue
		// need to add support for files within the root user folder
		if dir, _ := IsDirectory(barname); !dir {
			continue
		}
		// creates bar
		bar := p.AddBar(int64(totalsize),
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
		// adds name a key and bar pointer for the map
		(*barlist)[barname] = bar
	}
}

// records how long the executed function took
//usage func funcName(){
// defer timeTrack(time.Now(),function name)
//}
// does not work well for recusive function unless you are trying to record each iteration
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// Will be removed once viable solution is developed
func DirSize(path string) (int64, error) {
	// Total size
	var size int64
	// Filepath.Walk to get directory size
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		// add error writer if error arises
		if err != nil {
			return err
		}
		// verifies that info is a file and not a dir
		if !info.IsDir() {
			// adds info's size to the total
			size += info.Size()
		}
		// returns error if error arises
		return err
	})
	// returns size and error (should be nil)
	return size, err
}

// CHecks if the povided path is a directory
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		// needs error writer here
		return false, err
	}
	//returns if file is dire and (ideally) nil
	return fileInfo.IsDir(), err
}
