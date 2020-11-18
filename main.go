//+ build windows darvin linux

package main

import(
  "path/filepath"
  "os"
  "strings"
  "fmt"
  "io/ioutil"
  "github.com/BurntSushi/toml"
  //"golang.org/x/crypto/ssh/terminal"
)

var(
  // User accounts that don't need to be included in the options
  SKIPPABLE = []string{"C:\\Users\\Default","C:\\Users\\Public","C:\\Users\\All Users","C:\\Users\\Default User"}
  users []user
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

func UsersJoin(users []user,seperator string)string{
  var newstring string
  for _,user := range users{
    newstring += user.path+seperator
  }
  return newstring
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
func GetUsers()[]user{
  var users []user
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
                users = append(users,NewUser(file))
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
func DirSize(path string,isRoot ...bool) (size int64) {
  entries, err := ioutil.ReadDir(path)
  if err != nil {
    return 0
  }
  for _, entry := range entries {
    if strings.ToLower(entry.Name()) == "appdata" && len(isRoot) > 0 {
      continue
    }
    if strings.ToLower(entry.Name()) == "library" && len(isRoot) > 0 {
      continue
    }
    if entry.IsDir() {
      size += DirSize(filepath.Join(path, entry.Name()))
    } else {
      size += int64(entry.Size())
    }
  }
  return
}

func init(){
  users = GetUsers()
  if _, err := toml.DecodeFile("C:\\go\\src\\github.com\\michaeldcanady\\Project01\\Config\\Settings.toml", &conf); err != nil {
    panic(err)
  }
}

func main(){
  // sets backup's information to none
  binfo := backup{"","","",[]string{},"",""}
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
          for _,path := range binfo.Source{
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
