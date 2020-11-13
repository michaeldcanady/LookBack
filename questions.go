package main

import(
  "github.com/AlecAivazis/survey/v2"
  term"github.com/AlecAivazis/survey/v2/terminal"
  "errors"
  "fmt"
  "strings"
  "path/filepath"
  "os"
)

func getUserName(binfo *backup){
  Heading(binfo)
  err := survey.AskOne(&survey.Input{Message: "What is your Username?"},&binfo.Technician)
  if err == term.InterruptErr {
	fmt.Println("interrupted")

	os.Exit(0)
  } else if err != nil {
	   panic(err)
  }
  if binfo.Technician == "" {
    panic(errors.New("Empty Username"))
  }
}

func getCSNumber(binfo *backup){
  // ENTER CS NUMBER
  Heading(binfo)
  err := survey.AskOne(&survey.Input{Message: "CSNumber (EX: CS1234567):"},&binfo.CSNumber)
  if err == term.InterruptErr {
	fmt.Println("interrupted")

	os.Exit(0)
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
	fmt.Println("interrupted")

	os.Exit(0)
  } else if err != nil {
	   panic(err)
  }
  binfo.Task = bORr[BorR]
}

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
	fmt.Println("interrupted")

	os.Exit(0)
  } else if err != nil {
	   panic(err)
  }
  return options[BorR]
}

func getSource(binfo *backup){
  // SELECT Source
  Heading(binfo)
  users := GetUsers()
  binfo.source = nil
  err := survey.AskOne(
    &survey.MultiSelect{
      Message: "Select what users you would like to backup: ",
      Options: users,
      }, &binfo.source)
  if err == term.InterruptErr {
    fmt.Println("interrupted")
    os.Exit(0)
  } else if err != nil {
    panic(err)
  }
}

func getDestination(binfo *backup){
  // SELECT destination
  Heading(binfo)
  fmt.Println(binfo.source)
  drives := getDrives()
  drive := strings.Join(drives,",")
  binfo.dest = ""
  err := survey.AskOne(
    &survey.Input{
      Message: fmt.Sprintf("Avaliable drives %s:",drive),
      Suggest: func (toComplete string) []string {
          files, _ := filepath.Glob(toComplete + "*")
          for i,file := range files{
            fi, err := os.Stat(file); if err !=nil{
              panic(err)
            }
            switch mode := fi.Mode(); {
              case mode.IsDir():
                files[i] = file+"\\"
              case mode.IsRegular():
                continue
              }
          }
          return files
      },
    }, &binfo.dest)
  if err == term.InterruptErr {
  	fmt.Println("interrupted")

  	os.Exit(0)
  } else if err != nil {
    panic(err)
  }
}

func getConfirmation(binfo *backup)bool{
  //CONFIRM Information
  Heading(binfo)
  confirnation := false
  err := survey.AskOne(
    &survey.Confirm{
      Message: "Confirm above information?",
    }, &confirnation)
  if err == term.InterruptErr {
  	fmt.Println("interrupted")

  	os.Exit(0)
  } else if err != nil {
  	panic(err)
  }
    return confirnation
}

func SelectChange(binfo *backup)[]string{
  Heading(binfo)
  selected := []string{}
  err := survey.AskOne(
    &survey.MultiSelect{
      Message: "Select what users you would like to backup: ",
      Options: []string{"Username","Ticket Number","Task","Source","Destination"},
      }, &selected)
  if err == term.InterruptErr {
    fmt.Println("interrupted")

    os.Exit(0)
  } else if err != nil {
    panic(err)
  }
    return selected
}
