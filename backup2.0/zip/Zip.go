package zipfile

import (
	"archive/zip"
	"io"
	"log"
	"time"
	"os"
	"path/filepath"
	"strings"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func AddFileToZip(zipWriter *zip.Writer, filename, volume string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	zipPath,err := filepath.Rel(volume, filename)
	if err != nil {
		return err
	}
	zipPath = strings.Replace(zipPath, "\\", "/", -1)
	zipPath = strings.TrimLeft(zipPath, "/")

	zipFile, err := zipWriter.Create(zipPath)
	if err != nil {
		return err
	}


	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Store

	_, err = io.Copy(zipFile, fileToZip)
	return err
}
