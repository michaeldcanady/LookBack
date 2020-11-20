// +build windows darvin linux dev

package main

import(
  "fmt"
  "sync"
  "time"
  "os"
  "path/filepath"
  "log"
  //"github.com/michaeldcanady/SliceTools"
  "github.com/AlecAivazis/survey"
  "github.com/blend/go-sdk/crypto"
  "encoding/hex"
)

var (
  //Different Loggers
  WarningLogger *log.Logger
  InfoLogger    *log.Logger
  ErrorLogger   *log.Logger
  MissedLogger  *log.Logger

  //Channels
  read = make(chan string)

  //Users specific
  USER string
  BASE string

  //Hash slices
  Orignialhash []file
  Newhash []file

  //WaitGroup
  wg sync.WaitGroup
)

func Setlog(dst string){
  f, err := os.OpenFile(filepath.Join(dst,"UserData","MissingFiles.log"), os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  if err != nil {
    log.Fatalf("error opening file: %v", err)
  }
  log.SetOutput(f)
  //InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
  //WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
  //ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
  MissedLogger = log.New(f,"",log.Ltime)

}

func InLineCopy(binfo *backup) {
  number := binfo.CSNumber
  srcs := []string{}
  for _,src := range binfo.Source{
    srcs = append(srcs,src)
  }
  dst := filepath.Join(binfo.Dest,number)
  os.Mkdir(filepath.Join(dst,USER),os.ModePerm)
  Header()
  wg.Add(1)
  start := time.Now()
  go Gatherer(srcs,read,&Orignialhash,&wg)
  key, err := crypto.CreateKey(32)
  if err != nil {
    panic(fmt.Sprintf("KEY CREATION ERROR: %S",err))
  }
  if conf.Advanced_Settings.Use_Ecryption == true{
    fmt.Printf("ENCRYPTION KEY IS: \n%s\n",hex.EncodeToString(key))
    fmt.Println("Please make note of and store safely, this keep will NOT be saved")
    fmt.Println("Press enter to continue")
    var Test string
    fmt.Scanln(&Test)
  }
  Header()
  for i:=0;i<6;i++{
    wg.Add(1)
    go copy(dst,read,&wg,&Newhash,key)
  }
  //PrintChan(hashChan,&wg)
  wg.Wait()
  Setlog(dst)

  percent,_,missed := VerifyFile(Newhash,Orignialhash)
  MissedLogger.Printf("Transfered %s.",ByteCountSI(SliceSizef(Newhash)))
  MissedLogger.Printf("Successfully transfered %v%.\n",percent)
  if percent != 100{
    var response string
    var times int
    fmt.Println("Not all files succesfully transfered.")
    survey.AskOne(&survey.Input{Message: "Would you like to retry them (y/n): "}, &response)
    if response == "y"{
      survey.AskOne(&survey.Input{Message: "How many times (1-3): "}, &times)
      for i:=0; i < times; i++{
        var tempNewHash []file
        temp := make(chan string)
        //for loop to retry copying specific files
        //removes file from slice if successful repeats until slice is empty or until try count is over
        go func(){
          wg.Add(1)
          for _, files := range missed{
            read <- files.filepath
          }
          wg.Done()
        }()
        go copy(filepath.Join(dst,USER),temp,&wg,&tempNewHash,key)
        wg.Wait()
      }
    }else if response == "n"{
      fmt.Println("Check the missing files log for files not verified as successful. You can manually copy them yourself.")
    }
    //MissedLogger.Println("Missed files:")
    for _,file := range missed{
      MissedLogger.Printf("%s\n",file.filepath)
    }
    err := GetInstalledPrograms(filepath.Join(dst,USER,"UserData"))
    if err != nil{
      panic(fmt.Sprintln("Gathering Installed Programs Error:",err))
    }
  }else{
    fmt.Println("Verified 100% successfully transfered!")
  }

  duration := time.Since(start)

  // Formatted string, such as "2h3m0.5s" or "4.503μs"
  fmt.Println(duration)


}
