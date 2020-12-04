package main

import(
  "github.com/BurntSushi/toml"
  "github.com/michaeldcanady/Project01/REWRITE/Project01/libs"
  //"github.com/michaeldcanady/SliceTools"
  "fmt"
  "time"
  "os"

)

var(
  conf libs.Config
)

func init(){
  if _, err := toml.DecodeFile("C:\\Go\\src\\github.com\\michaeldcanady\\Project01\\REWRITE\\Project01\\libs\\settings.toml", &conf); err != nil {
    panic(err)
  }

}

func main(){
  if len(os.Args) == 2 && os.Args[1] == "update" {
    equinoxUpdate()
  }
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
