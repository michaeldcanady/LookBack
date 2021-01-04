package main

import (
	"github.com/michaeldcanady/Project01/backup2.0/logging"
)

func main() {

	l := logging.New("C:\\Users\\dmcanady\\Desktop\\testlogs\\test.log")
	l.UserSrc("dmcanady", 1000, 500, 1000000000, 20000, 10000, 50000, 40000)
	l.UserDst("dmcanady", 10000, 10000, 50000, 40000)
	l.UserSrc("dhunter63", 1000, 500, 1000000000, 20000, 10000, 50000, 40000)
	l.UserDst("dhunter63", 20000, 10000, 50000, 40000)
	l.SrcTot()
	l.DstTot()
	l.Print()
}
