package file

import (
	"fmt"
	"os"
	"runtime"

	"github.com/michaeldcanady/LookBack/backup2.0/conversion"
	"github.com/michaeldcanady/LookBack/backup2.0/hash"
	"github.com/michaeldcanady/LookBack/backup2.0/path"
)

var UNIT int64

func init() {
	if runtime.GOOS == "windows" {
		UNIT = 1024
	} else {
		UNIT = 1000
	}
}

type File struct {
	Path *filestruct.FilePath
	File *os.FileInfo
	Hash string
}

func New(path string, file *os.FileInfo) *File {
	h, _ := hash.HashFile(path, UNIT)
	return &File{
		Path: filestruct.New(path),
		File: file,
		Hash: h,
	}
}

//Methods for Path
//Returns a string version of the path
func (F *File) PathStr() string {
	return fmt.Sprintf("%v", F.Path.Join())
}

//Returns the path volume
func (F *File) PathVol() string {
	return F.Path.Volume
}

//Returns the path head
func (F *File) PathHead() string {
	return F.Path.Head
}

//Returns the path tail
func (F *File) PathTail() string {
	return F.Path.Tailf()
}

//Returns the path user
func (F *File) PathUserf() string {
	return F.Path.Userf()
}
func (F *File) PathUserp() string {
	return F.Path.User
}
func (F *File) PathFile() string {
	return F.Path.File
}

//Returns the bytes of a file
func (F *File) RawSize() int64 {
	return (*F.File).Size()
}

//Returns the size formatted by OS
func (F *File) Size() string {
	return conversion.ByteCountSI(F.RawSize(), UNIT, 0)
}
