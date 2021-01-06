package backup

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	logPath    = "C:\\TEST\\OUTLOG"
	file1      = createdst(logPath, ".txt")
	InfoLogger = log.New(file1, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
)

//FileCheck validates files using Included and Excluded. Use_Exclusions/Use_Inclusions define which lists are being used,
//File_Types is the list of excluded files types
func FileCheck(file string, UseExclusions, UseInclusions bool, Include, Exclude, FileTypes []string) bool {
	excluded := false
	included := true

	if UseExclusions {
		excluded = IsSlice(Exclude, file)
	}
	if UseInclusions {
		included = IsSlice(Include, file)
	}
	invalidExt := InvalidExtension(FileTypes, file)

	InfoLogger.Println("Name:", file)
	InfoLogger.Println("excluded:", excluded)
	InfoLogger.Println("included:", included)
	InfoLogger.Println("invalid Ext:", invalidExt)

	if excluded && !included || invalidExt {
		return false
	}
	return true

}

//Checks if values in p are in the file path at all
func IsSlice(p []string, file string) bool {
	pattern := strings.ToLower(strings.Join(p, "|"))
	file = strings.ToLower(file)
	file = strings.Replace(file, "\\", "/", -1)
	result, _ := regexp.MatchString(pattern, file)
	if result {
		return true
	}

	for _, a := range p {
		a = strings.Replace(a, "./", "", -1)
		a = strings.Replace(a, "/.", "", -1)
		dir, _ := filepath.Split(file)
		file = strings.ToLower(dir + a)
		if results, _ := regexp.MatchString(pattern, file); results {
			InfoLogger.Println(file)
			if _, err := os.Stat(file); err == nil {
				return true
			} else if os.IsNotExist(err) {
				continue
			}
		}
	}
	return false
}

func InvalidExtension(extensions []string, file string) bool {
	pattern := strings.ToLower(strings.Join(extensions, "|"))
	file = strings.ToLower(file)
	result, _ := regexp.MatchString(pattern, file)
	if result {
		return true
	}
	return false
}
