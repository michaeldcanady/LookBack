package main

//INSTALL COMMAND: go install -ldflags -H=windowsgui "github.com\michaeldcanady\Project01\REWRITE\Lookback - NotificationService"

import(
  "fmt"
  "time"
  "sync"
  "strconv"
  //"os/exec"

  "github.com/michaeldcanady/Project01/REWRITE/Project01/libs"
  "github.com/gen2brain/beeep"
  "github.com/robfig/cron/v3"
  "github.com/BurntSushi/toml"
)

const(
  Title = "LookBack"
)

// Create custom duration type to cater to extended needs
type Duration time.Duration

// Create customer duration string type to cater to extended needs
func (d Duration) String() string {
	// Largest time is 2540400h10m10.000000000s
	var buf [32]byte
	w := len(buf)

	u := uint64(d)
	neg := d < 0
	if neg {
		u = -u
	}

	if u < uint64(time.Second) {
		// Special case: if duration is smaller than a second,
		// use smaller units, like 1.2ms
		var prec int
		w--
		buf[w] = 's'
		w--
		switch {
		case u == 0:
			return "0s"
		case u < uint64(time.Microsecond):
			// print nanoseconds
			prec = 0
			buf[w] = 'n'
		case u < uint64(time.Millisecond):
			// print microseconds
			prec = 3
			// U+00B5 'µ' micro sign == 0xC2 0xB5
			w-- // Need room for two bytes.
			copy(buf[w:], "µ")
		default:
			// print milliseconds
			prec = 6
			buf[w] = 'm'
		}
		w, u = fmtFrac(buf[:w], u, prec)
		w = fmtInt(buf[:w], u)
	} else {
		w--
		buf[w] = 's'

		w, u = fmtFrac(buf[:w], u, 9)

		// u is now integer seconds
		w = fmtInt(buf[:w], u%60)
		u /= 60

		// u is now integer minutes
		if u > 0 {
			w--
			buf[w] = 'm'
			w = fmtInt(buf[:w], u%60)
			u /= 60

			// u is now integer hours
			// Stop at hours because days can be different lengths.
			if u > 0 {
				w--
				buf[w] = 'h'
				w = fmtInt(buf[:w], u%24)
        u /= 24

        if u > 0 {
  				w--
  				buf[w] = 'd'
          w = fmtInt(buf[:w], u%7)
          u /= 7

          if u > 0 {
    				w--
    				buf[w] = 'w'
            w = fmtInt(buf[:w], u%52)
            u /= 52

            if u > 0 {
      				w--
      				buf[w] = 'y'
              w = fmtInt(buf[:w], u)
            }
          }
        }
			}
		}
	}

	if neg {
		w--
		buf[w] = '-'
	}

	return string(buf[w:])
}

func fmtFrac(buf []byte, v uint64, prec int) (nw int, nv uint64) {
	// Omit trailing zeros up to and including decimal point.
	w := len(buf)
	print := false
	for i := 0; i < prec; i++ {
		digit := v % 10
		print = print || digit != 0
		if print {
			w--
			buf[w] = byte(digit) + '0'
		}
		v /= 10
	}
	if print {
		w--
		buf[w] = '.'
	}
	return w, v
}

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}

var(
  wg sync.WaitGroup
  conf libs.Config
  Week = Day*7
  Day = Hour*24
  Hour = Duration(time.Hour)
  Minute = Duration(time.Minute)
  hourlyNoticationTimes = []Duration{Minute*57,Minute*30, Minute*15, Minute*5, Minute*0}
  dailyNoticationTimes = append([]Duration{Hour*12, Hour*6, Hour*1, Minute*50},hourlyNoticationTimes...)
  weeklyNoticationTimes = append([]Duration{Day*3,Day}, dailyNoticationTimes...)
  monthlyNoticationTimes = append([]Duration{Week*2,Week}, weeklyNoticationTimes...)
  yearlyNoticationTimes = append([]Duration{Day*60,Day*30},monthlyNoticationTimes...)
)

// Creates a notification for the user when backups will be happening
func Notification(t Duration){
  var t1 string
  fmt.Println(t)
  if t == Duration(-Minute*0){
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

// Adjusts notification cycles to match user selection
func notTimes(timeFrame string)[]Duration{
  switch timeFrame{
  case "yearly":
    return yearlyNoticationTimes
  case "monthly":
    return monthlyNoticationTimes
  case "weekly":
    return weeklyNoticationTimes
  case "daily":
    return dailyNoticationTimes
  case "hourly":
    return hourlyNoticationTimes
  default:
    panic(fmt.Errorf("%s is not a valid selection.",timeFrame))
  }
}

func date(backupTime time.Time, d Duration,recursion string)string{
  backupTime = backupTime.Add(time.Duration(-d))
  minute,hour,day,month,weekday := backupTime.Minute(),backupTime.Hour(), backupTime.Day(), backupTime.Month(), backupTime.Weekday()
  // For user to set increments
  switch recursion{
  case "yearly":
    return fmt.Sprintf("%d %d %02d %02d %s",minute,hour,day,month,"*")
  case "monthly":
    return fmt.Sprintf("%d %d %02d %s %s",minute,hour,day,"*","*")
  case "weekly":
    return fmt.Sprintf("%d %d %s %s %d",minute,hour,"*","*",weekday)
  case "daily":
    return fmt.Sprintf("%d %d %s %s %s",minute,hour,"*","*","*")
  case "hourly":
    return fmt.Sprintf("%d %s %s %s %s",minute,"*","*","*","*")
  default:
    return fmt.Sprintf("%d %d %02d %02d %d",minute,hour,day,month,weekday)
  }
}

// Gets rules for backups from settings
// Potential will remove
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

func contains(e int,s []int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func RangeCheck(s []string){
  newSlice := []int{}
  for _,i := range s{
    i1,err := strconv.Atoi(i)
    if err != nil{
      panic(err)
    }
    newSlice = append(newSlice,i1)
  }
  for _,ss := range newSlice{
    if contains(ss+1,newSlice){
      fmt.Println("contained")
    }else if contains(ss-1,newSlice){
      
    }
  }
}

func init(){
  if _, err := toml.DecodeFile("C:\\Go\\src\\github.com\\michaeldcanady\\Project01\\REWRITE\\Project01\\libs\\settings.toml", &conf); err != nil {
    panic(err)
  }
}

func main(){
  // Gets when the backup is preformed
  backupTime := CheckTime()
  // Conduit for task sequencing
  //conf.Timing.TimeOfDay
  RangeCheck(conf.Timing.Dates)
  c := cron.New()

  // Creates all needed notifications
  for _,t := range notTimes(conf.Timing.Type){
    t:=t
    fmt.Println(t)
    if time.Now().After(backupTime.Add(-time.Duration(t))){
      continue
    }
    // sets up task, when it will be preformed and what will be preformed
    _, err := c.AddFunc(date(backupTime,t,conf.Timing.Type),func(){Notification(t)})
    if err != nil {
      panic(err)
    }
  }
  // Stops the scheduler
  c.Stop()
  // Waits for scheduled events
  c.Run()

}
