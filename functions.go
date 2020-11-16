package main

import(
  "sync"
  "path/filepath"
  "path"
  "fmt"
  "os"
  "io"
  "strings"
  "errors"
)

// empty struct (0 bytes)
type void struct{}

var(
  conf Config
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

func Is(slice []string, value string)bool{
  for _,elem := range slice{
    if elem == value{
      return true
    }else if strings.Contains(value,elem){
      return true
    }
  }
  return false
}

func IsSlice(sliceA []string , file string)bool{
  files := strings.Split(file,"\\")
  file = strings.Join(files[3:],"\\")
  for _,elemA := range sliceA{
    if strings.Contains(file,elemA){
      return true//,strings.Replace(elemA,elemB,"",1)
    }else if strings.Contains(elemA,file){
      return true
    }
  }
  return false//,""
}

type settings struct{
  Use_Exclusions bool `toml: "Use_Exclusions"`
  Use_Inclusions bool `toml: "Use_Inclusions"`
}

type adsettings struct{
  use_ecryption bool `toml: "Use_Encryption"`
  domain string `toml: "Domain"`
}

type exclusion struct{
  General_Exclusions []string `toml: "General_Exclusions"`
  Profile_Exclusions []string `toml: "Profile_Exclusions"`
}

type inclusion struct{
  General_Inclusions []string `toml: "General_Inclusions"`
  Profile_Inclusions []string `toml: "Profile_Inclusions"`
}

type Config struct{
  Settings settings
  Exclusions exclusion
  Inclusions inclusion
  Advanced_Settings adsettings
}

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

func GetFiles(src string, read chan string, hashSlice *[]file, recusive bool,Settings settings,Inclusions inclusion,Exclusions exclusion){
  Use_Exclusions := Settings.Use_Exclusions
  Use_Inclusions := Settings.Use_Inclusions
  Excluded := Exclusions.General_Exclusions
  Included := Inclusions.General_Inclusions

  files,_ := filepath.Glob(path.Join(src,"*"))

  for _,file := range files{
    if !Use_Exclusions && !Use_Inclusions{
      //Backup All Files
    }else if !Use_Exclusions && Use_Inclusions{
      // Only backup if included
      ok := IsSlice(Included,file)
      if !ok{
        continue
      }else{
      }
    }else if Use_Exclusions && !Use_Inclusions{
      // Only backup if not excluded
      if Is(Excluded,file){
        continue
      }else{
      }
    }else if Use_Exclusions && Use_Inclusions{
      //Backup if not exluded unless explicitly included
      ok := IsSlice(Included,file)
      if Is(Excluded,file) && ok{

      }else if Is(Excluded,file){
        continue
      }
    }else{
      panic(errors.New(fmt.Sprintf("Error: The combinantion of %t,%t is not possible",Settings.Use_Exclusions,Settings.Use_Inclusions)))
    }

    // Gets file stats
      fi, err := os.Stat(file); if os.IsNotExist(err) {
        fmt.Println("No exist",err)
      }else if err != nil {
        fmt.Println("Stat",err)
     }
    switch mode := fi.Mode(); {
      case mode.IsDir():
          empty,err := IsDirEmpty(file); if err != nil{
            fmt.Println("DirEmpty Error:",err)
          }
          if recusive == false{
            read <- file
            continue
          }else{
            if empty{
              fmt.Println("Empty:",file)
            }else{
              GetFiles(src,read,hashSlice,true,conf.Settings,conf.Inclusions,conf.Exclusions)
            }
          }
      case mode.IsRegular():
        // Adds filepath to file channel
        read <- file
    }
  }
}

func Gatherer(srcs []string,read chan string,hashSlice *[]file,wg *sync.WaitGroup){
  defer wg.Done()
  defer close(read)

  for _,src := range srcs{
    tempsrc := src
    dirs := strings.Split(tempsrc,PATHSEPARATOR)
    fmt.Println(dirs)
    tempsrc = dirs[len(dirs)-1]
    if _, err := os.Stat(src); os.IsNotExist(err) {
      continue
    }else{
      files,_ := filepath.Glob(path.Join(src,"*"))
      // Creates all files in user folder, directories are empty
      for _,file := range files{
        read <- file
      }
      GetFiles(src,read,hashSlice,true,conf.Settings,conf.Inclusions,conf.Exclusions)
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
      fmt.Println(f)
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
