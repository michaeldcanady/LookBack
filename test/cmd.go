package main

import(
  "io/ioutil"
  "path/filepath"
  //"os"
  "fmt"
  "time"
  "os/exec"
)

func GetUsersBackground(user, copyloc string)(string,error){
  backgroundLoc := filepath.Join(user,"\\AppData\\Roaming\\Microsoft\\Windows\\Themes\\TranscodedWallpaper")
  //Read all the contents of the  original file
   bytesRead, err := ioutil.ReadFile(backgroundLoc)
   if err != nil {
       return "", err
   }
   base := filepath.Base(backgroundLoc)
   if err != nil{
     return "", err
   }
   //Copy all the contents to the desitination file
   copyloc = filepath.Join(copyloc, base+".jpg")
   err = ioutil.WriteFile(copyloc, bytesRead, 0755)
   if err != nil {
       return "", err
   }
  return copyloc,nil

}

func RunCommand(command string)(string,error){
  output, err := exec.Command("Powershell", "-Command", command ).CombinedOutput()
  if err != nil{
    return "", err
  }
  return string(output),nil
}

func ChangeBackground(file string)error{
  command := fmt.Sprintf(`reg add "HKEY_CURRENT_USER\Control Panel\Desktop" /v Wallpaper /t REG_SZ /d %s /f`, file)
  fmt.Println(command)
  _,err := RunCommand(command)
  if err != nil {
    return err
  }
  _,err = RunCommand("RUNDLL32.EXE user32.dll,UpdatePerUserSystemParameters")
  if err != nil {
    return err
  }
  return nil

}

func main(){
  copyloc,err := GetUsersBackground("C:\\Users\\dmcanady","C:\\Users\\dmcanady\\desktop")
  if err != nil {
    panic(fmt.Sprintf("Getting background Error: %s",err))
  }
  fmt.Println("Sleeping")
  time.Sleep(10 * time.Second)
  err = ChangeBackground(copyloc)
  if err != nil {
    panic(fmt.Sprintf("Setting background Error: %s",err))
  }
}
