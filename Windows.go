// +build windows

package main

import(
  "fmt"
  "path/filepath"
  "io/ioutil"
)

const(
  PATHSEPARATOR = "\\"
  PATHLISTSEPARATOR = ';'
  UNIT = 1024
)

func ByteCountSI(b int64) string {
    if b < UNIT {
        return fmt.Sprintf("%d B", b)
    }
    div, exp := int64(UNIT), 0
    for n := b / UNIT; n >= UNIT; n /= UNIT {
        div *= UNIT
        exp++
    }
    return fmt.Sprintf("%.1f %cB",
        float64(b)/float64(div), "kMGTPE"[exp])
}

// Retrieves the current user's background
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
