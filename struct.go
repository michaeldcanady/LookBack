package main

import(
  "path/filepath"
)


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
