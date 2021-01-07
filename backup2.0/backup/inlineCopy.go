package backup

import (
	"fmt"
	"path/filepath"

	"github.com/michaeldcanady/Project01/backup2.0/dispatcher"
	"github.com/michaeldcanady/Project01/backup2.0/file"
	"github.com/michaeldcanady/Project01/backup2.0/worker"
	"github.com/vbauerster/mpb"
)

// InLineCopy copies all files gathered in Gatherer and sends them directly to thier new location
func InLineCopy(backup bool, dd *dispatcher.Dispatcher, barlist map[string]*mpb.Bar, bars bool, dst string, output chan *file.File) (int64, int64) {

	var i, s int64
copyloop:
	for {
		select {
		case file, ok := (<-output):
			if ok {
				dd.Submit(worker.NewJob(i, file, dst, backup))
				i++
				s += (*file.File).Size()
				b := barlist[filepath.Join((*file).PathVol(), (*file).PathUserf(), (*file).PathHead())]
				if b == nil {
					continue
				}
				if bars {
					b.IncrInt64(file.RawSize())
				}
			} else {
				fmt.Println("Channel closed!")
				break copyloop
			}
		default:
			continue
		}
	}
	return i, s
}
