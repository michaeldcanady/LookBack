package main

//INSTALL COMMAND: go install -ldflags -H=windowsgui "github.com\michaeldcanady\Project01\REWRITE\Lookback - NotificationService"

import(
  "github.com/robfig/cron/v3"
  "fmt"
  "time"
  "github.com/gen2brain/beeep"
  "sync"
  "os/exec"
)

const(
  Title = "LookBack"
)

var(
  wg sync.WaitGroup
  noticationTimes = []time.Duration{time.Hour*12, time.Hour*6, time.Hour*1, time.Minute*30, time.Minute*15, time.Minute*5, time.Minute*0}
)

// Creates a notification for the user when backups will be happening
func Notification(t time.Duration){
  var t1 string
  if t == time.Duration(-time.Minute*0){
    t1 = fmt.Sprintf("Your backup will begin now")
    //exec.Command("LookBack.exe","-backup")
  }else{
    t1 = fmt.Sprintf("Your backup will begin in %v", t)
  }
  err := beeep.Notify(Title, t1, "TEMP_LOGO.png")
  if err != nil {
    panic(err)
  }
}

func main(){
  // Gets when the backup is preformed
  backupTime := CheckTime()
  // Conduit for task sequencing
  c := cron.New()
  // Creates all needed notifications
  for _,t := range noticationTimes{
    t:=t
    if time.Now().After(backupTime.Add(-t)){
      continue
    }
    // sets up task, when it will be preformed and what will be preformed
    _, err := c.AddFunc(date(backupTime,t),func(){Notification(t)})
    if err != nil {
      panic(err)
    }
  }
  // Stops the scheduler
  c.Stop()
  // Waits for scheduled events
  c.Run()

}

func date(backupTime time.Time, d time.Duration)string{
  backupTime = backupTime.Add(-d)
  return fmt.Sprintf("%d %d %02d %02d %d",backupTime.Minute(),backupTime.Hour(), backupTime.Day(), backupTime.Month(), backupTime.Weekday())
}

// Gets rules for backups from settings

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
