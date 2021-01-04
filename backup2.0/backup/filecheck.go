package backup

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

func FileCheck(file string, Use_Exclusions, Use_Inclusions bool, Included, Excluded, File_Types []string) bool {
	if InvalidExtension(File_Types, file) && Use_Exclusions {
		return false
	}
	if !Use_Exclusions && !Use_Inclusions {
		return true
	} else if !Use_Exclusions && Use_Inclusions {
		// Only backup if included
		ok := IsSlice(Included, file)
		if !ok {
			return false
		}
	} else if Use_Exclusions && !Use_Inclusions {
		// Only backup if not excluded
		if IsSlice(Excluded, file) {
			return false
		} else {
			return true
		}
	} else if Use_Exclusions && Use_Inclusions {
		//Backup if not exluded unless explicitly included
		ok := IsSlice(Included, file)
		exclude := Is(Excluded, file)

		if !exclude && ok {
			return true
		} else if exclude && ok {
			return true
		} else if !ok && !exclude {
			return true
		} else {
			return false
		}
	} else {
		panic(errors.New(fmt.Sprintf("Error: The combinantion of %t,%t is not possible", Use_Exclusions, Use_Inclusions)))
	}
	return false
}

func Is(slice []string, value string) bool {
	for _, elem := range slice {
		if elem == value {
			return true
		} else if strings.Contains(value, elem) {
			return true
		}
	}
	return false
}

//Checks if values in sliceA are in the file path at all
func IsSlice(sliceA []string, file string) bool {
	files := strings.Split(file, "\\")
	file = strings.Join(files[3:], "\\")
	for _, elemA := range sliceA {
		if strings.Contains(file, elemA) {
			return true //,strings.Replace(elemA,elemB,"",1)
		} else if strings.Contains(elemA, file) {
			return true
		}
	}
	return false //,""
}

func InvalidExtension(extensions []string, file string) bool {
	for _, ext := range extensions {
		if ext == filepath.Ext(file) {
			return true
		}
	}
	return false
}
