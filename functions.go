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
  "github.com/michaeldcanady/Project01/fileEncryption"
  "os/exec"
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

func GetFiles(src string, read chan string, hashSlice *[]file, recusive bool,Settings settings,Inclusions inclusion,Exclusions exclusion){
  Use_Exclusions := Settings.Use_Exclusions
  Use_Inclusions := Settings.Use_Inclusions
  Excluded := Exclusions.General_Exclusions
  ExcludedFiles := Exclusions.File_Type_Exclusions
  Included := Inclusions.General_Inclusions
  files,err := filepath.Glob(path.Join(src,"*"))
  if err!= nil{
    fmt.Println("Glob error",err)
  }
  // Logic to see if files match requirements
  for _,file := range files{
    if !FileCheck(file, Use_Exclusions, Use_Inclusions, Included, Excluded, ExcludedFiles){
      continue
    }else{
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
              read <- file
            }else{
              GetFiles(file,read,hashSlice,true,conf.Settings,conf.Inclusions,conf.Exclusions)
            }
          }
      case mode.IsRegular():
        // Hash for verification
        *hashSlice = append(*hashSlice,newFile(file))
        read <- file
      }
    }
  }
}

func GetInstalledPrograms(path string)error{
  f,err := os.Create(filepath.Join(path,"InstalledPrograms.txt"))
  if err != nil{
    return err
  }
  commandString := "wmic product get name"
  output, err := exec.Command("Powershell", "-Command", commandString ).CombinedOutput()
  if err != nil{
    return err
  }
  _,err = f.WriteString(string(output))
  if err != nil{
    return err
  }
  return nil

}

func InvalidExtension(extensions []string, file string)bool{
  for _,ext := range extensions{
    if ext == filepath.Ext(file){
      return true
    }
  }
  return false
}

func FileCheck(file string,Use_Exclusions,Use_Inclusions bool, Included, Excluded,File_Types []string)bool{
    if InvalidExtension(File_Types,file) && Use_Exclusions{
      return false
    }
    if !Use_Exclusions && !Use_Inclusions{
      return true
    }else if !Use_Exclusions && Use_Inclusions{
      // Only backup if included
      ok := IsSlice(Included,file)
      if !ok{
        return false
      }
    }else if Use_Exclusions && !Use_Inclusions{
      // Only backup if not excluded
      if IsSlice(Excluded,file){
        return false
      }else{
        return true
      }
    }else if Use_Exclusions && Use_Inclusions{
      //Backup if not exluded unless explicitly included
      ok := IsSlice(Included,file)
      exclude := Is(Excluded,file)

      if !exclude && ok{
        return true
      }else if exclude && ok{
        return true
      }else if !ok && !exclude{
        return true
      }else{
        return false
      }
    }else{
      panic(errors.New(fmt.Sprintf("Error: The combinantion of %t,%t is not possible",Use_Exclusions,Use_Inclusions)))
    }
  return false
}

func Gatherer(srcs []string,read chan string,hashSlice *[]file,wg *sync.WaitGroup){
  defer wg.Done()
  defer close(read)

  for _,src := range srcs{
    tempsrc := src
    dirs := strings.Split(tempsrc,PATHSEPARATOR)
    tempsrc = dirs[len(dirs)-1]
    if _, err := os.Stat(src); os.IsNotExist(err) {
      continue
    }else{
      files,_ := filepath.Glob(path.Join(src,"*"))
      // Creates all files in user folder, directories are empty
      for _,file := range files{
        if IsSlice(conf.Exclusions.General_Exclusions,file){
          continue
        }else{
          fi, err := os.Stat(file); if os.IsNotExist(err) {
            fmt.Println("No exist",err)
          }else if err != nil {
            fmt.Println("Stat",err)
          }
          switch mode := fi.Mode(); {
            case mode.IsDir():
              read <- file
            case mode.IsRegular():
              // Hash for verification
              *hashSlice = append(*hashSlice,newFile(file))
              read <- file
          }
        }
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
        panic(fmt.Sprintf("bad error: %s",err))
      }

      if sourceFileStat.Mode().IsDir() {
        os.Mkdir(dst,os.ModePerm)
        continue
        //panic(fmt.Errorf("%s is not a regular file", f))
      }

      source, err := os.Open(f)
      if err != nil {
        panic(fmt.Sprintf("dst copy: %s",err))
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
      if err != nil{
        fmt.Println("Copy Error",err)
      }
      //Add check if it is a file or a folder, if a folder, do not hash.
      *Newfile = append(*Newfile,newFile(dst))
      if conf.Advanced_Settings.Use_Ecryption == true && filepath.Ext(dst) == ".txt"{






        //panic("Encrypting file")
        fmt.Println(dst)
        err = Encrypt_file(dst)
        if err != nil{
          panic(fmt.Sprintf("Encryption error %s",err))
        }
      }
    }
  }
}

// returns assumed to be encrypted bits needs more more
func Encrypt_file(file string)error{
  key := "C:\\go\\src\\github.com\\michaeldcanady\\Project01\\.ssh\\id_rsa.public"
  publicKey1, err := encryption.RetrievePublic(key); if err != nil{
    return errors.New(fmt.Sprintf("Retrieval err: %s", err))
  }
  fileBytes := encryption.Encrypt(file,publicKey1)
  fmt.Println(fileBytes)
  source, err := os.Open(file)
  if err != nil {
    return errors.New(fmt.Sprintf("dst copy: %s",err))
  }
  _, err = source.WriteString(string(fileBytes))
  if err != nil{
    return errors.New(fmt.Sprintf("Writing error: %s",err))
  }
  source.Close()
  return nil
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
