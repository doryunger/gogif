package worker

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"os"
	"strconv"

	"gocv.io/x/gocv"
)

// Job represents a single entity that should be processed.
type Job struct {
	NUM int
	ID  string
}

type JobChannel chan Job
type JobQueue chan chan Job

// Worker is a a single processor. Typically its possible to
type Worker struct {
	ID      int           
	JobChan JobChannel    
	Queue   JobQueue     
	Quit    chan struct{} 
}
//A function that initializes a worker instance
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
			// when available, put the JobChan again on the JobPool
			// and wait to receive a job
			wr.Queue <- wr.JobChan
			select {
			case job := <-wr.JobChan:
				// when a job is received, process
				SaveGif(job.NUM, job.ID)
			case <-wr.Quit:
				// a signal on this channel means someone triggered
				// a shutdown for this worker
				close(wr.JobChan)
				return
			}
		}
	}()
}

// stop closes the Quit channel on the worker.
func (wr *Worker) Stop() {
	close(wr.Quit)
}
//Defines the method which will be assigned to task 
func SaveGif(pos int, uuid string) {
	mat := gocv.NewMat()
	img := gocv.NewMat()
	avi, _ := gocv.VideoCaptureFile(uuid + "/" + "test.mp4")
	avi.Set(1, float64(pos))
	avi.Read(&img)
	gocv.Resize(img, &mat, image.Point{int(int(avi.Get(3)) / 2), int(int(avi.Get(4)) / 2)}, float64(0), float64(0), 0)
	buf, _ := gocv.IMEncode(".jpg", mat)
	out, err := os.Create(uuid + "/" + strconv.Itoa(pos) + ".gif")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	jpy, err := jpeg.Decode(bytes.NewReader(buf))
	err = gif.Encode(out, jpy, &gif.Options{NumColors: 256})
}
