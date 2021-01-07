// +build windows

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/michaeldcanady/LookBack/OLD/MapDrive"
	structure "github.com/michaeldcanady/LookBack/backup2.0/struct"
)

const (
	PATHSEPARATOR     = "\\"
	PATHLISTSEPARATOR = ';'
	UNIT              = 1024
)

var ()

func Clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func ByteCountSI(b int64) string {
	if b < UNIT {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(UNIT), 0
	for n := b / UNIT; n >= UNIT; n /= UNIT {
		div *= UNIT
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
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
	var r []string
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		f, err := os.Open(string(drive) + ":\\")
		if err == nil {
			r = append(r, string(drive))
			f.Close()
		}
	}
	return r
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

func getName(path string) string {
	var drive, name string
	volume := filepath.VolumeName(path)
	command := fmt.Sprintf("vol %s", volume)
	if c, err := exec.Command("cmd", "/c", command).CombinedOutput(); err != nil {
		log.Fatal(err)
	} else {
		str := strings.Fields(string(c))
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
	}
	return fmt.Sprintf("%s (%s)", name, drive)
}

func mapDrive(drive, blank, b1 string) error {
	err := MapDrive.WNetAddConnection2(drive, blank, b1)
	return err
}

//HAVE IT FACTOR IN FILES THAT NEED TO BE SKIPPED
//Gets size of specified directory
func DirSize(path string, isRoot ...bool) (size int64) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return 0
	}
	for _, entry := range entries {
		if strings.ToLower(entry.Name()) == "appdata" && len(isRoot) > 0 {
			continue
		}
		if strings.ToLower(entry.Name()) == "library" && len(isRoot) > 0 {
			continue
		}
		if entry.IsDir() {
			size += DirSize(filepath.Join(path, entry.Name()))
		} else {
			size += int64(entry.Size())
		}
	}
	return
}
