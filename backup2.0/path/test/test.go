package main

import (
	"fmt"

	"github.com/michaeldcanady/Project01/backup2.0/path"
)

func main() {
	test := filestruct.New("C:\\Users\\dmcanady\\Desktop\\Compare")
	fmt.Println(test)
	fmt.Println(test.ZipPath())
}
