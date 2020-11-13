// +build windows

package main

import(
  "os/exec"
  "os"
)

func Clear(){
  cmd := exec.Command("cmd", "/c", "cls")
  cmd.Stdout = os.Stdout
  cmd.Run()
}
