package main

import(
  "fmt"
  "os"
)

func AvaliableDrives()[]string{
  var r []string
  for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ"{
      f, err := os.Open(string(drive)+":\\")
      if err != nil {
          r = append(r, string(drive))
          f.Close()
      }
  }
  return r
  }
}


func main(){
  fmt.Println(AvaliableDrives())
}
