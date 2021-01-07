package logging

import (
	"bufio"
	"io"
	"os"

	"github.com/michaeldcanady/LookBack/backup2.0/conversion"
	slicetils "github.com/michaeldcanady/SliceTils/SliceTils"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Logging struct {
	Path     io.Writer
	Users    []*UserData
	SrcTotal Total
	DstTotal Total
}

type Total struct {
	Files   int64
	Folders int64
	Size    int64
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func write(s string, L *Logging) {
	w := bufio.NewWriter(L.Path)
	_, err := w.WriteString(s)
	check(err)
	w.Flush()
}

// Creates a new file
func New(path string) *Logging {
	f, err := os.Create(path)
	check(err)
	return &Logging{Path: f}
}

//Adds basic user information for the log
func (L *Logging) UserSrc(username string, files, folders, size, desk, down, pics, docs int64) {
	u := &UserData{
		Username: username,
		Files:    files,
		Folders:  folders,
		Size:     size,
		Src: BreakDown{
			Desk: desk,
			Docs: docs,
			Down: down,
			Pics: pics,
		},
	}
	L.Users = append(L.Users, u)
}

func (L *Logging) UserDst(username string, desk, down, pics, docs int64) {
	for _, users := range L.Users {
		if users.Username == username {
			users.Dst = BreakDown{
				Desk: desk,
				Docs: docs,
				Down: down,
				Pics: pics,
			}
		}
	}
}

func (L *Logging) SrcTot() {
	for _, users := range L.Users {
		L.SrcTotal.Files += users.Files
		L.SrcTotal.Folders += users.Folders
		L.SrcTotal.Size += users.Size
	}
}

func (L *Logging) DstTot() {
	for _, users := range L.Users {
		L.DstTotal.Files += users.Files
		L.DstTotal.Folders += users.Folders
		L.DstTotal.Size += users.Size
	}
}

func (L *Logging) CPrint(U *UserData) {
	test := []int{CountDigits(U.Src.Desk / UNIT), CountDigits(U.Src.Docs / UNIT), CountDigits(U.Src.Down / UNIT), CountDigits(U.Src.Pics / UNIT)}
	size := slicetils.Max(test).(int) + 2
	p := message.NewPrinter(language.English)
	s := p.Sprintf(`
  %s:
       Desktop: %v (%06.2f%%)
     Documents: %v (%06.2f%%)
     Downloads: %v (%06.2f%%)
      Pictures: %v (%06.2f%%)`, U.Username, conversion.ByteCountSI(U.Src.Desk, UNIT, size), float64(U.Dst.Desk)/float64(U.Src.Desk)*100,
		conversion.ByteCountSI(U.Src.Docs, 1024, size), float64(U.Dst.Docs)/float64(U.Src.Docs)*100,
		conversion.ByteCountSI(U.Src.Down, 1024, size), float64(U.Dst.Down)/float64(U.Src.Down)*100,
		conversion.ByteCountSI(U.Src.Pics, 1024, size), float64(U.Dst.Pics)/float64(U.Src.Pics)*100)
	w := bufio.NewWriter(L.Path)
	_, err := w.WriteString(s)
	check(err)
	w.Flush()
}

func (L *Logging) Print() {
	for _, users := range L.Users {
		L.CPrint(users)
	}
}

func CountDigits(i int64) (count int) {
	for i != 0 {

		i /= 10
		count = count + 1
	}
	return count
}
