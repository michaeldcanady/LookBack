//+ build windows darwin linux

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/michaeldcanady/LookBack/backup2.0/backup"
	"github.com/michaeldcanady/LookBack/backup2.0/conversion"
	filestruct "github.com/michaeldcanady/LookBack/backup2.0/path"
	"github.com/michaeldcanady/LookBack/backup2.0/servicenow"
	structure "github.com/michaeldcanady/LookBack/backup2.0/struct"
)

var (
	//SKIPPABLE User accounts that don't need to be included in the options
	SKIPPABLE = []string{"C:\\Users\\Default", "C:\\Users\\Public", "C:\\Users\\All Users", "C:\\Users\\Default User"}
	users     []structure.User
	conf      structure.Config
	MAX       int64
)

func init() {
	SetupCloseHandler()
	test()
	if runtime.GOOS == "windows" {
		MAX = conf.Settings.WinServerBackupMax
	} else {
		MAX = conf.Settings.MacServerBackupMax
	}
	structure.Conf = conf
}

func main() {
	equinoxUpdate()
	// sets backup's information to none
	var binfo structure.Backup
	// Explained in questions.go
	getUserName(&binfo)
	getPassword(&binfo)
	getCSNumber(&binfo)
	getTask(&binfo)
	getSource(&binfo)
	getDestination(&binfo)
	// Explained in questions.go
	for {
		// Checks if user confirmes data or not
		confirm := getConfirmation(&binfo)
		if confirm == false {
			// checks what fields user wants to change
			selected := SelectChange(&binfo)
			for _, s := range selected {
				if s == "Username" {
					getUserName(&binfo)
				} else if s == "Ticket Number" {
					getCSNumber(&binfo)
				} else if s == "Task" {
					getTask(&binfo)
				} else if s == "Source" {
					getSource(&binfo)
				} else if s == "Destination" {
					getDestination(&binfo)
				} else {
					continue
				}
			}
		} else {
			break
		}
	}
	method := getBackupMethod(&binfo, true)
	var volume string
	if !strings.HasSuffix(binfo.Dest, string(os.PathSeparator)) {
		volume = binfo.Dest + string(os.PathSeparator)
	} else {
		volume = binfo.Dest
	}

	Heading(&binfo)
	var i, s int64
	name := getName(binfo.Dest, true)
	servicenow.Start(binfo.Client, binfo.Task, name)
	b := true
	dst := filepath.Join(volume, binfo.CSNumber)
	if binfo.Task == "Restore" {
		b = false
		dst = volume
	}
	i, s = backup.Backup(binfo.Source, dst, method, conf, b)

	servicenow.Finish(binfo.Client, binfo.Task, name, map[string]interface{}{"Files": i, "Size": conversion.ByteCountSI(s, UNIT, 0)})
	src := filestruct.New(binfo.Source[0].Path)

	const format = "01-02-2006"
	t := time.Now()
	newdir := t.Format(format) + "_" + src.Head
	err := os.Rename(filepath.Join(src.Volume, src.Head), filepath.Join(src.Volume, newdir))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Copied %v files \n", i)
	fmt.Printf("Total size: %v\n", conversion.ByteCountSI(s, UNIT, 0))
	exit()
}
