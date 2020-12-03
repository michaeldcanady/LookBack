package main

import(
  "github.com/BurntSushi/toml"
  "github.com/michaeldcanady/Project01/REWRITE/Project01/libs"
  "github.com/michaeldcanady/SliceTools"
  "fmt"
  "time"
  "github.com/inconshreveable/go-update"
  "net/http"

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
  var err error
  currentTime := time.Now()
  //fmt.Println(conf.Timing.Type)
  //fmt.Println(conf.Timing.Dates)

  dayOfMonth := currentTime.Format("02")
  weekDay := currentTime.Format("Mon")
  time := currentTime.Format("03:04 am")

  switch t := conf.Timing.Type; t{
  case "weekly":
    if SliceTools.SliceIndex(len(conf.Timing.Dates),func(i int)bool{return conf.Timing.Dates[i] == weekDay}) == -1{
      fmt.Println("Not the right backup day")
    }
  case "daily":
    if time != conf.Timing.TimeOfDay{
      fmt.Println("It is not time to update")
    }
  case "monthly":
    if SliceTools.SliceIndex(len(conf.Timing.Dates),func(i int)bool{return conf.Timing.Dates[i] == dayOfMonth}) == -1{
      fmt.Println("Not the right backup day")
    }
  default:
    err = fmt.Errorf("'%s' is not a valid value for [Timing]. Please correct",t)
  }
  if err != nil{
    panic(err)
  }
  if time != conf.Timing.TimeOfDay{
    fmt.Println("It is not time to update")
  }

  //fmt.Println(dayOfMonth,weekDay,time)
}
