package main

import(
  "github.com/BurntSushi/toml"
  "fmt"
  "time"
  "os"

  "github.com/michaeldcanady/Project01/REWRITE/Project01/libs"
  "github.com/michaeldcanady/SliceTools"
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
  var err error
  currentTime := time.Now()
  //fmt.Println(conf.Timing.Type)
  //fmt.Println(conf.Timing.Dates)

  dayOfMonth := currentTime.Format("02")
  weekDay := currentTime.Format("Mon")
  time := currentTime.Format("03:04 am")

  //fmt.Println(dayOfMonth,weekDay,time)
}
