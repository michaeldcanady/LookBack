//+ build windows darvin linux

package main

import(
  "fmt"
  "path/filepath"

  "github.com/BurntSushi/toml"
  "github.com/michaeldcanady/Project01/OLD/restore"
  //"golang.org/x/crypto/ssh/terminal"
)

var(
  // User accounts that don't need to be included in the options
  SKIPPABLE = []string{"C:\\Users\\Default","C:\\Users\\Public","C:\\Users\\All Users","C:\\Users\\Default User"}
  users []user
)

func init(){
  users = GetUsers()
  if _, err := toml.DecodeFile(filepath.Join("C:\\","go","src","github.com","michaeldcanady","Project01","OLD","Config","Settings.toml"), &conf); err != nil {
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
      // Checks if user confirms data or not
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
