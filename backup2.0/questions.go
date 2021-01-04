package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/AlecAivazis/survey"
	term "github.com/AlecAivazis/survey/terminal"
	"github.com/michaeldcanady/Project01/backup2.0/servicenow"
	"github.com/michaeldcanady/Project01/backup2.0/struct"
)

// requests username from user
func getUserName(binfo *structure.Backup) {
	Heading(binfo)
	err := survey.AskOne(&survey.Input{Message: "What is your Username?"}, &binfo.Technician, survey.WithValidator(survey.Required))
	errCheck(err)
}

// requests password from user
func getPassword(binfo *structure.Backup) {
	Heading(binfo)
	err := survey.AskOne(&survey.Password{Message: "What is your Password?"}, &binfo.Password, survey.WithValidator(survey.Required))
	errCheck(err)
}

// requests csnumber from user
func getCSNumber(binfo *structure.Backup) {
	// ENTER CS NUMBER
	Heading(binfo)
	q := survey.Question{Prompt: &survey.Input{Message: "CSNumber (EX: CS1234567):"},
		Validate: func(val interface{}) error {
			switch str, ok := val.(string); {
			case !ok:
				return errors.New("Invalid assertion")
			case len(str) != 9:
				return errors.New("Number not 9 digits")
			case !servicenow.Validate(servicenow.Create(binfo.Technician, binfo.Password, conf.Tktsystem.URL, str)):
				return fmt.Errorf("Invalid CS Number (i.e., ticket closed, wrong CS Number, or not assigned to %s)", binfo.Technician)
			case !strings.HasPrefix(str, "CS"):
				return errors.New("Number must have a prefix of CS")
			default:
				return nil
			}
		}}
	qs := []*survey.Question{&q}
	err := survey.Ask(qs, &binfo.CSNumber)
	binfo.Client = servicenow.Create(binfo.Technician, binfo.Password, conf.Tktsystem.URL, binfo.CSNumber)
	errCheck(err)
}

func getTask(binfo *structure.Backup) {
	// SELECT TO RESTORE OR Backup
	bORr := []string{"Backup", "Restore"}
	Heading(binfo)
	BorR := 0
	prompt := &survey.Select{
		Message: "Backup or Restore?",
		Options: bORr,
	}
	err := survey.AskOne(prompt, &BorR)
	errCheck(err)
	binfo.Task = bORr[BorR]
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

// User selects all user files to Backup
func getSource(binfo *structure.Backup) {
	switch task := binfo.Task; task {
	case "Backup":
		backupSource(binfo)
	case "Restore":
		restoreSource(binfo)
	default:

	}
}

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
		binfo.Dest = conf.Settings.Network_Path
	} else {
		err := survey.AskOne(&survey.Input{Message: "Enter network drive address:"}, &netdrive)
		errCheck(err)
		err = mapDrive(netdrive, binfo.Technician+"@"+conf.Settings.Email_Extension, binfo.Password)
		if err != nil {
			panic("Failed to bind drive")
		} else {
			binfo.Source = append(binfo.Source, netdrive)
		}
	}
}

func restoreLocDrive(binfo *structure.Backup) {
	Heading(binfo)

	//User selected file location
	// gets a string of available drives
	var drive string
	if runtime.GOOS == "Windows" {
		drive = strings.Join(getDrives(), ",")
	} else {

	}
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
	binfo.Source = append(binfo.Source, val)
}

func backupSource(binfo *structure.Backup) {
	// SELECT Source
	var Users []string
	Heading(binfo)
	for _, user := range users {
		Users = append(Users, fmt.Sprintf("%s", user.Path))
	}
	binfo.Source = nil
	err := survey.AskOne(
		&survey.MultiSelect{
			Message: "Select what users you would like to Backup: ",
			Options: Users,
		}, &binfo.Source)
	errCheck(err)
}

// SELECT destination
func getDestination(binfo *structure.Backup) {
	Heading(binfo)
	//Map to HDBackups
	var confirnation string
	err := survey.AskOne(
		&survey.Select{
			Message: "Backup to: network drive or local drive: ",
			Options: []string{"Network Drive", "Local Drive"},
		}, &confirnation)
	errCheck(err)
	if confirnation == "Network Drive" {
		binfo.DestType = "Network"
		NetworkDrive(binfo)
	} else if confirnation == "Local Drive" {
		binfo.DestType = "Local"
		LocalDrive(binfo)
	}
}

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
			panic(err)
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

// checks if user confirms entered information
func getConfirmation(binfo *structure.Backup) bool {
	//CONFIRM Information
	Heading(binfo)
	confirnation := false
	err := survey.AskOne(
		&survey.Confirm{
			Message: "Confirm above information?",
		}, &confirnation)
	if err == term.InterruptErr {
		exit()
	} else if err != nil {
		panic(err)
	}
	return confirnation
}

func getUpdateConf() bool {
	//CONFIRM Information
	confirnation := false
	err := survey.AskOne(
		&survey.Confirm{
			Message: "Newer version avaliable would you like to update?",
		}, &confirnation)
	if err == term.InterruptErr {
		exit()
	} else if err != nil {
		panic(err)
	}
	return confirnation
}

// allows user to select information they wish to change
func SelectChange(binfo *structure.Backup) []string {
	Heading(binfo)
	selected := []string{}
	err := survey.AskOne(
		&survey.MultiSelect{
			Message: "Select what users you would like to Backup: ",
			Options: []string{"Username", "Ticket Number", "Task", "Source", "Destination"},
		}, &selected)
	if err == term.InterruptErr {
		exit()
	} else if err != nil {
		panic(err)
	}
	return selected
}

func errCheck(err error) {
	switch err {
	case term.InterruptErr:
		exit()
	case nil:

	default:
		log.Fatal(err)
	}
}

// Proceedure for user initiated exit
func exit() {
	fmt.Println("interrupted")
	os.Exit(0)
}
