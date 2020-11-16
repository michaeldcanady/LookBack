package main

import(
  "fmt"
  "path/filepath"
  "path"
  "os"
  "io"
  "strings"
  "github.com/BurntSushi/toml"
  "errors"
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

var(
  conf Config
)

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
              GetFiles(file,true,Settings,Inclusions,Exclusions)
            }
          }
      case mode.IsRegular():
        // Adds filepath to file channel
        read <- file
    }
  }
}

func Gatherer(src []string,read chan string,hashSlice *[]file,wg *sync.WaitGroup){
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
      GetFiles(src,true,conf.Settings,conf.Inclusions,conf.Exclusions)
    }
  }
}

//Checks if a file/directory needs to be excluded
//if Is(exclude,file){
//  ok,path := IsSlice(include,file)
//  if ok{
//    fmt.Println(filepath.Join(file,path))
//  }else{
//    continue
//  }
//}else{
//}
//return nil

func init(){

  if _, err := toml.DecodeFile("C:\\go\\src\\github.com\\michaeldcanady\\Project01\\Config\\Settings.toml", &conf); err != nil {
    panic(err)
  }
}
