package backup

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/michaeldcanady/LookBack/backup2.0/file"
	"github.com/michaeldcanady/LookBack/backup2.0/zip"
	"github.com/vbauerster/mpb"
)

// ZipCopy Copies all users and thier files in to a single zip file
func ZipCopy(backup bool, barlist map[string]*mpb.Bar, bars bool, dst string, output chan *file.File) (int64, int64) {
	//Creates temporary zip file
	newZipFile, err := ioutil.TempFile("", "Users*.zip")
	if err != nil {
		panic(err)
	}
	//Closes zipfile when done
	defer newZipFile.Close()
	//Creates a zip writer used to zip files
	zipWriter := zip.NewWriter(newZipFile)
	//Closes zip writer when done
	defer zipWriter.Close()
	// i is number of files, s is size of files
	var i, s int64
	// Add files to zip
ziploop:
	for {
		select {
		// Checks is files are avaliable from output
		case file, ok := <-output:
			// if a file is recieved and the channel is still open
			if ok {
				//first parameter zip.Writer, full filepath, elements being removed
				if err = zipfile.AddFileToZip(zipWriter, file.PathStr(), file.PathVol()+"\\Users"); err != nil {
					panic(err)
				}
				// Adds another file to the count
				i++
				size := file.RawSize()
				// Adds the files size to the count
				s += size
				//Get the volume,
				path := filepath.Join((*file).PathVol(), (*file).PathUserf(), (*file).PathHead())
				// Updates the file's respective bar
				b := barlist[path]
				// If bar does not exist pass
				// Need to add info log write here
				// No files without a bar should exist
				if b == nil {
					continue
				}
				// if bars are enabled
				if bars {
					//increment by size of the file
					b.IncrInt64(size)
				}
				// if channel is closed
			} else {
				// change to info writer when the channel was closed
				fmt.Println("Channel closed")
				// breaks the for loop
				break ziploop
			}
			// if no files are avaliable loop continues
		default:
			continue
		}
	}
	// Gets the file name of the temp zip file
	_, name := filepath.Split(newZipFile.Name())
	// Change to info writer for location of the zip file
	fmt.Println(filepath.Join(dst, name))
	// Moves file from temp location to new location
	err = os.Rename(newZipFile.Name(), filepath.Join(dst, name))
	// if move failes delete the temp file
	if err != nil {
		// add error writer for what error
		// add info writer to validate that file was deleted
		os.Remove(newZipFile.Name())
	}
	// returns total files, and size
	return i, s
}
