package main

import(
  "github.com/blend/go-sdk/crypto"
  "os"
  "path/filepath"
  "io"
  "strings"
)

func main(){
  key := []byte{109, 39, 108, 91, 114, 113, 246, 167, 93, 168, 74, 189, 26, 104, 7, 67, 121, 232, 57, 221, 15, 105, 71, 67, 233, 70, 253, 234, 67, 149, 44, 97}
  path := "C:\\CS0001"
  recuse(path,"C:\\Test",key)
}

func recuse(path,dst string,key []byte){
  files,err := filepath.Glob(path)
  for _,file := range files{
    fi,_ := os.Stat(file)
    switch mode := fi.Mode(); {
      case mode.IsDir():
        recuse(file,dst,key)
      case mode.IsRegular():
        var name = file
        for _,i := range []string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}{
          //fmt.Println(i+":\\Users\\")
          name = strings.ReplaceAll(name,i+":\\Users\\","")
        }

        dst := filepath.Join(dst,name)
        dir,_ := filepath.Split(dst)

        source,_ := crypto.NewStreamDecrypter(key,file)
        io.Copy(dst,source)
      }
  }
}
