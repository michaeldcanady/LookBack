//+ build windows darwin linux

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/michaeldcanady/Project01/backup2.0/backup"
	"github.com/michaeldcanady/Project01/backup2.0/restore"
	"github.com/michaeldcanady/Project01/backup2.0/servicenow"
	"github.com/michaeldcanady/Project01/backup2.0/struct"
	//"golang.org/x/crypto/ssh/terminal"
)

var (
	//SKIPPABLE User accounts that don't need to be included in the options
	SKIPPABLE = []string{"C:\\Users\\Default", "C:\\Users\\Public", "C:\\Users\\All Users", "C:\\Users\\Default User"}
	users     []structure.User
	conf      structure.Config
)

func init() {
	users = GetUsers()
	test()
}

func main() {
	equinoxUpdate()
	// sets backup's information to none
	binfo := structure.Backup{"", "", servicenow.Back{}, "", "", []string{}, "", ""}
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
			fmt.Println("Beginning backing up data.")
			go func() {
				var users []structure.User
				for _, path := range binfo.Source {
					users = append(users, structure.NewUser(path))
				}
			}()
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
	servicenow.Start(binfo.Client, binfo.Task)
	fmt.Println(binfo.Task)
	if binfo.Task == "Backup" {
		backup.Backup(binfo.Source, binfo.Client, filepath.Join(volume, binfo.CSNumber), UNIT, conf, method)
	} else if binfo.Task == "Restore" {
		restore.Restore(binfo.Source, binfo.Client, binfo.Dest, UNIT, conf, method)
	}
}
