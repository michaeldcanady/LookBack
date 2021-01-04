package filestruct

import (
	"fmt"

	"os"

	"strings"

	"path/filepath"
)

//FilePath struct breaks a recieved filepath into usable compotents in a give situation

type FilePath struct {
	Volume string

	User string

	Head string

	Tail []string

	File string

	Type string
}

//New Creates a new FilePath struct from a given path

func New(path string) *FilePath {

	var file FilePath

	var dirs string

	//Ensures PathSeparator are for the OS

	p := filepath.Clean(path)

	//Get the volume name from the path

	file.Volume = filepath.VolumeName(p) + string(os.PathSeparator)

	p, _ = filepath.Rel(file.Volume, p)

	// Checks if path has a user

	file.User = hasUser(p)

	// Removes user from viable path

	p, _ = filepath.Rel(file.Userf(), p)

	// Checks if path has a file or not

	if isFile(path) {

		//Splits the file from the dirs

		file.Type = "File"

		dirs, file.File = filepath.Split(p)

	} else {

		//Sets path to have no file

		file.Type = "Dir"

		dirs = p

	}

	list := strings.Split(dirs, string(os.PathSeparator))

	file.Head = list[0]

	switch len(list) {

	case 1:

		file.Tail = []string{}

	case 2:

		file.Tail = []string{list[1]}

	default:

		file.Tail = list[1 : len(list)-1]

	}

	return &file

}

func hasUser(path string) string {

	list := strings.Split(path, string(os.PathSeparator))

	if list[0] == "Users" {

		return list[1]

	}

	return ""

}

func isFile(path string) bool {

	fi, err := os.Stat(path)

	if err != nil {

		fmt.Println(err)

	}

	if fi.Mode().IsRegular() {

		return true

	}

	return false

}

// Tailf return the tail after the head to the directory that houses the file

func (P *FilePath) Tailf() string {

	return filepath.Join(P.Tail...)

}

// Userf returns the user joined with Users

func (P *FilePath) Userf() string {

	if P.User != "" {

		return filepath.Join("Users", P.User)

	}

	return ""

}

//Join returns the user requested values of the path

func (P *FilePath) Join(values ...string) string {

	var base string

	if len(values) == 0 {

		return filepath.Join(P.Volume, P.Userf(), P.Head, P.Tailf(), P.File)

	}

	for i := 0; i < len(values); i++ {

		switch values[i] {

		case "head":

			filepath.Join(base, P.Head)

		case "volume":

			filepath.Join(base, P.Volume)

		case "user":

			filepath.Join(base, P.Userf())

		case "tail":

			filepath.Join(base, P.Tailf())

		case "file":

			filepath.Join(base, P.File)

		}

	}

	return base

}

func (P *FilePath) ZipPath() string {

	if P.Type == "Dir" {

		return filepath.ToSlash(filepath.Join(P.User, P.Head, P.Tailf(), P.File)) + "/"

	}

	return filepath.ToSlash(filepath.Join(P.User, P.Head, P.Tailf(), P.File))

}

func (P *FilePath) String() string {

	return fmt.Sprintf("Volume: %s\nUser: %s\nHead: %s\nTail: %s\nFile: %s\nType: %s", P.Volume, P.User, P.Head, P.Tail, P.File, P.Type)

}
