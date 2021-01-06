package worker

import (
	//"fmt"
	"log"
	"runtime"

	"github.com/michaeldcanady/Project01/backup2.0/copy"
	"github.com/michaeldcanady/Project01/backup2.0/file"
)

var UNIT int64

func init() {
	if runtime.GOOS == "windows" {
		UNIT = 1024
	} else {
		UNIT = 1000
	}
}

type Job struct {
	ID     int64
	File   *file.File
	Dst    string
	Backup bool
}

type JobChannel chan Job
type JobQueue chan chan Job

func NewJob(id int64, file *file.File, dst string, backup bool) Job {
	return Job{
		ID:     id,
		File:   file,
		Dst:    dst,
		Backup: backup,
	}
}

type Worker struct {
	ID      int
	JobChan JobChannel
	Queue   JobQueue // shared between all workers
	Quit    chan struct{}
}

func New(ID int, JobChan JobChannel, Queue JobQueue, Quit chan struct{}) *Worker {
	return &Worker{
		ID:      ID,
		JobChan: JobChan,
		Queue:   Queue,
		Quit:    Quit,
	}
}

func (wr *Worker) Start() {
	go func() {
		for {
			wr.Queue <- wr.JobChan
			select {
			case job := <-wr.JobChan:
				copy.Copy(job.Dst, job.File, UNIT, job.Backup)
			case <-wr.Quit:
				close(wr.JobChan)
				return
			}
		}
	}()
}

func (wr *Worker) Stop() {
	close(wr.Quit)
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
