// +build windows darvin linux

package main


import(
  "fmt"
  "os"
  "strings"
  "os/exec"

  "golang.org/x/crypto/ssh/terminal"
)

var(
  fd = int(os.Stdout.Fd())
	termWidth, termHeight, _ = terminal.GetSize(fd)
  LOGOWIDTH = 120
  SPACESIZE = ((termWidth-LOGOWIDTH)/2)
)

func Clear(){
  cmd := exec.Command("cmd", "/c", "cls")
  cmd.Stdout = os.Stdout
  cmd.Run()
}

func Header(){
  Clear()
  fmt.Println(logo(strings.Repeat(" ",SPACESIZE)))
}

func logo(space string)string{
return space+"┌──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐\n"+
       space+"│          ██            ████████      ████████    ██      ██  ██████████      ██████      ██████    ██      ██        │\n"+
       space+"│          ██          ██        ██  ██        ██  ██    ██    ██        ██  ██      ██  ██      ██  ██    ██          │\n"+
       space+"│          ██          ██        ██  ██        ██  ██  ██      ██        ██  ██      ██  ██          ██  ██            │\n"+
       space+"│          ██          ██        ██  ██        ██  ████        ██████████    ██      ██  ██          ████              │\n"+
       space+"│          ██          ██        ██  ██        ██  ██  ██      ██        ██  ██████████  ██          ██  ██            │\n"+
       space+"│          ██          ██        ██  ██        ██  ██    ██    ██        ██  ██      ██  ██      ██  ██    ██          │\n"+
       space+"│          ██████████    ████████      ████████    ██      ██  ██████████    ██      ██    ██████    ██      ██        │\n"+
       space+"├──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤\n"+
       space+"│                                      Developed by Michael Canady                                                     │\n"+
       space+"└──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘"
}
