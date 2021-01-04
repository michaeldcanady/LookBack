package dispatcher

import (
	"github.com/michaeldcanady/Project01/backup2.0/worker"
)

type Dispatcher struct {
	Workers  []*worker.Worker
	WorkChan worker.JobChannel
	Queue    worker.JobQueue
}

func New(num int) *Dispatcher {
	return &Dispatcher{
		Workers:  make([]*worker.Worker, num),
		WorkChan: make(worker.JobChannel),
		Queue:    make(worker.JobQueue),
	}
}

func (d *Dispatcher) Start() *Dispatcher {
	l := len(d.Workers)
	for i := 1; i <= l; i++ {
		wrk := worker.New(i, make(worker.JobChannel), d.Queue, make(chan struct{}))
		wrk.Start()
		d.Workers = append(d.Workers, wrk)
	}
	go d.Process()
	return d
}

func (d *Dispatcher) Process() {
	for {
		select {
		case job := <-d.WorkChan:
			jobChan := <-d.Queue
			jobChan <- job
		}
	}
}

func (d *Dispatcher) Submit(job worker.Job) {
	d.WorkChan <- job
}
