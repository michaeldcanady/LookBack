package structure

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	//"github.com/BurntSushi/toml"

	"github.com/michaeldcanady/Project01/backup2.0/servicenow"
)

type settings struct {
	Use_Exclusions  bool   `toml: "Use_Exclusions"`
	Use_Inclusions  bool   `toml: "Use_Inclusions"`
	Email_Extension string `toml: "Email_Extension"`
	Network_Path    string `toml: "Network_Path"`
}

type adsettings struct {
	Use_Ecryption bool   `toml: "Use_Encryption"`
	domain        string `toml: "Domain"`
}

type tktsystem struct {
	Provider string `toml: "Provider"`
	URL      string `toml: "URL"`
}

type exclusion struct {
	General_Exclusions   []string `toml: "General_Exclusions"`
	Profile_Exclusions   []string `toml: "Profile_Exclusions"`
	File_Type_Exclusions []string `toml: "File_Type_Exclusions"`
}

type inclusion struct {
	General_Inclusions []string `toml: "General_Inclusions"`
	Profile_Inclusions []string `toml: "Profile_Inclusions"`
}

type Config struct {
	Settings          settings
	Tktsystem         tktsystem
	Exclusions        exclusion
	Inclusions        inclusion
	Advanced_Settings adsettings
}

// struct used to store all data for a Backup
type Backup struct {
	Technician string
	Password   string
	Client     servicenow.Back
	CSNumber   string
	Task       string
	Source     []string
	DestType   string
	Dest       string
}

// struct for storing User data
type User struct {
	Path string
	Size int64
}

// struct for storing File data
type File struct {
	Name     string
	Filepath string
	Hash     string
}

// function for creating new File struct
func newFile(path string) File {
	File := File{Filepath: path}
	//File.hash,_ = HashFile(path)
	File.Name = filepath.Base(path)
	return File
}

// function for creating new User struct
func NewUser(path string) User {
	var u User
	u.Path = path
	u.Size = DirSize(path)
	return u
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
