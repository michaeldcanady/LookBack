package main

import (
	"errors"
	"log"
	"strings"

	"github.com/AlecAivazis/survey"
	term "github.com/AlecAivazis/survey/terminal"
	structure "github.com/michaeldcanady/LookBack/backup2.0/struct"
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
	//binfo.Client = servicenow.Create(binfo.Technician, binfo.Password, conf.Tktsystem.URL, binfo.CSNumber)
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

// SELECT destination
func getDestination(binfo *structure.Backup) {
	switch task := binfo.Task; task {
	case "Backup":
		backupDest(binfo)
	case "Restore":
		restoreDest(binfo)
	default:

	}
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
