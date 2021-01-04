package main

import(
  "fmt"
  "path/filepath"
)

func checkerr(err error){
  if err != nil {
    panic(err)
  }
}

func main(){
  volumes,err := filepath.Glob("/Volumes/**")
  checkerr(err)
  fmt.Println(volumes)
}
