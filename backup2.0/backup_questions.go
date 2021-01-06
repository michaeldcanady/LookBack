package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey"
	term "github.com/AlecAivazis/survey/terminal"
	"github.com/michaeldcanady/Project01/backup2.0/conversion"
	structure "github.com/michaeldcanady/Project01/backup2.0/struct"
)

var (
	TOTALBACKUPSIZE int64
)

// Mapping a network drive for
func NetworkDrive(binfo *structure.Backup) {
	Heading(binfo)
	var netdrive string
	if conf.Settings.Network_Path != "" {
		binfo.Dest = conf.Settings.Network_Path
	} else {
		err := survey.AskOne(&survey.Input{Message: "Enter network drive address:"}, &netdrive)
		errCheck(err)
		err = mapDrive(netdrive, "", "")
		if err != nil {
			err = mapDrive(netdrive, binfo.Technician, binfo.Password)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
		} else {
			binfo.Dest = netdrive
		}
	}
}

func DropboxLoc(binfo *structure.Backup) {
	//Save to shared DropBox files
}

//Let users Backup thier data to a local drive
func LocalDrive(binfo *structure.Backup) {
	Heading(binfo)

	//User selected file location
	// gets a string of available drives
	drive := strings.Join(getDrives(), ",")
	binfo.Dest = ""
	err := survey.AskOne(
		&survey.Input{
			Message: fmt.Sprintf("Avaliable drives %s:", drive),
			// Auto complete function for file path
			Suggest: func(toComplete string) []string {
				// gets entered path and gets subfolders and files
				files, _ := filepath.Glob(toComplete + "*")
				// enumerates over the files/folder
				for i, file := range files {
					fi, err := os.Stat(file)
					if err != nil {
						panic(err)
					}
					switch mode := fi.Mode(); {
					// if it is a dir it appends \ to the end
					case mode.IsDir():
						files[i] = file + "/"
					// if its a file it does not output it
					case mode.IsRegular():
						continue
					}
				}
				// returns a list of directories
				return files
			},
		}, &binfo.Dest)
	errCheck(err)
}

func backupSource(binfo *structure.Backup) {
	// SELECT Source
	users = GetUsers()
	var Users []string
	var num []int
	Heading(binfo)
	for _, user := range users {
		Users = append(Users, fmt.Sprintf("%s - %v", user.Path, conversion.ByteCountSI(user.Size, UNIT, 0)))
	}
	binfo.Source = nil
	err := survey.AskOne(
		&survey.MultiSelect{
			Message: "Select what users you would like to Backup: ",
			Options: Users,
		}, &num)

	for _, n := range num {
		binfo.Source = append(binfo.Source, users[n].Path)
		TOTALBACKUPSIZE += users[n].Size
	}

	errCheck(err)
}

// gets how the user would like to Backup information
func getBackupMethod(binfo *structure.Backup, tar bool) string {
	// SELECT Backup METHOD
	Heading(binfo)
	BorR := 0
	var options = []string{"InLine Copy", "Compress"}
	if tar {
		options = append(options, "Zip")
	}
	prompt := &survey.Select{
		Message: strings.Join(options, ", "),
		Options: options,
	}
	err := survey.AskOne(prompt, &BorR)
	errCheck(err)
	return options[BorR]
}

func backupDest(binfo *structure.Backup) {
	Heading(binfo)
	var confirnation string
	err := survey.AskOne(
		&survey.Select{
			Message: "Backup to: network drive or local drive: ",
			Options: []string{"Network Drive", "Local Drive"},
		}, &confirnation)
	errCheck(err)
	if confirnation == "Network Drive" {
		var max int64
		if runtime.GOOS == "windows" {
			max = conf.Settings.WinServerBackupMax
		} else {
			max = conf.Settings.MacServerBackupMax
		}
		if TOTALBACKUPSIZE > max {
			if !SizeWarn(binfo, TOTALBACKUPSIZE, max) {
				backupDest(binfo)
			}
			binfo.DestType = "Network"
			NetworkDrive(binfo)
		}
	} else if confirnation == "Local Drive" {
		binfo.DestType = "Local"
		LocalDrive(binfo)
	}
}

func SizeWarn(binfo *structure.Backup, size, max int64) bool {
	//CONFIRM Information
	Heading(binfo)
	confirm := false
	err := survey.AskOne(
		&survey.Confirm{
			Message: fmt.Sprintf("Backup Size is %s which is greater than maximum (%v)", conversion.ByteCountSI(size, UNIT, 0),
				conversion.ByteCountSI(max, UNIT, 0)+"\nPlease confirm you have recieved permission from T3+:"),
		}, &confirm)
	if err == term.InterruptErr {
		exit()
	} else if err != nil {
		panic(err)
	}
	return confirm
}
