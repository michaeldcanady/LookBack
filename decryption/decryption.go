package main

import(
  "github.com/blend/go-sdk/crypto"
  "os"
  "path/filepath"
  "io"
  "strings"
  "bytes"
  "fmt"
  "path"
  "io/ioutil"
)

func main(){
  key1 := []byte{100, 146, 51, 146, 145, 205, 138, 24, 155, 160, 89, 143, 251, 113, 208, 117, 152, 104, 145, 236, 181, 15, 189, 143, 187, 31, 78, 202, 85, 206, 25, 26}
  //key := []byte{109, 39, 108, 91, 114, 113, 246, 167, 93, 168, 74, 189, 26, 104, 7, 67, 121, 232, 57, 221, 15, 105, 71, 67, 233, 70, 253, 234, 67, 149, 44, 97}

  path := "C:\\CS004\\micha"
  os.MkdirAll("C:\\Test",os.ModePerm)
  recuse(path,"C:\\Test",key1)
}

func recuse(src,dst string,key []byte){
  files,err := filepath.Glob(path.Join(src,"*"))
  if err != nil{
    fmt.Println("Getting files Error:",err)
  }
  for _,file := range files{
    fi,_ := os.Stat(file)
    switch mode := fi.Mode(); {
      case mode.IsDir():
        fmt.Println(file)
        recuse(file,dst,key)
      case mode.IsRegular():
        var name = file
        for _,i := range []string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}{
          //fmt.Println(i+":\\Users\\")
          name = strings.ReplaceAll(name,i+":\\CS004\\","")
        }

        dst := filepath.Join(dst,name)
        dir,_ := filepath.Split(dst)
        fmt.Println(dir)
        os.MkdirAll(dir,os.ModePerm)
        destination, err := os.Create(dst)
        if err != nil {
          fmt.Println("CREATION ERROR",err)
          panic(err)
        }
        defer destination.Close()
        f,_ := ioutil.ReadFile(file)
        source,err := crypto.Decrypt(key,f)
        if err != nil{
          panic(err)
        }
        //fmt.Println(string(source))
        r := bytes.NewReader(source)
        _,err = io.Copy(destination,r)
        if err != nil{
          fmt.Println("Copy Error",err)
        }
        //fmt.Println(total)
      }
  }
}
