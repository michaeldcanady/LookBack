package zipfile

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/michaeldcanady/LookBack/backup2.0/file"
)

func ZipWriter(files chan *file.File) *sync.WaitGroup {
	f, err := os.Create("C:\\out2.zip")
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	zw := zip.NewWriter(f)
	go func() {
		// Note the order (LIFO):
		defer wg.Done() // 2. signal that we're done
		defer f.Close() // 1. close the file
		var err error
		var fw io.Writer
		for f2 := range files {
			// Loop until channel is closed.
			name, _ := filepath.Rel("C:\\Users\\", f2.PathStr())
			if fw, err = zw.Create(name); err != nil {
				panic(err)
			}

			f1, _ := os.Open(f2.PathStr())

			io.Copy(fw, f1)
			f1.Close()
		}
		// The zip writer must be closed *before* f.Close() is called!
		if err = zw.Close(); err != nil {
			panic(err)
		}
	}()
	return &wg
}
