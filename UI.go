// +build Liberty

package main

import(
  "path/filepath"
  "os"
  "strings"
  "fmt"
  //"golang.org/x/crypto/ssh/terminal"
)

type backup struct{
  Technician string
  Task string
  source []string
  dest string
  CSNumber string
}

var(
  bORr = []string{"Backup", "Restore"}
  // User accounts that don't need to be included in the options
  SKIPPABLE = []string{"C:\\Users\\Default","C:\\Users\\Public","C:\\Users\\All Users","C:\\Users\\Default User"}
)

// returns a slice of active drives
func getDrives()[]string{
var r []string
for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ"{
    f, err := os.Open(string(drive)+":\\")
    if err == nil {
        r = append(r, string(drive))
        f.Close()
    }
}
return r
}

// Creates format for selection made
func Heading(binfo *backup){
  Header()
  fmt.Println(strings.Repeat(" ",SPACESIZE)+"Currently selected options:")
  fmt.Printf(strings.Repeat(" ",SPACESIZE)+"   Technician: %s\n",binfo.Technician)
  fmt.Printf(strings.Repeat(" ",SPACESIZE)+"Ticket Number: %s\n",binfo.CSNumber)
  fmt.Printf(strings.Repeat(" ",SPACESIZE)+"         Task: %s\n",binfo.Task)
  fmt.Printf(strings.Repeat(" ",SPACESIZE)+"       Source: %s\n",strings.Join(binfo.source,","))
  fmt.Printf(strings.Repeat(" ",SPACESIZE)+"  Destination: %s\n",binfo.dest)
  fmt.Println("")
}

// Checks if selection is within skippable slice above
func Skippable(selection string)bool{
  for _,skip := range SKIPPABLE{
    if skip == selection{
      return true
    }
  }
  return false
}

// returns a list of users on all devices connected to machine
func GetUsers()[]string{
  var users []string
  drives := getDrives()
  for _,drive := range drives{
    if _, err := os.Stat(drive+":\\Users"); os.IsNotExist(err) {
      continue
    }else{
      files, _ := filepath.Glob(drive+":\\Users\\" + "*")
      for _,file := range files{
        fi, err := os.Stat(file); if err !=nil{
          panic(err)
        }
        switch mode := fi.Mode(); {
          case mode.IsDir():
            if Skippable(file){
              continue
              }else{
                users = append(users,file)
              }
            case mode.IsRegular():
              continue
        }
      }
    }
  }
  return users
}

//HAVE IT FACTOR IN FILES THAT NEED TO BE SKIPPED
//Gets size of specified directory
func DirSize(path string) (int64, error) {
    var size int64
      _ = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            size += info.Size()
        }
        return err
      })
  return size, nil
}

// UNUSED
type user struct{
  path string
  size int64

}

// UNUSED
func NewUser(path string)user{
  var u user
  u.path = path
  u.size,_ = DirSize(path)
  return u
}

func main(){
  // sets backup's information to none
  binfo := backup{"","",[]string{},"",""}
  // Explained in questions.go
  getUserName(&binfo)
  getCSNumber(&binfo)
  getTask(&binfo)
  getSource(&binfo)
  getDestination(&binfo)
  // Explained in questions.go
  // Checks if user confirmes data or not
  confirm := getConfirmation(&binfo)
  for{
    if confirm == false{
      // checks what fields user wants to change
      selected:= SelectChange(&binfo)
      for _,s := range selected{
        if s == "Username"{
          getUserName(&binfo)
        }else if s == "Ticket Number"{
          getCSNumber(&binfo)
        }else if s == "Task"{
          getTask(&binfo)
        }else if s == "Source"{
          getSource(&binfo)
        }else if s == "Destination"{
          getDestination(&binfo)
        }else{
          continue
        }
      }
      // Checks if user confirmes data or not
      confirm = getConfirmation(&binfo)
      }else{
        fmt.Println("Beginning backing up data.")
        go func(){
          var users []user
          for _,path := range binfo.source{
            users = append(users,NewUser(path))
          }
        }()
        break
      }
    }
    method := getBackupMethod(&binfo,true)
    if method == "InLine Copy"{
      InLineCopy(&binfo)
    }else if method == "Compress"{
      fmt.Println("Currently Not Supported")
    }else if method == "Tar"{
      fmt.Println("Currently Not Supported")
    }
}
