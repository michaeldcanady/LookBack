package main

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	path := "J:\\Users\\dmcanady"
	volume := filepath.VolumeName(path)
	command := fmt.Sprintf("vol %s", volume)
	if c, err := exec.Command("cmd", "/c", command).CombinedOutput(); err != nil {
		log.Fatal(err)
	} else {
		str := strings.Fields(string(c))
		var drive, name string
		for i, t := range str {
			if i == 0 {

			} else if str[i-1] == "drive" {
				drive = t
			} else if i > 1 && str[i-2] == drive {
				name = t
			} else {
				continue
			}
		}
		fmt.Printf("%s (%s)", name, drive)
		//fmt.Println(strings.Spli(string(c), ""))
	}

}
