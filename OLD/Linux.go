// +build linux

package main

import(

)

const(
  PATHSEPARATOR = '/'
  PATHLISTSEPARATOR = ':'
)

func Clear(){
  c := exec.Command("clear")
  c.Stdout = os.Stdout
  c.Run()
}
