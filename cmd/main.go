package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/burhankangsi/LetsYouTube/content"
	"github.com/burhankangsi/LetsYouTube/flash_api"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	host     = "youtubeclone1.c7itawripsrx.us-east-1.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "Postgres"
	dbname   = "youtubeclone"
)

type Video struct {
	videoID      int    `json:"videoid"`
	videoName    string `json:"videoname"`
	duration     int64  `json:"duration"`
	channelID    int    `json:"channelid"`
	title        string `json:"title"`
	channelImage string `json:"channelimage"`
	views        string `json:"views"`
	timestamp    string `json:"timestamp"`
	channelName  string `json:"channelname"`
	uploadDate   string `json:"uploaddate"`
	uploadTime   string `json:"uploadtime"`
	thumbnail    string `json:"thumbnail"`
}

type JsonResponse struct {
	Type    string  `json:"type"`
	Data    []Video `json:"data"`
	Message string  `json:"message"`
}

type JsonResponse1 struct {
	Type    string `json:"type"`
	Data    Video  `json:"data"`
	Message string `json:"message"`
}

type JsonResponse2 struct {
	Type    string   `json:"type"`
	Data    []string `json:"data"`
	Message string   `json:"message"`
}

const (
	AWS_S3_REGION = "us-east-1"
	AWS_S3_BUCKET = "youtube-clone-bk"
)

func uploadImageToS3(filepath string) error {
	upFile, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer upFile.Close()

	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	session, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION)})
	if err != nil {
		log.Fatal(err)
	}

	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(AWS_S3_BUCKET),
		Key:                  aws.String(filepath),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

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
	log.Info("About to download video...please wait")
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

func getConnection() *sql.DB {
	var err error
	var db *sql.DB
	//connStr := "postgres://postgres:password@localhost/retrieveTest?sslmode=disable"
	pg_con := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", pg_con)

	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	// this will be printed in the terminal, confirming the connection to the database
	fmt.Println("Connected to database")
	return db
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welecome to my api")
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
	myRouter.HandleFunc("/", ListOfAllVideos).Methods("GET")
	myRouter.HandleFunc("/trending", Trending).Methods("GET")
	myRouter.HandleFunc("/shorts", YoutubeShorts).Methods("GET")
	myRouter.HandleFunc("/feed/subscriptions/{channelId}", Subscriptions).Methods("GET")
	myRouter.HandleFunc("/feed/libray/watchlater/{channelId}", WatchLater).Methods("GET")
	myRouter.HandleFunc("/feed/libray/likedvideos/{channelId}", LikedVideos).Methods("GET")
	myRouter.HandleFunc("/feed/history/{channelId}", History).Methods("GET")
	myRouter.HandleFunc("/feed/playlist/{channelId}/{playlistId}", Playlist).Methods("GET")
	myRouter.HandleFunc("/feed/yourvideos/{channelId}", YourVideos).Methods("GET")
	myRouter.HandleFunc("/feed/subscriptionlist/{channelId}", SubscriptionList).Methods("GET")
	myRouter.HandleFunc("/{channelId}/listplaylist", Playlist).Methods("GET")
	myRouter.HandleFunc("/{channelId}/channellist", Channellist).Methods("GET")
	myRouter.HandleFunc("/{channelId}/videos", ListChannelVideos).Methods("GET")
	myRouter.HandleFunc("/video/addTrending/{channelId}", AddTrending).Methods("POST")
	myRouter.HandleFunc("/video/addHistory/{channelId}", AddHistory).Methods("POST")
	myRouter.HandleFunc("/video/addShorts/{channelId}", AddShorts).Methods("POST")
	myRouter.HandleFunc("/video/addSubscriptions/{channelId}", AddSubscriptions).Methods("POST")
	myRouter.HandleFunc("/video/addWatchLater/{channelId}", AddWatchLater).Methods("POST")
	myRouter.HandleFunc("/video/addLikedVideos/{channelId}", AddLikedVideos).Methods("POST")
	myRouter.HandleFunc("/video/addPlaylist/{channelId}", AddPlaylist).Methods("POST")
	myRouter.HandleFunc("/video/addSubscriptionList/{channelId}", AddSubscriptionList).Methods("POST")
	myRouter.HandleFunc("/video/addChannelList/{channelId}", AddChannelList).Methods("POST")
	myRouter.HandleFunc("/video/addYourVideos/{channelId}", AddYourVideos).Methods("POST")
	myRouter.HandleFunc("/video/addListPlayList/{channelId}", AddListPlaylist).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", myRouter))
	log.Info("Server started")

	log.Info("Server stopped")
}

func main() {
	initServer()

	log.Infof("Application closed")
}
