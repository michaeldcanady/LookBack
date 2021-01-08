package structure

import (
	"os"
	"path/filepath"

	//"github.com/BurntSushi/toml"

	"github.com/michaeldcanady/LookBack/backup2.0/servicenow"
)

type settings struct {
	Use_Exclusions     bool   `toml: "Use_Exclusions"`
	Use_Inclusions     bool   `toml: "Use_Inclusions"`
	Email_Extension    string `toml: "Email_Extension"`
	Network_Path       string `toml: "Network_Path"`
	NetworkFolderPath  string `toml: "NetworkFolderPath"`
	WinServerBackupMax int64  `toml: "WinServerBackupMax"`
	MacServerBackupMax int64  `toml: "MacServerBackupMax"`
	NetworkDriveName   string `toml: "NetworkDriveName"`
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
	Source     []User
	DestType   string
	Dest       string
}

// struct for storing User data
type User struct {
	Path     string
	Size     int64
	RootDirs map[string]int64
	Files    map[string]os.FileInfo
}

// function for creating new User struct
func NewUser(path string) User {
	var u User
	rootDirs, files, size := UserSize(path)
	u.Path = path
	u.Size = size
	// Use this sizing to get the rootdirs
	u.RootDirs = rootDirs
	u.Files = files
	return u
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

func UserSize(homedir string) (map[string]int64, map[string]os.FileInfo, int64) {

	Use_Exclusions := Conf.Settings.Use_Exclusions
	Use_Inclusions := Conf.Settings.Use_Inclusions
	Excluded := Conf.Exclusions.General_Exclusions
	ExcludedFiles := Conf.Exclusions.File_Type_Exclusions
	Included := Conf.Inclusions.General_Inclusions

	rootDirs := make(map[string]int64)
	var total int64
	files := make(map[string]os.FileInfo)
	subDirectories, _ := filepath.Glob(homedir + "/**")
	for _, subDirectory := range subDirectories {
		if !FileCheck(subDirectory, Use_Exclusions, Use_Inclusions, Included, Excluded, ExcludedFiles) {
		} else {
			size, err := DirSize(subDirectory, &files)
			if err != nil {
				panic(err)
			}
			rootDirs[subDirectory] = size
			total += size
		}
	}
	return rootDirs, files, total
}

//HAVE IT FACTOR IN FILES THAT NEED TO BE SKIPPED
//Gets size of specified directory
func DirSize(path string, files *map[string]os.FileInfo) (int64, error) {

	Use_Exclusions := Conf.Settings.Use_Exclusions
	Use_Inclusions := Conf.Settings.Use_Inclusions
	Excluded := Conf.Exclusions.General_Exclusions
	ExcludedFiles := Conf.Exclusions.File_Type_Exclusions
	Included := Conf.Inclusions.General_Inclusions

	var size int64
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !FileCheck(filePath, Use_Exclusions, Use_Inclusions, Included, Excluded, ExcludedFiles) {
		} else if !info.IsDir() {
			size += info.Size()
			(*files)[filePath] = info
		}
		return err
	})
	return size, err
}
