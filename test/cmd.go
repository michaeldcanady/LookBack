package main

import(
  "os/exec"
  "fmt"
)



func main(){
  commandString := "RUNAS /trustlevel:0x20000 'net use'"
  output, err := exec.Command("Powershell", "-Command", commandString ).CombinedOutput()
  if err != nil{
    fmt.Println(err)
  }
  fmt.Println(string(output))
}
