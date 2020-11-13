package main

import(
  "fmt"
  "path/filepath"
  "path"
  "os"
  "io"
)

func IsDirEmpty(name string) (bool, error) {
        f, err := os.Open(name)
        if err != nil {
                return false, err
        }
        defer f.Close()

        // read in ONLY one file
        _, err = f.Readdir(1)

        // and if the file is EOF... well, the dir is empty.
        if err == io.EOF {
                return true, nil
        }
        return false, err
}

func GetFiles(src string, recusive bool, exclude... string){
  fmt.Println(src)
  for _,ex := range exclude{
    if src == ex{
      break
    }
    file, err := filepath.Glob(path.Join(src,"*"))
    if err != nil {
      fmt.Println(err)
    }
    for _, s := range file{
      fi, err := os.Stat(s)
      if err != nil {
        return
      }
    // checks if file or directory
      switch mode := fi.Mode(); {
        case mode.IsDir():
            empty,_ := IsDirEmpty(s)
            if recusive == false{
              fmt.Println("NOT RECUSIVE")
              fmt.Println(s)
              continue
            }
            if empty{
              fmt.Println(s)
            }else{
              fmt.Println("RECUSIVE")
              GetFiles(s,true,exclude)
            }
        case mode.IsRegular():
          if s != "C:\\Users\\dmcanady\\NTUSER.DAT" && s != "C:\\Users\\dmcanady\\ntuser.dat.LOG1" && s != "C:\\Users\\dmcanady\\ntuser.dat.LOG2"{
            // Add hash to hash channel
            //*hashSlice = append(*hashSlice,newFile(s))
            // Adds filepath to file channel
            fmt.Println(s)
          }
      }
    }
  }
}

func main(){
  src := filepath.Clean("C:\\Users\\dmcanady")
  GetFiles(src,true,"C:\\Users\\AppData")
}
