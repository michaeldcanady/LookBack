package main

import(
  "time"
  "fmt"
)

func main(){

  OnStart("LookBack","LookBack","LookBack.exe")
  base := CreateStructure("LookBack",".ssh","Config")
  time.Sleep(time.Second*20)
  fmt.Println("Creating shortcut")
  err := CreateShortcut(base,"LookBack")
  if err != nil{
    panic(err)
  }
  // Initial settings questions here
}
