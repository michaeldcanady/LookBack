/*
test purposes
*/

package main

import (
	"runtime"
	"os/exec"
	"fmt"
	"os"
)

func exit(err error){
	input := ""
	fmt.Println("error in dropbox.go")
	fmt.Println(err)
	fmt.Println("Press enter to exit")
	fmt.Scanln(&input)
	os.Exit(0)
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
	  exit(err)
	}

}

func main() {

openbrowser("https://dropbox.com/login/")
/*

*/
login_confirmation := ""
fmt.Println("Please log into dropbox")
fmt.Println("Press enter when logged into dropbox")
fmt.Scanln(&login_confirmation)
	}
