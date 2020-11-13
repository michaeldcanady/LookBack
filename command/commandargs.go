package main

import(
  "flag"
  "fmt"
  "os/user"
  "strings"
)

var(
  users string
  userSlice []string
)

func main(){
  currentUser,_ := user.Current()
  User := strings.Split(currentUser.Username,"\\")
  flag.StringVar(&users, "users",User[1], "List all users you want backed up.")

  flag.Parse()
  
  userSlice = strings.Split(users," ")
  fmt.Println(userSlice)
}
