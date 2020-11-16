package main

import(
  "path/filepath"
  //"github.com/BurntSushi/toml"
)

// struct used to store all data for a backup
type backup struct{
  Technician string
  CSNumber string
  Task string
  Source []string
  DestType string
  Dest string

}

// struct for storing user data
type user struct{
  path string
  size int64

}

// struct for storing file data
type file struct {
  name string
  filepath string
  hash string
}

// function for creating new file struct
func newFile(path string) file{
  file := file{filepath: path}
  file.hash,_ = HashFile(path)
  file.name = filepath.Base(path)
  return file
}

// function for creating new user struct
func NewUser(path string)user{
  var u user
  u.path = path
  u.size = DirSize(path)
  return u
}
