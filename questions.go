package main

import(
  "github.com/AlecAivazis/survey/v2"
  term"github.com/AlecAivazis/survey/v2/terminal"
  "github.com/michaeldcanady/Project01/MapDrive"
  "errors"
  "fmt"
  "strings"
  "path/filepath"
  "os"
)

// requests username from user
func getUserName(binfo *backup){
  Heading(binfo)
  err := survey.AskOne(&survey.Input{Message: "What is your Username?"},&binfo.Technician)
  if err == term.InterruptErr {
	   exit()
  } else if err != nil {
	   panic(err)
  }
  if binfo.Technician == "" {
    panic(errors.New("Empty Username"))
  }
}

// requests csnumber from user
func getCSNumber(binfo *backup){
  // ENTER CS NUMBER
  Heading(binfo)
  err := survey.AskOne(&survey.Input{Message: "CSNumber (EX: CS1234567):"},&binfo.CSNumber)
  if err == term.InterruptErr {
	   exit()
  } else if err != nil {
	   panic(err)
  }
  if binfo.CSNumber == "" {
    panic(errors.New("Empty CSNumber"))
  }
}

func getTask(binfo *backup){
  // SELECT TO RESTORE OR BACKUP
  Heading(binfo)
  BorR := 0
  prompt := &survey.Select{
    Message: "Backup or Restore?",
    Options: bORr,
  }
  err := survey.AskOne(prompt, &BorR)
  if err == term.InterruptErr {
	   exit()
  } else if err != nil {
	   panic(err)
  }
  binfo.Task = bORr[BorR]
}

// gets how the user would like to backup information
func getBackupMethod(binfo *backup,tar bool)string{
  // SELECT BACKUP METHOD
  Heading(binfo)
  BorR := 0
  var options = []string{"InLine Copy","Compress"}
  if tar{
    options = append(options,"Tar")
  }
  prompt := &survey.Select{
    Message: strings.Join(options,", "),
    Options: options,
  }
  err := survey.AskOne(prompt, &BorR)
  if err == term.InterruptErr {
	   exit()
  } else if err != nil {
	   panic(err)
  }
  return options[BorR]
}

// User selects all user files to backup
func getSource(binfo *backup){
  // SELECT Source
  var Users []string
  Heading(binfo)
  for _,user := range users{
    Users = append(Users,fmt.Sprintf("%s - size: %s",user.path,ByteCountSI(user.size)))
  }
  binfo.Source = nil
  err := survey.AskOne(
    &survey.MultiSelect{
      Message: "Select what users you would like to backup: ",
      Options: Users,
      }, &binfo.Source)
  if err == term.InterruptErr {
    exit()
  } else if err != nil {
    panic(err)
  }
}

// SELECT destination
func getDestination(binfo *backup){
  Heading(binfo)
  //Map to HDBackups
  var confirnation string
  err := survey.AskOne(
    &survey.Select{
      Message: "backup to: network drive or local drive: ",
      Options: []string{"Network Drive","Local Drive"},
    }, &confirnation)
  if err == term.InterruptErr {
  	exit()
  } else if err != nil {
  	panic(err)
  }
  if confirnation == "Network Drive"{
    binfo.DestType = "Network"
    NetworkDrive(binfo)
  }else if confirnation == "Local Drive"{
    binfo.DestType = "Local"
    LocalDrive(binfo)
  }
}

// Mapping a network drive for
func NetworkDrive(binfo *backup){
  Heading(binfo)
  var netdrive string
  err := survey.AskOne(&survey.Input{Message: "Enter network drive address:"},&netdrive)
  if err == term.InterruptErr {
    exit()
  } else if err != nil {
    panic(err)
  }
  success := MapDrive.MapHDBackupsWindows(binfo.Technician,netdrive)
  if success{
    binfo.Dest = netdrive
  }
}

func DropboxLoc(binfo *backup){
  //Save to shared DropBox files
}

//Let users backup thier data to a local drive
func LocalDrive(binfo *backup){
  Heading(binfo)

  //User selected file location
  // gets a string of available drives
  drive := strings.Join(getDrives(),",")
  binfo.Dest = ""
  err := survey.AskOne(
    &survey.Input{
      Message: fmt.Sprintf("Avaliable drives %s:",drive),
      // Auto complete function for file path
      Suggest: func (toComplete string) []string {
          // gets entered path and gets subfolders and files
          files, _ := filepath.Glob(toComplete + "*")
          // enumerates over the files/folder
          for i,file := range files{
            fi, err := os.Stat(file); if err !=nil{
              panic(err)
            }
            switch mode := fi.Mode(); {
              // if it is a dir it appends \ to the end
              case mode.IsDir():
                files[i] = file+"\\"
              // if its a file it does not output it
              case mode.IsRegular():
                continue
              }
          }
          // returns a list of directories
          return files
      },
    }, &binfo.Dest)
  if err == term.InterruptErr {
  	exit()
  } else if err != nil {
    panic(err)
  }
}

// checks if user confirms entered information
func getConfirmation(binfo *backup)bool{
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

// allows user to select information they wish to change
func SelectChange(binfo *backup)[]string{
  Heading(binfo)
  selected := []string{}
  err := survey.AskOne(
    &survey.MultiSelect{
      Message: "Select what users you would like to backup: ",
      Options: []string{"Username","Ticket Number","Task","Source","Destination"},
      }, &selected)
  if err == term.InterruptErr {
    exit()
  } else if err != nil {
    panic(err)
  }
    return selected
}

// Proceedure for user initiated exit
func exit(){
  fmt.Println("interrupted")
  os.Exit(0)
}
