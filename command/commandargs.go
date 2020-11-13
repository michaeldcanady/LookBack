package main

import(
  "flag"
  "fmt"
  "os/user"
  "strings"
  "path/filepath"
)

var(
  users string
  destination string
  silent bool
  BorR string
)

func main(){
  // Getting current user for default
  currentUser,_ := user.Current()
  User := strings.Split(currentUser.Username,"\\")
  homedir := currentUser.HomeDir
  desktop := filepath.Join(homedir,"/Desktop/")
  // defining flags
  flag.StringVar(&users, "users",User[1], "List all users you want backed up.")
  flag.StringVar(&destination, "dest",desktop, "Set backup destination.")
  flag.StringVar(&BorR, "type","backup", "Backup or Restore")
  flag.BoolVar(&silent, "s",false, "silent command output.")

  flag.Parse()

  if BorR != "backup" || BorR != "restore"{
    panic("-type must be either BACKUP or RESTORE")
  }

  if !silent{
    userSlice := strings.Split(users," ")
    fmt.Println(userSlice)
    fmt.Println(destination)
  }
}
