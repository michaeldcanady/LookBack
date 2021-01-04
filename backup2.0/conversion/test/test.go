package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/michaeldcanady/Project01/backup2.0/conversion"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func reportGen(files, folders, src, dst int64) {
	def := fmt.Sprintf(
		`Note: /1024 is for size estimation on a PC. /1000 is for Mac.
[Username]:
      files: %v
    folders: %v
        Src: %v bytes
        Dst: %v bytes
   Windows:
        Src: %v
        Dst: %v
       Mac:
        Src: %v
        Dst: %v`, files, folders, src, dst,
		conversion.ByteCountSI(src, 1000), conversion.ByteCountSI(dst, 1000),
		conversion.ByteCountSI(src, 1024), conversion.ByteCountSI(dst, 1024))
	f, err := os.Create("C:\\Users\\dmcanady\\Desktop\\testlogs\\test.log")
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	n4, err := w.WriteString(def)
	check(err)
	fmt.Printf("wrote %d bytes\n", n4)
	w.Flush()
}

func main() {
	reportGen(1000, 500, 25000000, 25000000)
}
