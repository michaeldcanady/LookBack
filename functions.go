package main

import(
  "sync"
  "path/filepath"
  "path"
  "fmt"
  "os"
  "io"
  "strings"
)

// empty struct (0 bytes)
type void struct{}

var(

)

func Checkwhitelist(path string)bool{
  for _,files := range whitelist{
    dir,_ := filepath.Split(path)
    if filepath.Join(BASE,USER,files) == dir{
      return true
    }else{
      continue
    }
  }
  return false
}

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

func GetFiles(src string,read chan string,hashSlice *[]file,recusive bool){
  fmt.Println(src)
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
              read <- s
              continue
            }
            if empty{
              read <- s
            }else{
              fmt.Println("RECUSIVE")
              GetFiles(s,read,hashSlice,true)
            }
        case mode.IsRegular():
          if s != "C:\\Users\\dmcanady\\NTUSER.DAT" && s != "C:\\Users\\dmcanady\\ntuser.dat.LOG1" && s != "C:\\Users\\dmcanady\\ntuser.dat.LOG2"{
            // Add hash to hash channel
            *hashSlice = append(*hashSlice,newFile(s))
            // Adds filepath to file channel
            read <- s
          }
      }
    }
  }

func Gather(srcs []string,read chan string,hashSlice *[]file,wg *sync.WaitGroup){
  defer wg.Done()
  defer close(read)
  for _,src := range srcs{
    tempsrc := src
    dirs := strings.Split(tempsrc,PATHSEPARATOR)
    fmt.Println(dirs)
    tempsrc = dirs[len(dirs)-1]
    if _, err := os.Stat(src); os.IsNotExist(err) {
      continue
    }else if tempsrc == "Favorites" || tempsrc == "Desktop" || tempsrc == "Documents" || tempsrc == "Contacts" || tempsrc == "Music" || tempsrc == "Pictures" || tempsrc == "Videos"{
      fmt.Println(tempsrc,"RECUSIVE")
      GetFiles(src,read,hashSlice,true)
    }else{
      fmt.Println(tempsrc,"NOT RECUSIVE")
      GetFiles(src,read,hashSlice,false)
    }
  }
}

func copy(dst string, read chan string,wg *sync.WaitGroup,Newfile *[]file){
  defer wg.Done()
  for{
    f,ok := <- read
    if ok == false{
      break
    }else{
      var name = f
      for _,i := range []string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}{
        //fmt.Println(i+":\\Users\\")
        name = strings.ReplaceAll(name,i+":\\Users\\","")
      }
      dst := filepath.Join(dst,name)
      dir,_ := filepath.Split(dst)
      sourceFileStat, err := os.Stat(f)
      if err != nil {
        panic(err)
      }

      if sourceFileStat.Mode().IsDir() {
        os.Mkdir(dst,os.ModePerm)
        continue
        //panic(fmt.Errorf("%s is not a regular file", f))
      }

      source, err := os.Open(f)
      if err != nil {
        panic(err)
      }
      defer source.Close()
      os.MkdirAll(dir,os.ModePerm)
      destination, err := os.Create(dst)
      if err != nil {
        fmt.Println("CREATION ERROR",err)
        panic(err)
      }
      defer destination.Close()
      _, err = io.Copy(destination, source)
      *Newfile = append(*Newfile,newFile(dst))
    }
  }
}

func containsfile(s []file, e string)(bool){
  for _, a := range s {
    if a.hash == e {
      return true
    }
  }
  return false
}

func VerifyFile(NewFiles,Orignialfiles []file)(float32,[]file,[]file){
  var Success []file
  var Failed []file
  fmt.Println("Comparing copied files to originial")
  for _,Ofile := range Orignialfiles{
    if containsfile(NewFiles,Ofile.hash){
      Success = append(Success,Ofile)
    }else{
      Failed = append(Failed,Ofile)
    }
  }
  newlen := float32(len(Success))
  oldlen := float32(len(Orignialfiles))
  fmt.Printf("Successfully merged %v\n",newlen)
  fmt.Printf("Orginial length %v\n",oldlen)
  return ((newlen/oldlen)*100),Success,Failed
}

func SliceSizef(slice []file) int64{
  var totalSize int64
  for _,files := range slice{
    file,err := os.Stat(files.filepath); if err != nil{
      panic(err)
    }else{
      totalSize += file.Size()
    }
  }
  return totalSize
}
