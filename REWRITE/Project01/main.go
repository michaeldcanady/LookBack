package main

import(
  "fmt"
  "time"
  "os"
  "flag"

  "github.com/BurntSushi/toml"
  "github.com/michaeldcanady/Project01/REWRITE/Project01/libs"
  "github.com/AlecAivazis/survey/v2"
  term"github.com/AlecAivazis/survey/v2/terminal"
  //"github.com/michaeldcanady/SliceTools"
)

var(
  conf libs.Config
)

func init(){
  var backup,silent,update,restore,image bool

  flag.BoolVar(&backup,"backup", false, "Start Backup Without Prompts")
  flag.BoolVar(&silent,"s", false, "No GUI displayed")
  flag.BoolVar(&update,"update", false, "Check for and update program if needed")
  flag.BoolVar(&restore,"restore", false, "Start Restore Without Prompts")
  flag.BoolVar(&image,"image", false, "Images the computer after backup")
  flag.Parse()

  if update{
    equinoxUpdate()
  }

  if _, err := toml.DecodeFile("C:\\Go\\src\\github.com\\michaeldcanady\\Project01\\REWRITE\\Project01\\libs\\settings.toml", &conf); err != nil {
    panic(err)
  }

}

func HomeScreen()string{
  Header()
  // SELECT TO RESTORE OR BACKUP
  bORr :=  []string{"Start Backup", "Start Restore", "Update Check()", "Settings","Exit"}
  BorR := 0
  prompt := &survey.Select{
    Message: "",
    Options: bORr,
  }
  err := survey.AskOne(prompt, &BorR)
  if err == term.InterruptErr {
	   exit()
  } else if err != nil {
	   panic(err)
  }
  if bORr[BorR] == "Exit"{
    os.Exit(0)
  }else{
    return bORr[BorR]
  }
  return ""
}

// Proceedure for user initiated exit
func exit(){
  fmt.Println("interrupted")
  os.Exit(0)
}

func main(){
  HomeScreen()

  //var err error
  currentTime := time.Now()
  //fmt.Println(conf.Timing.Type)
  //fmt.Println(conf.Timing.Dates)

  dayOfMonth := currentTime.Format("02")
  weekDay := currentTime.Format("Mon")
  Time := currentTime.Format("03:04 am")
  fmt.Println(dayOfMonth,weekDay)

//  switch t := conf.Timing.Type; t{
//  case "weekly":
//    if SliceTools.SliceIndex(len(conf.Timing.Dates),func(i int)bool{return conf.Timing.Dates[i] == weekDay}) == -1{
//      fmt.Println("Not the right backup day")
//    }
//  case "daily":
//    if time != conf.Timing.TimeOfDay{
//      fmt.Println("It is not time to update")
//    }
//  case "monthly":
//    if SliceTools.SliceIndex(len(conf.Timing.Dates),func(i int)bool{return conf.Timing.Dates[i] == dayOfMonth}) == -1{
//      fmt.Println("Not the right backup day")
//    }
//  default:
//    err = fmt.Errorf("'%s' is not a valid value for [Timing]. Please correct",t)
//  }
//  if err != nil{
//    panic(err)
//  }
  if Time != conf.Timing.TimeOfDay{
    fmt.Printf("Backup will begin at %s",conf.Timing.TimeOfDay)
  }

  //fmt.Println(dayOfMonth,weekDay,time)
}
