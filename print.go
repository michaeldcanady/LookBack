// +build windows darvin linux

package main


import(
  "fmt"
  "golang.org/x/crypto/ssh/terminal"
  "os"
  "strings"
)

var(
  fd = int(os.Stdout.Fd())
	termWidth, termHeight, _ = terminal.GetSize(fd)
  LOGOWIDTH = 120
  SPACESIZE = ((termWidth-LOGOWIDTH)/2)
)

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
