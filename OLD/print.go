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

// Creates format for selection made
func Heading(binfo *backup){
  Header()
  fmt.Println(strings.Repeat(" ",SPACESIZE)+"Currently selected options:")
  fmt.Printf(strings.Repeat(" ",SPACESIZE)+"   Technician: %s\n",binfo.Technician)
  fmt.Printf(strings.Repeat(" ",SPACESIZE)+"Ticket Number: %s\n",binfo.CSNumber)
  fmt.Printf(strings.Repeat(" ",SPACESIZE)+"         Task: %s\n",binfo.Task)
  fmt.Printf(strings.Repeat(" ",SPACESIZE)+"       Source: %s\n",strings.Join(binfo.Source,","))
  fmt.Printf(strings.Repeat(" ",SPACESIZE)+"  Destination: %s\n",binfo.Dest)
  fmt.Println("")
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
