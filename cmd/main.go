package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/burhankangsi/LetsYouTube/content"
	"github.com/burhankangsi/LetsYouTube/flash_api"
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

	// params := mux.Vars(r)
	// chanId := params["channelId"]

	// //Creating a new file object
	// var item content.File
	// item.Id = strconv.Itoa(rand.Intn(1000000))
	// item.ChannelId = chanId
	// item.FileName = handler.Filename
	// item.UploadTime = strconv.FormatInt(time.Now().Unix(), 10)
	// item.UploadDate = strconv.Itoa(time.Now().Day())

	// jsonFile, err2 := json.Marshal(item)
	// if err2 != nil {
	// 	log.Fatal("Error converting file struct to json. %v", err2)
	// }

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

//func initServer(ctx context.Context) (*http.Server, context.Context) {
func initServer() {
	//Creating the routers
	log.Infof("Starting server")

	myRouter := mux.NewRouter().StrictSlash(true)
	log.Infof("Calling handler")
	myRouter.HandleFunc("/{channelId}/{videoId}/video.ts", GetVideoObjectHandler).Methods("GET")
	myRouter.HandleFunc("/video/{channelId}", UploadVideoHandler).Methods("POST")
	myRouter.HandleFunc("/video", ListOfVideos).Methods("GET")
	myRouter.HandleFunc("/", HomePage).Methods("GET")

	// srv := &http.Server{
	// 	Addr:    "https://vocal-starship-53117c.netlify.app/endpoint",
	// 	Handler: myRouter,
	// }

	//go func() {
	log.Fatal(http.ListenAndServe(":8000", myRouter))
	//}()
	//http.ListenAndServe("https://vocal-starship-53117c.netlify.app/endpoint", myRouter)
	log.Info("Server started")

	//<-ctx.Done()

	log.Info("Server stopped")

	// ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer func() {
	// 	cancel()
	// }()

	//log.Info("Server exited properly")

	//return srv, ctxShutDown
}

func main() {
	//ctx, _ := context.WithCancel(context.Background())
	//server, ctxShutDown := initServer(ctx)
	initServer()
	// if err := server.Shutdown(ctxShutDown); err != nil {
	// 	log.Fatalf("server Shutdown Failed:%+s", err)
	// }
	log.Infof("Application closed")
}
