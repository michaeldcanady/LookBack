// +build windows

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	structure "github.com/michaeldcanady/LookBack/backup2.0/struct"
	winapi "github.com/michaeldcanady/Windows-Api/Windows-Api"
	"github.com/michaeldcanady/Windows-Api/Windows-Api/kernel32"
)

const UNIT = 1024

func Clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// Retrieves the current user's background
func GetUsersBackground(user, copyloc string) (string, error) {
	backgroundLoc := filepath.Join(user, "\\AppData\\Roaming\\Microsoft\\Windows\\Themes\\TranscodedWallpaper")
	//Read all the contents of the  original file
	bytesRead, err := ioutil.ReadFile(backgroundLoc)
	if err != nil {
		return "", err
	}
	base := filepath.Base(backgroundLoc)
	if err != nil {
		return "", err
	}
	//Copy all the contents to the desitination file
	copyloc = filepath.Join(copyloc, base+".jpg")
	err = ioutil.WriteFile(copyloc, bytesRead, 0755)
	if err != nil {
		return "", err
	}
	return copyloc, nil

}

// returns a slice of active drives
func getDrives() []string {
	drives, err := kernel32.GetLogicalDrives()
	if err != nil {
		panic(err)
	}
	return drives
}

func UsersJoin(users []structure.User, seperator string) string {
	var newstring string
	for _, user := range users {
		newstring += user.Path + seperator
	}
	return newstring
}

// Checks if selection is within skippable slice above
func Skippable(selection string) bool {
	for _, skip := range SKIPPABLE {
		if skip == selection {
			return true
		}
	}
	return false
}

//GetUsers returns a list of users on all devices connected to machine
func GetUsers() []structure.User {
	var users []structure.User
	drives := getDrives()
	for _, drive := range drives {
		userdir := drive + ":\\Users"
		if _, err := os.Stat(userdir); os.IsNotExist(err) {
			continue
		} else {
			files, _ := filepath.Glob(userdir + "/**")
			for _, file := range files {
				fi, err := os.Stat(file)
				if err != nil {
					panic(err)
				}
				switch mode := fi.Mode(); {
				case mode.IsDir():
					if Skippable(file) {
						continue
					} else {
						users = append(users, structure.NewUser(file))
					}
				case mode.IsRegular():
					continue
				}
			}
		}
	}
	return users
}

func getName(path string, withLetter bool) string {
	info := kernel32.GetVolumeInformationW(path)
	drive := strings.Replace(info.PathName, "\\", "", -1)
	name := info.VolumeLabel

	if withLetter {
		return fmt.Sprintf("%s (%s)", name, drive)
	}
	return fmt.Sprintf("%s", name)
}

func mapDrive(drive, username, password, volume string) error {
	return winapi.WNetAddConnection2(drive, username, password, volume)
}

//HAVE IT FACTOR IN FILES THAT NEED TO BE SKIPPED
//Gets size of specified directory
func DirSize(path string, isRoot ...bool) (size int64) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return 0
	}
	for _, entry := range entries {
		if entry.IsDir() {
			size += DirSize(filepath.Join(path, entry.Name()))
		} else {
			size += int64(entry.Size())
		}
	}
	return
}
