package backup

import (
	//"fmt"
	"path/filepath"

	"github.com/michaeldcanady/LookBack/backup2.0/dispatcher"
	"github.com/michaeldcanady/LookBack/backup2.0/file"
	"github.com/michaeldcanady/LookBack/backup2.0/worker"
	"github.com/vbauerster/mpb"
)

// InLineCopy copies all files gathered in Gatherer and sends them directly to thier new location
func InLineCopy(backup bool, dd *dispatcher.Dispatcher, barlist map[string]*mpb.Bar, bars bool, dst string, output chan *file.File) (int64, int64) {

	var i, s int64
	for file := range output{
		dd.Submit(worker.NewJob(i, file, dst, backup))
		i++
		s += (*file.File).Size()
		b := barlist[filepath.Join((*file).PathVol(), (*file).PathUserf(), (*file).PathHead())]
		if b == nil {
			continue
		}
		if bars {
			b.IncrInt64((*file.File).Size())
		}
	}
	return i, s
}
