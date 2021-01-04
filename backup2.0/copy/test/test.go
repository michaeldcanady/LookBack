package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	if c, err := exec.Command("powershell", "-Command", `echo F | xcopy "C:\Users\dmcanady\AppData\Local\Google\Chrome\User Data\Default\Cache\data_2" "C:\CS0085333\dmcanady\AppData\Local\Google\Chrome\User Data\Default\Cache\data_2" /Y`).CombinedOutput(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", c)
	}
}
