package factory

import (
	"C"
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"serverTest/dispatcher"
	"serverTest/worker"
	"strconv"
	"time"
	"gocv.io/x/gocv"
)

//Initialize global variables
var count int
var files []string
var uuid string
var framesCount float64
var frames []int
var err error
var num int


//Generates a gif image out of the data in the request
func CreateGif(r *http.Request, id string) {
	//Creates a folder with a unique uuid value
	uuid = id
	frames = frames[:0]
	files = files[:0]
	_, err := os.Stat("folder")

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(uuid, 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}
	count = 0
	//Reads the request's data
	file, header, err := r.FormFile("file")
	println(header)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//Convert the data into buffer type
	buf, err := ioutil.ReadAll(file)
	err = ioutil.WriteFile(uuid+"/"+"test.mp4", buf, 0644)
	if err != nil {
		panic(err)
	}

	img := gocv.NewMat()
	defer img.Close()
	//Starting a conversion session by capturing a video file
	avi, _ := gocv.VideoCaptureFile(uuid + "/" + "test.mp4")
	//Arraging thr frames in an array for later use
	framesCount := avi.Get(7)
	for i := 1; i < int(framesCount); i++ {
		if i%2 == 0 {
			frames = append(frames, i)
			files = append(files, strconv.Itoa(i)+".gif")
		}
	}
	//Dispatching extraction tasks concurrently (using GO's worker pool)
	dd := dispatcher.New(25).Start()

	for _, s := range frames {

		dd.Submit(worker.Job{
			NUM: s,
			ID:  uuid,
		})
		time.Sleep(250 * time.Millisecond)
	}
	time.Sleep(350 * time.Millisecond)
	
	//After all the frames have been extracted from the video file, those frames are encoded to a single gif file
	outGif := &gif.GIF{}
	for _, name := range files {
		f, _ := os.Open(uuid + "/" + name)
		inGif, _ := gif.Decode(f)
		f.Close()

		outGif.Image = append(outGif.Image, inGif.(*image.Paletted))
		outGif.Delay = append(outGif.Delay, 0)
	}

	//Saving the output on disk
	f, _ := os.OpenFile(uuid+"/"+"out.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, outGif)

}
