package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		exit()
		//DeleteFiles()
	}()
}

// Proceedure for user initiated exit
func exit() {
	path := netdrive
	if conf.Settings.Network_Path != "" {
		path = conf.Settings.Network_Path
	}
	fmt.Println("\nctrl+c pressed...")
	command := fmt.Sprintf("net use /delete %s", path)
	//out, err := exec.Command("net", "use", `\\fs3.liberty.edu\hdbackups`).CombinedOutput()
	//fmt.Println(err)
	_, err := exec.Command("net", "use", "/delete", path).CombinedOutput()
	if err != nil {
		fmt.Printf("Attempted to disconnect from %s using:\n'%s'.\n", path, command)
		fmt.Printf("Verifying no drives with name %s exists.\n", conf.Settings.NetworkDriveName)
		for _, drive := range getDrives() {
			drive := drive + ":"
			if getName(drive, false) == conf.Settings.NetworkDriveName {
				_, err := exec.Command("net", "use", "/delete", drive).CombinedOutput()
				if err != nil {
					fmt.Printf("Error: removing located drive letter (%s) for %s.\n", drive, path)
					os.Exit(0)
				}
				fmt.Println(drive, "was successfully removed.\n")
				os.Exit(0)
			}
		}
		fmt.Println("Could not locate a drive with the name", conf.Settings.NetworkDriveName)
		fmt.Println("Run 'net use' from an elevated prompt and verify the resource was removed")
		os.Exit(0)
	}
}
