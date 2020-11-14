package main

import(
  "path/filepath"
)

type backup struct{
  Technician string
  CSNumber string
  Task string
  Source []string
  DestType string
  Dest string

}

func NewUser(path string)user{
  var u user
  u.path = path
  u.size = DirSize(path)
  return u
}

type user struct{
  path string
  size int64

}

type file struct {
  name string
  filepath string
  hash string
}

func newFile(path string) file{
  file := file{filepath: path}
  file.hash,_ = HashFile(path)
  file.name = filepath.Base(path)
  return file
}
