package MapDrive

import(
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
)

// Returns a slice of avaliable drives
func AvaliableDrives(available bool)[]string{
  var r []string
  for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ"{
      f, err := os.Open(string(drive)+":\\")
      if err != nil && available {
          r = append(r, string(drive))
          f.Close()
      }else if err == nil && available == false{
        r = append(r, string(drive))
        f.Close()
      }
  }
  return r
}

//WORK IN PROGRESS
//need to find a way to get the name for the drive such as D://AnyUse
func CheckMapped(){
  Drives := AvaliableDrives(false)
  for _,drive := range Drives{
    fmt.Println(drive,filepath.VolumeName(drive))
  }
}

func MapHDBackupsWindows(username ,netdrive string) bool {
	temp := "/user:" + username
  Drives := AvaliableDrives(true)
	_, err := exec.Command("net", "use", Drives[0]+":", netdrive, temp).CombinedOutput()
	if err != nil {
		fmt.Println("Error mapping server. Open a CMD", "AS ADMIN", "and run this command:")
		fmt.Println()
		fmt.Println("net use "+Drives[0]+":",netdrive, temp)
		fmt.Println()
		fmt.Println("It will have you type out your password. If it completed successfully, leave the window open")
		fmt.Println("and rerun this program")
		fmt.Println()
    return false
	}else{
    return true
  }
}
