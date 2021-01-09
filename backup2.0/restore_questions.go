package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey"
	structure "github.com/michaeldcanady/LookBack/backup2.0/struct"
)

func restoreSource(binfo *structure.Backup) {
	Heading(binfo)
	var confirm string
	err := survey.AskOne(
		&survey.Select{
			Message: "Restore from: network drive or local drive: ",
			Options: []string{"Network Drive", "Local Drive"},
		}, &confirm)
	errCheck(err)
	switch confirm {
	case "Network Drive":
		restoreNetDrive(binfo)
	case "Local Drive":
		restoreLocDrive(binfo)
	default:
	}
}

func restoreNetDrive(binfo *structure.Backup) {
	Heading(binfo)
	var netdrive string
	if conf.Settings.Network_Path != "" {
		netdrive = filepath.Join(conf.Settings.Network_Path, conf.Settings.NetworkFolderPath)
	} else {
		err := survey.AskOne(&survey.Input{Message: "Enter network drive address:"}, &netdrive)
		errCheck(err)
	}
	err := mapDrive(netdrive, binfo.Technician+"@"+conf.Settings.Email_Extension, binfo.Password, "M:")
	if err != nil {
		panic("Failed to bind drive")
	} else {
		files, err := filepath.Glob(netdrive + "/**")
		errCheck(err)
		var loc string
		for _, file := range files {
			if _, CS := filepath.Split(file); CS == binfo.CSNumber {
				loc = file
			}
		}
		if loc == "" {
			fmt.Printf("Could not file %s in %s. Please try another CS Number", binfo.CSNumber, netdrive)
			restoreSource(binfo)
		} else {
			users, _ := filepath.Glob(loc + "/**")
			for _, user := range users {
				if strings.Contains(user, "logs") {
					continue
				}
				u := structure.NewUser(user)
				binfo.Source = append(binfo.Source, u)
			}
		}
	}
}

func restoreLocDrive(binfo *structure.Backup) {
	Heading(binfo)

	//User selected file location
	// gets a string of available drives
	drive := strings.Join(getDrives(), ",")
	var val string
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
						files[i] = file + "\\"
					// if its a file it does not output it
					case mode.IsRegular():
						continue
					}
				}
				// returns a list of directories
				return files
			},
		}, &val)
	errCheck(err)
	if !strings.HasSuffix(val, string(os.PathSeparator)) {
		val = val + string(os.PathSeparator)
	}
	files, err := filepath.Glob(val + "/**")
	errCheck(err)
	var loc string
	for _, file := range files {
		if _, CS := filepath.Split(file); CS == binfo.CSNumber {
			loc = file
		}
	}
	if loc == "" {
		fmt.Printf("Could not file %s in %s. Please try another CS Number", binfo.CSNumber, val)
		restoreSource(binfo)
	} else {
		users, _ := filepath.Glob(loc + "/**")
		for _, user := range users {
			if strings.Contains(user, "logs") {
				continue
			}
			u := structure.NewUser(user)
			binfo.Source = append(binfo.Source, u)
		}
	}
}

func restoreDest(binfo *structure.Backup) {
	Heading(binfo)
	binfo.Dest = "C:\\Users"
}
