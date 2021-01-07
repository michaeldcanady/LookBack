// +build linux

package main

import (
	"fmt"
	"os/exec"
)

func mapDrive(loc, user, pass string) error {
	mediaName := "/media/name"
	command := fmt.Sprintf("'sudo mount -t cifs -o username=%s %s %s'", user, loc, mediaName)
	_, err := exec.Command("/bin/sh", "-c", "sudo mkdir "+mediaName).Output()
	if err != nil {
		return err
	}
	_, err := exec.Command("/bin/sh", "-c", "sudo mkdir "+command+" "+mediaName).Output()
	if err != nil {
		return err
	}
	return nil
}
