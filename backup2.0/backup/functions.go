package backup

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/michaeldcanady/Project01/backup2.0/file"
	"github.com/michaeldcanady/Project01/backup2.0/struct"
)

func gather(path string, output chan *file.File, conf structure.Config) {
	Use_Exclusions := conf.Settings.Use_Exclusions
	Use_Inclusions := conf.Settings.Use_Inclusions
	Excluded := conf.Exclusions.General_Exclusions
	ExcludedFiles := conf.Exclusions.File_Type_Exclusions
	Included := conf.Inclusions.General_Inclusions

	dirs, _ := filepath.Glob(path + "/**")
	for _, dir := range dirs {
		if !FileCheck(dir, Use_Exclusions, Use_Inclusions, Included, Excluded, ExcludedFiles) {
			continue
		} else {
			fi, err := os.Stat(dir)
			if err != nil {
				check(err, "error")
				continue
			}
			switch mode := fi.Mode(); {
			case mode.IsDir():
				//fmt.Println(dir)
				gather(dir, output, conf)
			case mode.IsRegular():
				output <- file.New(dir, &fi)
			}
		}
	}
}

func getUser(input string) *user.User {
	var count = 2
	if runtime.GOOS != "windows" {
		count = 1
	}
	if User, err := user.LookupId(input); err == nil {
		return User
	} else if User, err := user.Lookup(strings.Split(input, string(os.PathSeparator))[count]); err == nil {
		return User
	} else {
		fmt.Println(strings.Split(input, string(os.PathSeparator)))
		return &user.User{
			Uid:      "",
			Gid:      "",
			Username: strings.Split(input, string(os.PathSeparator))[1],
			Name:     "",
			HomeDir:  input,
		}
	}

}

func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}
