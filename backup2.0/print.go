// +build windows darwin linux

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	structure "github.com/michaeldcanady/LookBack/backup2.0/struct"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	fd                       = int(os.Stdout.Fd())
	termWidth, termHeight, _ = terminal.GetSize(fd)
	LOGOWIDTH                = 120
	SPACESIZE                = ((termWidth - LOGOWIDTH) / 2)
	OLD                      structure.Backup
	red                      = color.New(color.FgGreen).SprintFunc()
	white                    = color.New(color.FgWhite).SprintFunc()
	color1                   = white
	color2                   = white
	color3                   = white
	color4                   = white
	color5                   = white
)

func createDisplaySourceFormat(users []structure.User, dilimeter string) string {
	var joinedString string
	for i, user := range users {
		if i != 0 {
			joinedString += ", "
		}
		joinedString += user.Path
	}
	return joinedString
}

func ColorCheck(binfo *structure.Backup) {
	if (*binfo).Technician != OLD.Technician {
		color1 = red
	} else {
		color1 = white
	}

	if (*binfo).CSNumber != OLD.CSNumber {
		color2 = red
	} else {
		color2 = white
	}

	if (*binfo).Task != OLD.Task {
		color3 = red
	} else {
		color3 = white
	}

	if createDisplaySourceFormat((*binfo).Source, ",") != createDisplaySourceFormat(OLD.Source, ",") {
		color4 = red
	} else {
		color4 = white
	}

	if (*binfo).Dest != OLD.Dest {
		color5 = red
	} else {
		color5 = white
	}
}

// Creates format for selection made
func Heading(binfo *structure.Backup) {
	if SPACESIZE < 0 {
		SPACESIZE = 1
	}
	Header()
	ColorCheck(binfo)

	fmt.Println(strings.Repeat(" ", SPACESIZE) + "Currently selected options:")
	fmt.Printf(strings.Repeat(" ", SPACESIZE)+"   Technician: %s\n", color1(binfo.Technician))
	fmt.Printf(strings.Repeat(" ", SPACESIZE)+"Ticket Number: %s\n", color2(binfo.CSNumber))
	fmt.Printf(strings.Repeat(" ", SPACESIZE)+"         Task: %s\n", color3(binfo.Task))
	fmt.Printf(strings.Repeat(" ", SPACESIZE)+"       Source: %s\n", color4(createDisplaySourceFormat(binfo.Source, ",")))
	fmt.Printf(strings.Repeat(" ", SPACESIZE)+"  Destination: %s\n", color5(binfo.Dest))
	fmt.Println("")

	OLD = (*binfo)
}

func Header() {
	Clear()
	fmt.Println(logo(strings.Repeat(" ", SPACESIZE)))
}

func logo(space string) string {
	return space + "┌──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐\n" +
		space + "│          ██            ████████      ████████    ██      ██  ██████████      ██████      ██████    ██      ██        │\n" +
		space + "│          ██          ██        ██  ██        ██  ██    ██    ██        ██  ██      ██  ██      ██  ██    ██          │\n" +
		space + "│          ██          ██        ██  ██        ██  ██  ██      ██        ██  ██      ██  ██          ██  ██            │\n" +
		space + "│          ██          ██        ██  ██        ██  ████        ██████████    ██      ██  ██          ████              │\n" +
		space + "│          ██          ██        ██  ██        ██  ██  ██      ██        ██  ██████████  ██          ██  ██            │\n" +
		space + "│          ██          ██        ██  ██        ██  ██    ██    ██        ██  ██      ██  ██      ██  ██    ██          │\n" +
		space + "│          ██████████    ████████      ████████    ██      ██  ██████████    ██      ██    ██████    ██      ██        │\n" +
		space + "├──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤\n" +
		space + "│                                      Developed by Michael Canady                                                     │\n" +
		space + "│                                            with help from                                                            │\n" +
		space + "│                              David Hunter, Mike Rose, and Preston Gibbs                                              │\n" +
		space + "└──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘"
}
