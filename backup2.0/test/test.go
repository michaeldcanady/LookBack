package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if c, err := exec.Command("cmd", "/c", "vol C:").CombinedOutput(); err != nil {
		log.Fatal(err)
	} else {
		str := strings.Fields(string(c))
		for i, t := range str {
			if i == 0 {

			} else if str[i-1] == "drive" || str[i-1] == "is" {
				fmt.Println(t)
			}
		}
		//fmt.Println(strings.Spli(string(c), ""))
	}

}
