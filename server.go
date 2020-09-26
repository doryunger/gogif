package main

import (
	"os"
	"path/filepath"
	"serverTest/factory"
	guuid "github.com/google/uuid"
	"C"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)


type GifHandler struct{}

var uid string
var response string

func genUUID() string {
	id := guuid.New()
	return id.String()
}

//Ignites the conversion process, once finished the output gif file returns to client
func (handler *GifHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	uid := genUUID()
	factory.CreateGif(r, uid)
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote("out.gif"))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, uid+"/out.gif")

}

func (handler *CalcsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "It's alive!")
	
}

func main() {
	//The port our server will be litening on
	port := 8080

	router := mux.NewRouter()

	//Our Routes :
	router.Handle("/gifgen", &GifHandler{}).Methods("POST")
	router.Handle("/test", &TestHandler{}).Methods("GET")

	log.Printf("Starting server. Listening on port %d", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
