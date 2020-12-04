package main

//INSTALL COMMAND: go install -ldflags -H=windowsgui "github.com\michaeldcanady\Project01\REWRITE\Lookback - NotificationService"

import(
  "github.com/robfig/cron/v3"
  "fmt"
  "time"
  "github.com/gen2brain/beeep"
  "sync"
)

const(
  Title = "LookBack"
)

var(
  wg sync.WaitGroup
  hr12  = time.Hour*12
  hr6   = time.Hour*6
  hr    = time.Hour*1
  min30 = time.Minute*30
  min15 = time.Minute*15
  min5  = time.Minute*5
)

func Notification(t time.Duration){
}

func main(){
  backupTime := CheckTime()
  //format := "Mon Jul 9 15:00 2012"
  noticationTimes := []time.Duration{hr12, hr6, hr, min30, min15, min5}
  c := cron.New()

  for _,t := range noticationTimes{
    t:=t
    if time.Now().After(backupTime.Add(-t)){
      continue
    }
    _, err := c.AddFunc(date(backupTime,t),func(){
      var t1 string
      if t == time.Duration(time.Minute*0){
        t1 = fmt.Sprintf("Your backup will begin now")
      }else{
        t1 = fmt.Sprintf("Your backup will begin in %v", t)
      }
      err := beeep.Notify(Title, t1, "TEMP_LOGO.png")
      if err != nil {
        panic(err)
      }
    })
    if err != nil {
      panic(err)
    }
  }
  c.Stop()
  c.Run()

}

func date(backupTime time.Time, d time.Duration)string{
  backupTime = backupTime.Add(-d)
  return fmt.Sprintf("%d %d %02d %02d %d",backupTime.Minute(),backupTime.Hour(), backupTime.Day(), backupTime.Month(), backupTime.Weekday())
}

func CheckTime()time.Time{
    var t time.Time
    // 24 hr time
    t1 := "18:00 EST"

    currentTime := time.Now()

    format := "02 01 2006 15:04 MST"
    date := fmt.Sprintf("%s %s",currentTime.Format("02 01 2006"),t1)

    tm, _ := time.Parse(format, date)
    currenttime := currentTime.Format(format)
    current,_ := time.Parse(format, currenttime)

    if current.After(tm){
      currentTime = currentTime.AddDate(0,0,1)
      date = fmt.Sprintf("%s %s",currentTime.Format("02 01 2006"),t1)
      tm, _ = time.Parse(format, date)
      return tm
    }else if currentTime.Before(tm){
      return tm
    }
    return t
}
