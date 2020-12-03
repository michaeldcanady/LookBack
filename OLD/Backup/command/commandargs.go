package main

import(
  "flag"
  "fmt"
  "os/user"
  "strings"
  "path/filepath"
  "os"
)

var(
  //Flag variables
  UsersArgs string
  DestArgs string
  SilentArgs bool
  BorRArgs string
)

//verifys passed args
func VerifyArgs(){

  //Checks that -type is backup or restore
  if BorRArgs != "backup" && BorRArgs != "restore"{
    panic("-type must be either BACKUP or RESTORE")
  }

  // Discards all printf and println statements
  // if -s is included it is set to true
  if SilentArgs{os.Stdout,_ = os.Open(os.DevNull)}
}

func init(){
  // Getting current user for default
  currentUser,_ := user.Current()
  User := strings.Split(currentUser.Username,"\\")
  homedir := currentUser.HomeDir
  desktop := filepath.Join(homedir,"/Desktop/")
  // defining flags
  flag.StringVar(&UsersArgs, "users",User[1], "List all users you want backed up.")
  flag.StringVar(&DestArgs, "dest",desktop, "Set backup destination.")
  flag.StringVar(&BorRArgs, "type","backup", "Backup or Restore")
  flag.BoolVar(&SilentArgs, "s",false, "silent command output.")
}

func main(){
  flag.Parse()

  VerifyArgs()
  // Test output
  userSlice := strings.Split(UsersArgs," ")
  fmt.Println(userSlice)
  fmt.Println(DestArgs)

}
