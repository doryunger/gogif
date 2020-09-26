package dispatcher

import (
	"serverTest/worker"
)

// New returns a new dispatcher. A Dispatcher communicates between the client שnd the worker. 
func New(num int) *disp {
	return &disp{
		Workers:  make([]*worker.Worker, num),
		WorkChan: make(worker.JobChannel),
		Queue:    make(worker.JobQueue),
	}
}

// disp is the link between the client and the workers
type disp struct {
	Workers  []*worker.Worker  
	WorkChan worker.JobChannel 
	Queue    worker.JobQueue   
}

// Start creates pool of num count of workers.
func (d *disp) Start() *disp {
	l := len(d.Workers)
	for i := 1; i <= l; i++ {
		wrk := worker.New(i, make(worker.JobChannel), d.Queue, make(chan struct{}))
		wrk.Start()
		d.Workers = append(d.Workers, wrk)
	}
	go d.process()
	return d
}

// process listens to a job submitted on WorkChan and relays it to the WorkPool.
func (d *disp) process() {
	for {
		select {
		case job := <-d.WorkChan: // listen to any submitted job on the WorkChan
			jobChan := <-d.Queue

			// Once a jobChan is available, send the submitted Job on this JobChan
			jobChan <- job
		}
	}
}

func (d *disp) Submit(job worker.Job) {
	d.WorkChan <- job
}
