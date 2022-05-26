package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"content"

	"./flash_api"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func UploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Please upload the video")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 100 MB files.
	r.ParseMultipartForm(100 << 20)

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	log.Infof("Uploaded File: %+v\n", handler.Filename)
	log.Infof("File size: %+v\n", handler.Size)
	log.Infof("MIME Header: %+v\n", handler.Header)

	prod, err1 := content.ConfigureProducer()
	if err1 != nil {
		log.Errorf("Error creating sarama producer. %v", err1)
	}

	f, err := os.OpenFile("./uploads/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Infof("Could not open ts file. %v", err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	if err3 := content.UploadToTopic(prod, handler.Filename); err3 != nil {
		log.Errorf("Error uploading video to topic. %v", err3)
	}
}

func GetVideoObjectHandler(W http.ResponseWriter, R *http.Request) {
	log.Info("Fetching video...please wait")
	W.Header().Set("Content-Type", "application/json")
	params := mux.Vars(R)
	VideoId, ok := params["videoId"]
	if !ok {
		log.Errorf("Video ID is missing in parameters")
	}
	ChannelId, ok1 := params["channelId"]
	if !ok1 {
		log.Errorf("Channel ID is missing in parameters")
	}

	_, err := flash_api.GetVideoObject(W, R, VideoId, ChannelId)
	if err != nil {
		log.Fatal("Error in getting the video. Please try again")
		return
	}
	//json.NewEncoder(W).Encode(item)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welecome to my api")
}

func ListOfVideos(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos will be displayed here")
}

func initServer() {
	//Creating the routers
	log.Infof("Starting server")

	myRouter := mux.NewRouter().StrictSlash(true)
	log.Infof("Calling handler")
	myRouter.HandleFunc("/{channelId}/{videoId}/video.ts", GetVideoObjectHandler).Methods("GET")
	myRouter.HandleFunc("/video/{channelId}", UploadVideoHandler).Methods("POST")
	myRouter.HandleFunc("/video", ListOfVideos).Methods("GET")
	myRouter.HandleFunc("/", HomePage).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", myRouter))
	log.Info("Server started")

	log.Info("Server stopped")
}

func main() {
	initServer()
	log.Infof("Application closed")
}
