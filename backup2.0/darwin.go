// +build darwin linux

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	structure "github.com/michaeldcanady/Project01/backup2.0/struct"
)

const (
	PATHSEPARATOR     = '/'
	PATHLISTSEPARATOR = ':'
	UNIT              = 1000
)

var ()

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func getDrives() []string {
	volumes, err := filepath.Glob("/Volumes/**")
	checkerr(err)
	return volumes
}

func getName(path string) string {
	return "TEST"
}

func mapDrive(loc, user, pass string) error {
	command := fmt.Sprintf("'smb://%s:%s@%s'", user, pass, loc)
	_, err := exec.Command("/bin/sh", "-c", "open "+command).Output()
	if err != nil {
		return err
	}
	return nil
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

// Checks if selection is within skippable slice above
func Skippable(selection string) bool {
	for _, skip := range SKIPPABLE {
		if skip == selection {
			return true
		}
	}
	return false
}

func GetUsers() []structure.User {
	var users []structure.User
	userdir := "/Users/"
	if _, err := os.Stat(userdir); os.IsNotExist(err) {
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
	return users
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

func Clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
