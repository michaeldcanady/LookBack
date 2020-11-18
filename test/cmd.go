package main

import(
  "io/ioutil"
  "path/filepath"
  //"os"
  "fmt"
)

func GetUsersBackground(user, copyloc string)error{
  backgroundLoc := filepath.Join(user,"\\AppData\\Roaming\\Microsoft\\Windows\\Themes\\TranscodedWallpaper")
  //Read all the contents of the  original file
   bytesRead, err := ioutil.ReadFile(backgroundLoc)
   if err != nil {
       return err
   }
   base := filepath.Base(backgroundLoc)
   if err != nil{
     return err
   }
   //Copy all the contents to the desitination file
   err = ioutil.WriteFile(filepath.Join(copyloc, base+".jpg"), bytesRead, 0755)
   if err != nil {
       return err
   }
  return nil

}

func main(){
  err := GetUsersBackground("C:\\Users\\dmcanady","C:\\Users\\dmcanady\\desktop")
  if err != nil {
    panic(fmt.Sprintf("Getting background Error: %s",err))
  }
}
