//+ build windows darwin linux

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/michaeldcanady/Project01/backup2.0/backup"
	"github.com/michaeldcanady/Project01/backup2.0/conversion"
	"github.com/michaeldcanady/Project01/backup2.0/servicenow"
	structure "github.com/michaeldcanady/Project01/backup2.0/struct"
)

var (
	//SKIPPABLE User accounts that don't need to be included in the options
	SKIPPABLE = []string{"C:\\Users\\Default", "C:\\Users\\Public", "C:\\Users\\All Users", "C:\\Users\\Default User"}
	users     []structure.User
	conf      structure.Config
	MAX       int64
)

func init() {
	test()
	if runtime.GOOS == "windows" {
		MAX = conf.Settings.WinServerBackupMax
	} else {
		MAX = conf.Settings.MacServerBackupMax
	}
}

func main() {
	equinoxUpdate()
	// sets backup's information to none
	binfo := structure.Backup{"", "", servicenow.Back{}, "", "", []structure.User{}, "", ""}
	// Explained in questions.go
	getUserName(&binfo)
	getPassword(&binfo)
	getCSNumber(&binfo)
	getTask(&binfo)
	getSource(&binfo)
	getDestination(&binfo)
	// Explained in questions.go
	// Checks if user confirmes data or not
	confirm := getConfirmation(&binfo)
	for {
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
			// Checks if user confirms data or not
			confirm = getConfirmation(&binfo)
		} else {
			break
		}
		//	fmt.Println("Beginning backing up data.")
		//	go func() {
		//		var users []structure.User
		//		for _, path := range binfo.Source {
		//			users = append(users, structure.NewUser(path, map[string]int{"C:\\Users\\dmcanady": 0}))
		//		}
		//	}()
		//	break
		//}
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
	name := getName(binfo.Dest)
	servicenow.Start(binfo.Client, binfo.Task, name)
	b := true
	if binfo.Task == "Restore" {
		b = false
	}
	i, s = backup.Backup(binfo.Source, filepath.Join(volume, binfo.CSNumber), method, name, conf, b)

	servicenow.Finish(binfo.Client, binfo.Task, name, map[string]interface{}{"Files": i, "Size": conversion.ByteCountSI(s, UNIT, 0)})
	fmt.Printf("Copied %v files \n", i)
	fmt.Printf("Total size: %v\n", conversion.ByteCountSI(s, UNIT, 0))
}
