package servicenow

import (
	"fmt"
	"time"

	"github.com/mjmorell/GoServe"
)

var (
	client goserve.Client
)

type Back struct {
	Client   goserve.Client
	Username string
	CSNumber string
}

const format = "Mon Jan _2 2006 at 15:04 am"

func Validate(client Back) bool {
	tech, _ := client.Client.PULL("sys_user", goserve.Filter("user_name")+goserve.IS(client.Username))

	_, num := client.Client.PULL("u_computer_support", goserve.Filter("number")+goserve.IS(client.CSNumber)+
		goserve.AND()+goserve.Filter("state")+goserve.ISNOT("3")+goserve.AND()+goserve.Filter("assigned_to")+goserve.IS(tech[0]["sys_id"]))
	if num == 0 {
		return false
	}
	return true
}

func Create(username, password, instance, csnumber string) Back {

	client := Back{Client: goserve.Client{Username: username, Password: password, Instance: instance}, Username: username, CSNumber: csnumber}
	return client
}

func Start(client Back) {
	pushable := make(map[string]string)
	t := time.Now()
	pushable["work_notes"] = fmt.Sprintf("backup began on %v", t.Format(format))
	client.Client.PUSH("u_computer_support", goserve.Filter("number")+goserve.IS(client.CSNumber), pushable)
}

func Pause(client Back) {
	pushable := make(map[string]string)
	t := time.Now()
	pushable["work_notes"] = fmt.Sprintf("backup Paused on %v", t.Format(format))
	client.Client.PUSH("u_computer_support", goserve.Filter("number")+goserve.IS(client.CSNumber), pushable)
}

func Resume(client Back) {
	pushable := make(map[string]string)
	t := time.Now()
	pushable["work_notes"] = fmt.Sprintf("backup Resumed on %v", t.Format(format))
	client.Client.PUSH("u_computer_support", goserve.Filter("number")+goserve.IS(client.CSNumber), pushable)
}

func Stop(client Back, reason string) {
	pushable := make(map[string]string)
	t := time.Now()
	pushable["work_notes"] = fmt.Sprintf("backup Stop on %v because %s", t.Format(format), reason)
	client.Client.PUSH("u_computer_support", goserve.Filter("number")+goserve.IS(client.CSNumber), pushable)
}

func Finish(client Back, opts ...map[string]interface{}) {
	pushable := make(map[string]string)
	t := time.Now()
	o := "Breakdown\n"
	for _, opt := range opts {
		for k, v := range opt {
			o += fmt.Sprintf("%s: %v\n", k, v)
		}
	}
	pushable["work_notes"] = fmt.Sprintf("backup Finished on %v\n %s", t.Format(format), o)
	client.Client.PUSH("u_computer_support", goserve.Filter("number")+goserve.IS(client.CSNumber), pushable)
}
