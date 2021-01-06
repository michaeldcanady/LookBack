package copy

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/michaeldcanady/Project01/backup2.0/file"
	"github.com/michaeldcanady/Project01/backup2.0/hash"
)

var (
	WarningLogger   *log.Logger
	InfoLogger      *log.Logger
	ErrorLogger     *log.Logger
	HashErrorLogger *log.Logger
)

// Copy to copy file from src to dstbase (root file)
func Copy(dstbase string, file *file.File, UNIT int64, backup bool) {

	src := fmt.Sprintf("%v", file.Path.Join())
	var dst string

	if !backup {
		dst = filepath.Join(dstbase, file.PathTail(), file.PathFile())
	} else {
		dst = filepath.Join(dstbase, file.PathUserp(), file.PathHead(), file.PathTail(), file.PathFile())
	}

	logPath := filepath.Join(dstbase, "logs", "errorLog.txt")
	logPath1 := filepath.Join(dstbase, "logs", "Hasherror.txt")

	file1 := createdst(logPath, false)
	file2 := createdst(logPath1, false)

	InfoLogger = log.New(file1, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file1, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file1, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	HashErrorLogger = log.New(file2, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	source, err := os.Open(src)
	if err != nil {
		command := fmt.Sprintf(`echo F | xcopy "%s" "%s" /Y`, src, dst)
		_, err := exec.Command("Powershell", "-Command", command).CombinedOutput()
		if err != nil {
			check(fmt.Errorf("Coping %s failed", src))
		}
	} else {
		defer source.Close()
		destination := createdst(dst, false)
		b1 := make([]byte, UNIT)
		_, err = io.CopyBuffer(destination, source, b1)
		check(err)
		destination.Close()
	}
	h, err := hash.HashFile(dst, UNIT)
	if err != nil {
		missing := filepath.Join(dstbase, file.PathUserp()+"Files", "Missing.txt")
		file2 := createdst(missing, false)
		_, _ = file2.WriteString(src + "\n")
	} else {
		if hash.CompareHash(file.Hash, h) == false {
			missing := filepath.Join(dstbase, file.PathUserp()+"Files", "Missing.txt")
			file2 := createdst(missing, false)
			_, _ = file2.WriteString(src + "\n")
		}
	}
}

func check(err error) {
	if err != nil {
		ErrorLogger.Println(err)
	}
}

func createdst(dst string, encrypt bool) *os.File {
	ext := ""
	if encrypt {
		ext = ".temp"
	}
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		head, _ := filepath.Split(dst)
		os.MkdirAll(head, 0700)
	}
	destination, err := os.OpenFile(dst+ext, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	check(err)
	return destination
}
