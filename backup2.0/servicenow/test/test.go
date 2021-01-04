package main

import (
	"fmt"

	"github.com/michaeldcanady/Project01/backup2.0/servicenow"
)

func main() {
	fmt.Println(servicenow.Validate(servicenow.Create("dmcanady", "1b7Y9&qP)Mv375!", "libertydev.service-now.com/", "CS0085386")))
}
