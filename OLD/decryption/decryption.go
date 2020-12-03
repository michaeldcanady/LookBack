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
  "encoding/hex"
)

func main(){
  var key string
  var ticket string
  fmt.Printf("Enter key:")
  fmt.Scan(&key)
  fmt.Printf("Enter Ticket Number:")
  fmt.Scan(&ticket)
  key1,_ := hex.DecodeString(key)

  path := filepath.Join("C:\\",ticket)
  os.MkdirAll("C:\\Test",os.ModePerm)
  recuse(path,"C:\\Test",key1,ticket)
}

func recuse(src,dst string,key []byte,ticket string){
  files,err := filepath.Glob(path.Join(src,"*"))
  if err != nil{
    fmt.Println("Getting files Error:",err)
  }
  for _,file := range files{
    fi,_ := os.Stat(file)
    switch mode := fi.Mode(); {
      case mode.IsDir():
        //if file == filepath.Join("C:\\",ticket,"dmcanady\\UserData"){
        //  continue
        //}
        fmt.Println(file)
        recuse(file,dst,key,ticket)
      case mode.IsRegular():
        var name = file
        for _,i := range []string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}{
          //fmt.Println(i+":\\Users\\")
          name = strings.ReplaceAll(name,i+":\\"+ticket+"\\","")
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
