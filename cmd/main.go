package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

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

func WatchLater(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos will be displayed here")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	userChannelId, ok := params["channelId"]
	if !ok {
		log.Errorf("User channel ID is missing in parameters")
	}
	log.Infof("user channelId is %v", userChannelId)

	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, thumbnail, upload_date, upload_time FROM videos.history where channel_id=$1 order by upload_time desc", userChannelId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var videoArr []Video

	for rows.Next() {
		var feed Video
		err := rows.Scan(&feed.videoID, &feed.videoName, &feed.channelID, &feed.duration, &feed.title, &feed.channelImage, &feed.views, &feed.timestamp, &feed.channelName, &feed.thumbnail, &feed.uploadDate, &feed.uploadTime)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n", feed)
		videoArr = append(videoArr, feed)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Infof("history videos for user %v are %v", userChannelId, videoArr)
	// images, err := ListObjectsFromS3History(userChannelId)
	// if err != nil {
	// 	panic(err)
	// }

	var response = JsonResponse{Type: "success", Data: videoArr}
	log.Infof("response is %v", response)
	json.NewEncoder(w).Encode(response)
}

func ListOfVideos(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos will be displayed here")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	userChannelId, ok := params["channelId"]
	if !ok {
		log.Errorf("User channel ID is missing in parameters")
	}
	log.Infof("user channelId is %v", userChannelId)

	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, thumbnail, upload_date, upload_time FROM videos.history where channel_id=$1 order by upload_time desc", userChannelId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var videoArr []Video

	for rows.Next() {
		var feed Video
		err := rows.Scan(&feed.videoID, &feed.videoName, &feed.channelID, &feed.duration, &feed.title, &feed.channelImage, &feed.views, &feed.timestamp, &feed.channelName, &feed.thumbnail, &feed.uploadDate, &feed.uploadTime)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n", feed)
		videoArr = append(videoArr, feed)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Infof("history videos for user %v are %v", userChannelId, videoArr)
	// images, err := ListObjectsFromS3History(userChannelId)
	// if err != nil {
	// 	panic(err)
	// }

	var response = JsonResponse{Type: "success", Data: videoArr}
	log.Infof("response is %v", response)
	json.NewEncoder(w).Encode(response)
}

func LikedVideos(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos will be displayed here")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	userChannelId, ok := params["channelId"]
	if !ok {
		log.Errorf("User channel ID is missing in parameters")
	}
	log.Infof("user channelId is %v", userChannelId)

	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, thumbnail, upload_date, upload_time FROM videos.history where channel_id=$1 order by upload_time desc", userChannelId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var videoArr []Video

	for rows.Next() {
		var feed Video
		err := rows.Scan(&feed.videoID, &feed.videoName, &feed.channelID, &feed.duration, &feed.title, &feed.channelImage, &feed.views, &feed.timestamp, &feed.channelName, &feed.thumbnail, &feed.uploadDate, &feed.uploadTime)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n", feed)
		videoArr = append(videoArr, feed)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Infof("history videos for user %v are %v", userChannelId, videoArr)
	// images, err := ListObjectsFromS3History(userChannelId)
	// if err != nil {
	// 	panic(err)
	// }

	var response = JsonResponse{Type: "success", Data: videoArr}
	log.Infof("response is %v", response)
	json.NewEncoder(w).Encode(response)
}

func History(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos will be displayed here")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	userChannelId, ok := params["channelId"]
	if !ok {
		log.Errorf("User channel ID is missing in parameters")
	}
	log.Infof("user channelId is %v", userChannelId)

	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, thumbnail, upload_date, upload_time FROM videos.history where channel_id=$1 order by upload_time desc", userChannelId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var videoArr []Video

	for rows.Next() {
		var feed Video
		err := rows.Scan(&feed.videoID, &feed.videoName, &feed.channelID, &feed.duration, &feed.title, &feed.channelImage, &feed.views, &feed.timestamp, &feed.channelName, &feed.thumbnail, &feed.uploadDate, &feed.uploadTime)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n", feed)
		videoArr = append(videoArr, feed)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Infof("history videos for user %v are %v", userChannelId, videoArr)
	// images, err := ListObjectsFromS3History(userChannelId)
	// if err != nil {
	// 	panic(err)
	// }

	var response = JsonResponse{Type: "success", Data: videoArr}
	log.Infof("response is %v", response)
	json.NewEncoder(w).Encode(response)
}

func ListOfAllVideos(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos will be displayed here")
}

func Subscriptions(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "List of subscriptions videos will be displayed here")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	userChannelId, ok := params["channelId"]
	if !ok {
		log.Errorf("User channel ID is missing in parameters")
	}
	log.Infof("user channelId is %v", userChannelId)

	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, thumbnail, upload_date, upload_time FROM videos.history where channel_id=$1 order by upload_time desc", userChannelId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var videoArr []Video

	for rows.Next() {
		var feed Video
		err := rows.Scan(&feed.videoID, &feed.videoName, &feed.channelID, &feed.duration, &feed.title, &feed.channelImage, &feed.views, &feed.timestamp, &feed.channelName, &feed.thumbnail, &feed.uploadDate, &feed.uploadTime)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n", feed)
		videoArr = append(videoArr, feed)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Infof("history videos for user %v are %v", userChannelId, videoArr)
	// images, err := ListObjectsFromS3Subscriptions(userChannelId)
	// if err != nil {
	// 	panic(err)
	// }

	var response = JsonResponse{Type: "success", Data: videoArr}
	log.Infof("response is %v", response)
	json.NewEncoder(w).Encode(response)
}

func Trending(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos will be displayed here")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	userChannelId, ok := params["channelId"]
	if !ok {
		log.Errorf("User channel ID is missing in parameters")
	}
	log.Infof("user channelId is %v", userChannelId)

	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, thumbnail, upload_date, upload_time FROM videos.history where channel_id=$1 order by upload_time desc", userChannelId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var videoArr []Video

	for rows.Next() {
		var feed Video
		err := rows.Scan(&feed.videoID, &feed.videoName, &feed.channelID, &feed.duration, &feed.title, &feed.channelImage, &feed.views, &feed.timestamp, &feed.channelName, &feed.thumbnail, &feed.uploadDate, &feed.uploadTime)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n", feed)
		videoArr = append(videoArr, feed)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Infof("history videos for user %v are %v", userChannelId, videoArr)
	// images, err := ListObjectsFromS3Trending(userChannelId)
	// if err != nil {
	// 	panic(err)
	// }

	var response = JsonResponse{Type: "success", Data: videoArr}
	log.Infof("response is %v", response)
	json.NewEncoder(w).Encode(response)
}

func Recommended(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of recommended videos will be displayed here")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	userChannelId, ok := params["channelId"]
	if !ok {
		log.Errorf("User channel ID is missing in parameters")
	}
	log.Infof("user channelId is %v", userChannelId)

	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, thumbnail, upload_date, upload_time FROM videos.history where channel_id=$1 order by upload_time desc", userChannelId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var videoArr []Video

	for rows.Next() {
		var feed Video
		err := rows.Scan(&feed.videoID, &feed.videoName, &feed.channelID, &feed.duration, &feed.title, &feed.channelImage, &feed.views, &feed.timestamp, &feed.channelName, &feed.thumbnail, &feed.uploadDate, &feed.uploadTime)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n", feed)
		videoArr = append(videoArr, feed)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Infof("recommended videos for user %v are %v", userChannelId, videoArr)
	// images, err := ListObjectsFromS3Trending(userChannelId)
	// if err != nil {
	// 	panic(err)
	// }

	var response = JsonResponse{Type: "success", Data: videoArr}
	log.Infof("response is %v", response)
	json.NewEncoder(w).Encode(response)
}

func YoutubeShorts(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "List of videos will be displayed here")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	userChannelId, ok := params["channelId"]
	if !ok {
		log.Errorf("User channel ID is missing in parameters")
	}
	log.Infof("user channelId is %v", userChannelId)

	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, thumbnail, upload_date, upload_time FROM videos.history where channel_id=$1 order by upload_time desc", userChannelId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var videoArr []Video

	for rows.Next() {
		var feed Video
		err := rows.Scan(&feed.videoID, &feed.videoName, &feed.channelID, &feed.duration, &feed.title, &feed.channelImage, &feed.views, &feed.timestamp, &feed.channelName, &feed.thumbnail, &feed.uploadDate, &feed.uploadTime)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n", feed)
		videoArr = append(videoArr, feed)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Infof("history videos for user %v are %v", userChannelId, videoArr)
	// images, err := ListObjectsFromS3Shorts(userChannelId)
	// if err != nil {
	// 	panic(err)
	// }

	var response = JsonResponse{Type: "success", Data: videoArr}
	log.Infof("response is %v", response)
	json.NewEncoder(w).Encode(response)
}

func Playlist(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "List of videos will be displayed here")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	userChannelId, ok := params["channelId"]
	if !ok {
		log.Errorf("User channel ID is missing in parameters")
	}
	log.Infof("user channelId is %v", userChannelId)

	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, thumbnail, upload_date, upload_time FROM videos.history where channel_id=$1 order by upload_time desc", userChannelId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var videoArr []Video

	for rows.Next() {
		var feed Video
		err := rows.Scan(&feed.videoID, &feed.videoName, &feed.channelID, &feed.duration, &feed.title, &feed.channelImage, &feed.views, &feed.timestamp, &feed.channelName, &feed.thumbnail, &feed.uploadDate, &feed.uploadTime)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n", feed)
		videoArr = append(videoArr, feed)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Infof("history videos for user %v are %v", userChannelId, videoArr)
	// images, err := ListObjectsFromS3Playlist(userChannelId)
	// if err != nil {
	// 	panic(err)
	// }

	var response = JsonResponse{Type: "success", Data: videoArr}
	log.Infof("response is %v", response)
	json.NewEncoder(w).Encode(response)
}

func YourVideos(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos will be displayed here")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	userChannelId, ok := params["channelId"]
	if !ok {
		log.Errorf("User channel ID is missing in parameters")
	}
	log.Infof("user channelId is %v", userChannelId)

	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, thumbnail, upload_date, upload_time FROM videos.history where channel_id=$1 order by upload_time desc", userChannelId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var videoArr []Video

	for rows.Next() {
		var feed Video
		err := rows.Scan(&feed.videoID, &feed.videoName, &feed.channelID, &feed.duration, &feed.title, &feed.channelImage, &feed.views, &feed.timestamp, &feed.channelName, &feed.thumbnail, &feed.uploadDate, &feed.uploadTime)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n", feed)
		videoArr = append(videoArr, feed)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Infof("history videos for user %v are %v", userChannelId, videoArr)
	// images, err := ListObjectsFromS3YourVideos(userChannelId)
	// if err != nil {
	// 	panic(err)
	// }

	var response = JsonResponse{Type: "success", Data: videoArr}
	log.Infof("response is %v", response)
	json.NewEncoder(w).Encode(response)
}

func SubscriptionList(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos will be displayed here")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	userChannelId, ok := params["channelId"]
	if !ok {
		log.Errorf("User channel ID is missing in parameters")
	}
	log.Infof("user channelId is %v", userChannelId)

	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, thumbnail, upload_date, upload_time FROM videos.history where channel_id=$1 order by upload_time desc", userChannelId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var videoArr []Video

	for rows.Next() {
		var feed Video
		err := rows.Scan(&feed.videoID, &feed.videoName, &feed.channelID, &feed.duration, &feed.title, &feed.channelImage, &feed.views, &feed.timestamp, &feed.channelName, &feed.thumbnail, &feed.uploadDate, &feed.uploadTime)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n", feed)
		videoArr = append(videoArr, feed)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Infof("history videos for user %v are %v", userChannelId, videoArr)
	// images, err := ListObjectsFromS3SubscriptionList(userChannelId)
	// if err != nil {
	// 	panic(err)
	// }

	var response = JsonResponse{Type: "success", Data: videoArr}
	log.Infof("response is %v", response)
	json.NewEncoder(w).Encode(response)
}

func Channellist(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos will be displayed here")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	userChannelId, ok := params["channelId"]
	if !ok {
		log.Errorf("User channel ID is missing in parameters")
	}
	log.Infof("user channelId is %v", userChannelId)

	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, thumbnail, upload_date, upload_time FROM videos.history where channel_id=$1 order by upload_time desc", userChannelId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var videoArr []Video

	for rows.Next() {
		var feed Video
		err := rows.Scan(&feed.videoID, &feed.videoName, &feed.channelID, &feed.duration, &feed.title, &feed.channelImage, &feed.views, &feed.timestamp, &feed.channelName, &feed.thumbnail, &feed.uploadDate, &feed.uploadTime)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n", feed)
		videoArr = append(videoArr, feed)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Infof("history videos for user %v are %v", userChannelId, videoArr)
	// images, err := ListObjectsFromS3ListChannel(userChannelId)
	// if err != nil {
	// 	panic(err)
	// }

	var response = JsonResponse{Type: "success", Data: videoArr}
	log.Infof("response is %v", response)
	json.NewEncoder(w).Encode(response)
}

func ListChannelVideos(w http.ResponseWriter, r *http.Request) {
	//lists all the videos belonging to a particular channel
	//usually called when a channel is searched in search page
	fmt.Fprintf(w, "List of videos will be displayed here")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	channelId, ok := params["channelId"]
	if !ok {
		log.Errorf("User channel ID is missing in parameters")
	}
	log.Infof("user channelId is %v", channelId)

	db := getConnection()
	defer db.Close()
	rows, err := db.Query("SELECT video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, thumbnail, upload_date, upload_time FROM videos.history where video_channel_id=$1 order by upload_time desc", channelId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var videoArr []Video

	for rows.Next() {
		var feed Video
		err := rows.Scan(&feed.videoID, &feed.videoName, &feed.channelID, &feed.duration, &feed.title, &feed.channelImage, &feed.views, &feed.timestamp, &feed.channelName, &feed.thumbnail, &feed.uploadDate, &feed.uploadTime)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n", feed)
		videoArr = append(videoArr, feed)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	log.Infof("history videos for user %v are %v", channelId, videoArr)
	// images, err := ListObjectsFromS3ListChannel(userChannelId)
	// if err != nil {
	// 	panic(err)
	// }

	var response = JsonResponse{Type: "success", Data: videoArr}
	log.Infof("response is %v", response)
	json.NewEncoder(w).Encode(response)
}

func AddChannelList(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos added to history will be displayed here")

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	userChannelId, ok1 := params["channelId"]
	if !ok1 {
		log.Errorf("Channel ID is missing in parameters")
	}

	image, handler, err := r.FormFile("image")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer image.Close()

	log.Infof("Got the image from post request: %+v\n", handler.Filename)
	tmpfile, err := os.Create("./" + handler.Filename)
	defer tmpfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = uploadImageToS3(handler.Filename)
	if err != nil {
		log.Fatalf("Could not upload image to S3. %v", err)
	}

	log.Info("Preparing to add video")
	videoID := r.FormValue("videoid")
	videoName := r.FormValue("videoname")
	channelID := r.FormValue("channelid")
	duration := r.FormValue("duration")
	title := r.FormValue("title")
	chanImage := r.FormValue("channelImage")
	views := r.FormValue("views")
	timestamp := r.FormValue("timestamp")
	chanName := r.FormValue("channelImage")
	uploadDate := r.FormValue("uploaddate")
	uploadTime := r.FormValue("uploadtime")

	var vidId, chanId int

	if vidId, err = strconv.Atoi(videoID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	if chanId, err = strconv.Atoi(channelID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}
	var length int64
	if length, err = strconv.ParseInt(duration, 10, 64); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	var file Video
	file = Video{
		videoID:      vidId,
		videoName:    videoName,
		duration:     length,
		channelID:    chanId,
		title:        title,
		channelImage: chanImage,
		views:        views,
		timestamp:    timestamp,
		channelName:  chanName,
		uploadDate:   uploadDate,
		uploadTime:   uploadTime,
	}

	log.Infof("Got the json %v", file)
	var response = JsonResponse1{}

	if videoID == "" || channelID == "" {
		response = JsonResponse1{Type: "error", Message: "You are missing one or more important parameter."}
	} else {
		db := getConnection()
		defer db.Close()

		sqlStatement := `INSERT INTO videos.history (video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, upload_date, upload_time, channel_id) VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err := db.Exec(sqlStatement, file.videoID, file.videoName, file.channelID, file.duration, file.title, file.channelImage, file.views, file.timestamp, file.channelName, file.uploadDate, file.uploadTime, userChannelId)
		if err != nil {
			log.Infof("Error while inserting the record into history table %v", err)
			panic(err)
		}
		fmt.Fprintf(w, "Record Inserted: ")
		response = JsonResponse1{Type: "success", Message: "The video for history has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(200)
}

func AddHistory(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos added to history will be displayed here")

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	userChannelId, ok1 := params["channelId"]
	if !ok1 {
		log.Errorf("Channel ID is missing in parameters")
	}

	image, handler, err := r.FormFile("image")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer image.Close()

	log.Infof("Got the image from post request: %+v\n", handler.Filename)
	tmpfile, err := os.Create("./" + handler.Filename)
	defer tmpfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = uploadImageToS3(handler.Filename)
	if err != nil {
		log.Fatalf("Could not upload image to S3. %v", err)
	}

	log.Info("Preparing to add video")
	videoID := r.FormValue("videoid")
	videoName := r.FormValue("videoname")
	channelID := r.FormValue("channelid")
	duration := r.FormValue("duration")
	title := r.FormValue("title")
	chanImage := r.FormValue("channelImage")
	views := r.FormValue("views")
	timestamp := r.FormValue("timestamp")
	chanName := r.FormValue("channelImage")
	uploadDate := r.FormValue("uploaddate")
	uploadTime := r.FormValue("uploadtime")

	var vidId, chanId int
	if vidId, err = strconv.Atoi(videoID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	if chanId, err = strconv.Atoi(channelID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}
	var length int64
	if length, err = strconv.ParseInt(duration, 10, 64); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	var file Video
	file = Video{
		videoID:      vidId,
		videoName:    videoName,
		duration:     length,
		channelID:    chanId,
		title:        title,
		channelImage: chanImage,
		views:        views,
		timestamp:    timestamp,
		channelName:  chanName,
		uploadDate:   uploadDate,
		uploadTime:   uploadTime,
	}

	log.Infof("Got the json %v", file)
	var response = JsonResponse1{}

	if videoID == "" || channelID == "" {
		response = JsonResponse1{Type: "error", Message: "You are missing one or more important parameter."}
	} else {
		db := getConnection()
		defer db.Close()

		sqlStatement := `INSERT INTO videos.history (video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, upload_date, upload_time, channel_id) VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err := db.Exec(sqlStatement, file.videoID, file.videoName, file.channelID, file.duration, file.title, file.channelImage, file.views, file.timestamp, file.channelName, file.uploadDate, file.uploadTime, userChannelId)
		if err != nil {
			log.Infof("Error while inserting the record into history table %v", err)
			panic(err)
		}
		fmt.Fprintf(w, "Record Inserted: ")
		response = JsonResponse1{Type: "success", Message: "The video for history has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

func AddShorts(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos added to history will be displayed here")

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	userChannelId, ok1 := params["channelId"]
	if !ok1 {
		log.Errorf("Channel ID is missing in parameters")
	}
	image, handler, err := r.FormFile("image")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer image.Close()

	log.Infof("Got the image from post request: %+v\n", handler.Filename)
	tmpfile, err := os.Create("./" + handler.Filename)
	defer tmpfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = uploadImageToS3(handler.Filename)
	if err != nil {
		log.Fatalf("Could not upload image to S3. %v", err)
	}

	log.Info("Preparing to add video")
	videoID := r.FormValue("videoid")
	videoName := r.FormValue("videoname")
	channelID := r.FormValue("channelid")
	duration := r.FormValue("duration")
	title := r.FormValue("title")
	chanImage := r.FormValue("channelImage")
	views := r.FormValue("views")
	timestamp := r.FormValue("timestamp")
	chanName := r.FormValue("channelImage")
	uploadDate := r.FormValue("uploaddate")
	uploadTime := r.FormValue("uploadtime")

	var vidId, chanId int
	if vidId, err = strconv.Atoi(videoID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	if chanId, err = strconv.Atoi(channelID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}
	var length int64
	if length, err = strconv.ParseInt(duration, 10, 64); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	var file Video
	file = Video{
		videoID:      vidId,
		videoName:    videoName,
		duration:     length,
		channelID:    chanId,
		title:        title,
		channelImage: chanImage,
		views:        views,
		timestamp:    timestamp,
		channelName:  chanName,
		uploadDate:   uploadDate,
		uploadTime:   uploadTime,
	}

	log.Infof("Got the json %v", file)
	var response = JsonResponse1{}

	if videoID == "" || channelID == "" {
		response = JsonResponse1{Type: "error", Message: "You are missing one or more important parameter."}
	} else {
		db := getConnection()
		defer db.Close()

		sqlStatement := `INSERT INTO videos.history (video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, upload_date, upload_time, channel_id) VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err := db.Exec(sqlStatement, file.videoID, file.videoName, file.channelID, file.duration, file.title, file.channelImage, file.views, file.timestamp, file.channelName, file.uploadDate, file.uploadTime, userChannelId)
		if err != nil {
			log.Infof("Error while inserting the record into history table %v", err)
			panic(err)
		}
		fmt.Fprintf(w, "Record Inserted: ")
		response = JsonResponse1{Type: "success", Message: "The video for history has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

func AddLikedVideos(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos added to history will be displayed here")

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	userChannelId, ok1 := params["channelId"]
	if !ok1 {
		log.Errorf("Channel ID is missing in parameters")
	}
	image, handler, err := r.FormFile("image")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer image.Close()

	log.Infof("Got the image from post request: %+v\n", handler.Filename)
	tmpfile, err := os.Create("./" + handler.Filename)
	defer tmpfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = uploadImageToS3(handler.Filename)
	if err != nil {
		log.Fatalf("Could not upload image to S3. %v", err)
	}

	log.Info("Preparing to add video")
	videoID := r.FormValue("videoid")
	videoName := r.FormValue("videoname")
	channelID := r.FormValue("channelid")
	duration := r.FormValue("duration")
	title := r.FormValue("title")
	chanImage := r.FormValue("channelImage")
	views := r.FormValue("views")
	timestamp := r.FormValue("timestamp")
	chanName := r.FormValue("channelImage")
	uploadDate := r.FormValue("uploaddate")
	uploadTime := r.FormValue("uploadtime")

	var vidId, chanId int
	if vidId, err = strconv.Atoi(videoID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	if chanId, err = strconv.Atoi(channelID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}
	var length int64
	if length, err = strconv.ParseInt(duration, 10, 64); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	var file Video
	file = Video{
		videoID:      vidId,
		videoName:    videoName,
		duration:     length,
		channelID:    chanId,
		title:        title,
		channelImage: chanImage,
		views:        views,
		timestamp:    timestamp,
		channelName:  chanName,
		uploadDate:   uploadDate,
		uploadTime:   uploadTime,
	}

	log.Infof("Got the json %v", file)
	var response = JsonResponse1{}

	if videoID == "" || channelID == "" {
		response = JsonResponse1{Type: "error", Message: "You are missing one or more important parameter."}
	} else {
		db := getConnection()
		defer db.Close()

		sqlStatement := `INSERT INTO videos.history (video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, upload_date, upload_time, channel_id) VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err := db.Exec(sqlStatement, file.videoID, file.videoName, file.channelID, file.duration, file.title, file.channelImage, file.views, file.timestamp, file.channelName, file.uploadDate, file.uploadTime, userChannelId)
		if err != nil {
			log.Infof("Error while inserting the record into history table %v", err)
			panic(err)
		}
		fmt.Fprintf(w, "Record Inserted: ")
		response = JsonResponse1{Type: "success", Message: "The video for history has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

func AddSubscriptions(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos added to history will be displayed here")

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	userChannelId, ok1 := params["channelId"]
	if !ok1 {
		log.Errorf("Channel ID is missing in parameters")
	}
	image, handler, err := r.FormFile("image")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer image.Close()

	log.Infof("Got the image from post request: %+v\n", handler.Filename)
	tmpfile, err := os.Create("./" + handler.Filename)
	defer tmpfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = uploadImageToS3(handler.Filename)
	if err != nil {
		log.Fatalf("Could not upload image to S3. %v", err)
	}

	log.Info("Preparing to add video")
	videoID := r.FormValue("videoid")
	videoName := r.FormValue("videoname")
	channelID := r.FormValue("channelid")
	duration := r.FormValue("duration")
	title := r.FormValue("title")
	chanImage := r.FormValue("channelImage")
	views := r.FormValue("views")
	timestamp := r.FormValue("timestamp")
	chanName := r.FormValue("channelImage")
	uploadDate := r.FormValue("uploaddate")
	uploadTime := r.FormValue("uploadtime")

	var vidId, chanId int
	if vidId, err = strconv.Atoi(videoID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	if chanId, err = strconv.Atoi(channelID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}
	var length int64
	if length, err = strconv.ParseInt(duration, 10, 64); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	var file Video
	file = Video{
		videoID:      vidId,
		videoName:    videoName,
		duration:     length,
		channelID:    chanId,
		title:        title,
		channelImage: chanImage,
		views:        views,
		timestamp:    timestamp,
		channelName:  chanName,
		uploadDate:   uploadDate,
		uploadTime:   uploadTime,
	}

	log.Infof("Got the json %v", file)
	var response = JsonResponse1{}

	if videoID == "" || channelID == "" {
		response = JsonResponse1{Type: "error", Message: "You are missing one or more important parameter."}
	} else {
		db := getConnection()
		defer db.Close()

		sqlStatement := `INSERT INTO videos.history (video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, upload_date, upload_time, channel_id) VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err := db.Exec(sqlStatement, file.videoID, file.videoName, file.channelID, file.duration, file.title, file.channelImage, file.views, file.timestamp, file.channelName, file.uploadDate, file.uploadTime, userChannelId)
		if err != nil {
			log.Infof("Error while inserting the record into history table %v", err)
			panic(err)
		}
		fmt.Fprintf(w, "Record Inserted: ")
		response = JsonResponse1{Type: "success", Message: "The video for history has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

func AddTrending(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos added to history will be displayed here")

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	userChannelId, ok1 := params["channelId"]
	if !ok1 {
		log.Errorf("Channel ID is missing in parameters")
	}
	image, handler, err := r.FormFile("image")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer image.Close()

	log.Infof("Got the image from post request: %+v\n", handler.Filename)
	tmpfile, err := os.Create("./" + handler.Filename)
	defer tmpfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = uploadImageToS3(handler.Filename)
	if err != nil {
		log.Fatalf("Could not upload image to S3. %v", err)
	}

	log.Info("Preparing to add video")
	videoID := r.FormValue("videoid")
	videoName := r.FormValue("videoname")
	channelID := r.FormValue("channelid")
	duration := r.FormValue("duration")
	title := r.FormValue("title")
	chanImage := r.FormValue("channelImage")
	views := r.FormValue("views")
	timestamp := r.FormValue("timestamp")
	chanName := r.FormValue("channelImage")
	uploadDate := r.FormValue("uploaddate")
	uploadTime := r.FormValue("uploadtime")

	var vidId, chanId int
	if vidId, err = strconv.Atoi(videoID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	if chanId, err = strconv.Atoi(channelID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}
	var length int64
	if length, err = strconv.ParseInt(duration, 10, 64); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	var file Video
	file = Video{
		videoID:      vidId,
		videoName:    videoName,
		duration:     length,
		channelID:    chanId,
		title:        title,
		channelImage: chanImage,
		views:        views,
		timestamp:    timestamp,
		channelName:  chanName,
		uploadDate:   uploadDate,
		uploadTime:   uploadTime,
	}

	log.Infof("Got the json %v", file)
	var response = JsonResponse1{}

	if videoID == "" || channelID == "" {
		response = JsonResponse1{Type: "error", Message: "You are missing one or more important parameter."}
	} else {
		db := getConnection()
		defer db.Close()

		sqlStatement := `INSERT INTO videos.history (video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, upload_date, upload_time, channel_id) VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err := db.Exec(sqlStatement, file.videoID, file.videoName, file.channelID, file.duration, file.title, file.channelImage, file.views, file.timestamp, file.channelName, file.uploadDate, file.uploadTime, userChannelId)
		if err != nil {
			log.Infof("Error while inserting the record into history table %v", err)
			panic(err)
		}
		fmt.Fprintf(w, "Record Inserted: ")
		response = JsonResponse1{Type: "success", Message: "The video for history has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

func AddRecommended(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos added to recommended will be displayed here")

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	userChannelId, ok1 := params["channelId"]
	if !ok1 {
		log.Errorf("Channel ID is missing in parameters")
	}
	image, handler, err := r.FormFile("image")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer image.Close()

	log.Infof("Got the image from post request: %+v\n", handler.Filename)
	tmpfile, err := os.Create("./" + handler.Filename)
	defer tmpfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = uploadImageToS3(handler.Filename)
	if err != nil {
		log.Fatalf("Could not upload image to S3. %v", err)
	}

	log.Info("Preparing to add video")
	videoID := r.FormValue("videoid")
	videoName := r.FormValue("videoname")
	channelID := r.FormValue("channelid")
	duration := r.FormValue("duration")
	title := r.FormValue("title")
	chanImage := r.FormValue("channelImage")
	views := r.FormValue("views")
	timestamp := r.FormValue("timestamp")
	chanName := r.FormValue("channelImage")
	uploadDate := r.FormValue("uploaddate")
	uploadTime := r.FormValue("uploadtime")

	var vidId, chanId int
	if vidId, err = strconv.Atoi(videoID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	if chanId, err = strconv.Atoi(channelID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}
	var length int64
	if length, err = strconv.ParseInt(duration, 10, 64); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	var file Video
	file = Video{
		videoID:      vidId,
		videoName:    videoName,
		duration:     length,
		channelID:    chanId,
		title:        title,
		channelImage: chanImage,
		views:        views,
		timestamp:    timestamp,
		channelName:  chanName,
		uploadDate:   uploadDate,
		uploadTime:   uploadTime,
	}

	log.Infof("Got the json %v", file)
	var response = JsonResponse1{}

	if videoID == "" || channelID == "" {
		response = JsonResponse1{Type: "error", Message: "You are missing one or more important parameter."}
	} else {
		db := getConnection()
		defer db.Close()

		sqlStatement := `INSERT INTO videos.history (video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, upload_date, upload_time, channel_id) VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err := db.Exec(sqlStatement, file.videoID, file.videoName, file.channelID, file.duration, file.title, file.channelImage, file.views, file.timestamp, file.channelName, file.uploadDate, file.uploadTime, userChannelId)
		if err != nil {
			log.Infof("Error while inserting the record into recommended table %v", err)
			panic(err)
		}
		fmt.Fprintf(w, "Record Inserted: ")
		response = JsonResponse1{Type: "success", Message: "The video for recommended has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

func AddWatchLater(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos added to history will be displayed here")

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	userChannelId, ok1 := params["channelId"]
	if !ok1 {
		log.Errorf("Channel ID is missing in parameters")
	}
	image, handler, err := r.FormFile("image")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer image.Close()

	log.Infof("Got the image from post request: %+v\n", handler.Filename)
	tmpfile, err := os.Create("./" + handler.Filename)
	defer tmpfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = uploadImageToS3(handler.Filename)
	if err != nil {
		log.Fatalf("Could not upload image to S3. %v", err)
	}

	log.Info("Preparing to add video")
	videoID := r.FormValue("videoid")
	videoName := r.FormValue("videoname")
	channelID := r.FormValue("channelid")
	duration := r.FormValue("duration")
	title := r.FormValue("title")
	chanImage := r.FormValue("channelImage")
	views := r.FormValue("views")
	timestamp := r.FormValue("timestamp")
	chanName := r.FormValue("channelImage")
	uploadDate := r.FormValue("uploaddate")
	uploadTime := r.FormValue("uploadtime")

	var vidId, chanId int
	if vidId, err = strconv.Atoi(videoID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	if chanId, err = strconv.Atoi(channelID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}
	var length int64
	if length, err = strconv.ParseInt(duration, 10, 64); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	var file Video
	file = Video{
		videoID:      vidId,
		videoName:    videoName,
		duration:     length,
		channelID:    chanId,
		title:        title,
		channelImage: chanImage,
		views:        views,
		timestamp:    timestamp,
		channelName:  chanName,
		uploadDate:   uploadDate,
		uploadTime:   uploadTime,
	}

	log.Infof("Got the json %v", file)
	var response = JsonResponse1{}

	if videoID == "" || channelID == "" {
		response = JsonResponse1{Type: "error", Message: "You are missing one or more important parameter."}
	} else {
		db := getConnection()
		defer db.Close()

		sqlStatement := `INSERT INTO videos.history (video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, upload_date, upload_time, channel_id) VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err := db.Exec(sqlStatement, file.videoID, file.videoName, file.channelID, file.duration, file.title, file.channelImage, file.views, file.timestamp, file.channelName, file.uploadDate, file.uploadTime, userChannelId)
		if err != nil {
			log.Infof("Error while inserting the record into history table %v", err)
			panic(err)
		}
		fmt.Fprintf(w, "Record Inserted: ")
		response = JsonResponse1{Type: "success", Message: "The video for history has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

func AddSubscriptionList(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos added to history will be displayed here")

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	userChannelId, ok1 := params["channelId"]
	if !ok1 {
		log.Errorf("Channel ID is missing in parameters")
	}
	image, handler, err := r.FormFile("image")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer image.Close()

	log.Infof("Got the image from post request: %+v\n", handler.Filename)
	tmpfile, err := os.Create("./" + handler.Filename)
	defer tmpfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = uploadImageToS3(handler.Filename)
	if err != nil {
		log.Fatalf("Could not upload image to S3. %v", err)
	}

	log.Info("Preparing to add video")
	videoID := r.FormValue("videoid")
	videoName := r.FormValue("videoname")
	channelID := r.FormValue("channelid")
	duration := r.FormValue("duration")
	title := r.FormValue("title")
	chanImage := r.FormValue("channelImage")
	views := r.FormValue("views")
	timestamp := r.FormValue("timestamp")
	chanName := r.FormValue("channelImage")
	uploadDate := r.FormValue("uploaddate")
	uploadTime := r.FormValue("uploadtime")

	var vidId, chanId int
	if vidId, err = strconv.Atoi(videoID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	if chanId, err = strconv.Atoi(channelID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}
	var length int64
	if length, err = strconv.ParseInt(duration, 10, 64); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	var file Video
	file = Video{
		videoID:      vidId,
		videoName:    videoName,
		duration:     length,
		channelID:    chanId,
		title:        title,
		channelImage: chanImage,
		views:        views,
		timestamp:    timestamp,
		channelName:  chanName,
		uploadDate:   uploadDate,
		uploadTime:   uploadTime,
	}

	log.Infof("Got the json %v", file)
	var response = JsonResponse1{}

	if videoID == "" || channelID == "" {
		response = JsonResponse1{Type: "error", Message: "You are missing one or more important parameter."}
	} else {
		db := getConnection()
		defer db.Close()

		sqlStatement := `INSERT INTO videos.history (video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, upload_date, upload_time, channel_id) VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err := db.Exec(sqlStatement, file.videoID, file.videoName, file.channelID, file.duration, file.title, file.channelImage, file.views, file.timestamp, file.channelName, file.uploadDate, file.uploadTime, userChannelId)
		if err != nil {
			log.Infof("Error while inserting the record into history table %v", err)
			panic(err)
		}
		fmt.Fprintf(w, "Record Inserted: ")
		response = JsonResponse1{Type: "success", Message: "The video for history has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

func AddPlaylist(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos added to history will be displayed here")

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	userChannelId, ok1 := params["channelId"]
	if !ok1 {
		log.Errorf("Channel ID is missing in parameters")
	}
	image, handler, err := r.FormFile("image")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer image.Close()

	log.Infof("Got the image from post request: %+v\n", handler.Filename)
	tmpfile, err := os.Create("./" + handler.Filename)
	defer tmpfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = uploadImageToS3(handler.Filename)
	if err != nil {
		log.Fatalf("Could not upload image to S3. %v", err)
	}

	log.Info("Preparing to add video")
	videoID := r.FormValue("videoid")
	videoName := r.FormValue("videoname")
	channelID := r.FormValue("channelid")
	duration := r.FormValue("duration")
	title := r.FormValue("title")
	chanImage := r.FormValue("channelImage")
	views := r.FormValue("views")
	timestamp := r.FormValue("timestamp")
	chanName := r.FormValue("channelImage")
	uploadDate := r.FormValue("uploaddate")
	uploadTime := r.FormValue("uploadtime")

	var vidId, chanId int
	if vidId, err = strconv.Atoi(videoID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	if chanId, err = strconv.Atoi(channelID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}
	var length int64
	if length, err = strconv.ParseInt(duration, 10, 64); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	var file Video
	file = Video{
		videoID:      vidId,
		videoName:    videoName,
		duration:     length,
		channelID:    chanId,
		title:        title,
		channelImage: chanImage,
		views:        views,
		timestamp:    timestamp,
		channelName:  chanName,
		uploadDate:   uploadDate,
		uploadTime:   uploadTime,
	}

	log.Infof("Got the json %v", file)
	var response = JsonResponse1{}

	if videoID == "" || channelID == "" {
		response = JsonResponse1{Type: "error", Message: "You are missing one or more important parameter."}
	} else {
		db := getConnection()
		defer db.Close()

		sqlStatement := `INSERT INTO videos.history (video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, upload_date, upload_time, channel_id) VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err := db.Exec(sqlStatement, file.videoID, file.videoName, file.channelID, file.duration, file.title, file.channelImage, file.views, file.timestamp, file.channelName, file.uploadDate, file.uploadTime, userChannelId)
		if err != nil {
			log.Infof("Error while inserting the record into history table %v", err)
			panic(err)
		}
		fmt.Fprintf(w, "Record Inserted: ")
		response = JsonResponse1{Type: "success", Message: "The video for history has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

func AddYourVideos(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos added to history will be displayed here")

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	userChannelId, ok1 := params["channelId"]
	if !ok1 {
		log.Errorf("Channel ID is missing in parameters")
	}
	image, handler, err := r.FormFile("image")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer image.Close()

	log.Infof("Got the image from post request: %+v\n", handler.Filename)
	tmpfile, err := os.Create("./" + handler.Filename)
	defer tmpfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = uploadImageToS3(handler.Filename)
	if err != nil {
		log.Fatalf("Could not upload image to S3. %v", err)
	}
	log.Info("Preparing to add video")
	videoID := r.FormValue("videoid")
	videoName := r.FormValue("videoname")
	channelID := r.FormValue("channelid")
	duration := r.FormValue("duration")
	title := r.FormValue("title")
	chanImage := r.FormValue("channelImage")
	views := r.FormValue("views")
	timestamp := r.FormValue("timestamp")
	chanName := r.FormValue("channelImage")
	uploadDate := r.FormValue("uploaddate")
	uploadTime := r.FormValue("uploadtime")

	var vidId, chanId int
	if vidId, err = strconv.Atoi(videoID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	if chanId, err = strconv.Atoi(channelID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}
	var length int64
	if length, err = strconv.ParseInt(duration, 10, 64); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	var file Video
	file = Video{
		videoID:      vidId,
		videoName:    videoName,
		duration:     length,
		channelID:    chanId,
		title:        title,
		channelImage: chanImage,
		views:        views,
		timestamp:    timestamp,
		channelName:  chanName,
		uploadDate:   uploadDate,
		uploadTime:   uploadTime,
	}

	log.Infof("Got the json %v", file)
	var response = JsonResponse1{}

	if videoID == "" || channelID == "" {
		response = JsonResponse1{Type: "error", Message: "You are missing one or more important parameter."}
	} else {
		db := getConnection()
		defer db.Close()

		sqlStatement := `INSERT INTO videos.history (video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, upload_date, upload_time, channel_id) VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err := db.Exec(sqlStatement, file.videoID, file.videoName, file.channelID, file.duration, file.title, file.channelImage, file.views, file.timestamp, file.channelName, file.uploadDate, file.uploadTime, userChannelId)
		if err != nil {
			log.Infof("Error while inserting the record into history table %v", err)
			panic(err)
		}
		fmt.Fprintf(w, "Record Inserted: ")
		response = JsonResponse1{Type: "success", Message: "The video for history has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

func AddListPlaylist(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of videos added to history will be displayed here")

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	userChannelId, ok1 := params["channelId"]
	if !ok1 {
		log.Errorf("Channel ID is missing in parameters")
	}
	image, handler, err := r.FormFile("image")
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer image.Close()

	log.Infof("Got the image from post request: %+v\n", handler.Filename)
	tmpfile, err := os.Create("./" + handler.Filename)
	defer tmpfile.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(tmpfile, image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = uploadImageToS3(handler.Filename)
	if err != nil {
		log.Fatalf("Could not upload image to S3. %v", err)
	}

	log.Info("Preparing to add video")
	videoID := r.FormValue("videoid")
	videoName := r.FormValue("videoname")
	channelID := r.FormValue("channelid")
	duration := r.FormValue("duration")
	title := r.FormValue("title")
	chanImage := r.FormValue("channelImage")
	views := r.FormValue("views")
	timestamp := r.FormValue("timestamp")
	chanName := r.FormValue("channelImage")
	uploadDate := r.FormValue("uploaddate")
	uploadTime := r.FormValue("uploadtime")

	var vidId, chanId int
	if vidId, err = strconv.Atoi(videoID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	if chanId, err = strconv.Atoi(channelID); err != nil {
		log.Infof("Error converting to int. %v", err)
	}
	var length int64
	if length, err = strconv.ParseInt(duration, 10, 64); err != nil {
		log.Infof("Error converting to int. %v", err)
	}

	var file Video
	file = Video{
		videoID:      vidId,
		videoName:    videoName,
		duration:     length,
		channelID:    chanId,
		title:        title,
		channelImage: chanImage,
		views:        views,
		timestamp:    timestamp,
		channelName:  chanName,
		uploadDate:   uploadDate,
		uploadTime:   uploadTime,
	}

	log.Infof("Got the json %v", file)
	var response = JsonResponse1{}

	if videoID == "" || channelID == "" {
		response = JsonResponse1{Type: "error", Message: "You are missing one or more important parameter."}
	} else {
		db := getConnection()
		defer db.Close()

		sqlStatement := `INSERT INTO videos.history (video_id, video_name, video_channel_id, length, title, channelImage, views, timestamp, channelName, upload_date, upload_time, channel_id) VALUES ($1,$2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err := db.Exec(sqlStatement, file.videoID, file.videoName, file.channelID, file.duration, file.title, file.channelImage, file.views, file.timestamp, file.channelName, file.uploadDate, file.uploadTime, userChannelId)
		if err != nil {
			log.Infof("Error while inserting the record into history table %v", err)
			panic(err)
		}
		fmt.Fprintf(w, "Record Inserted: ")
		response = JsonResponse1{Type: "success", Message: "The video for history has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

//func initServer(ctx context.Context) (*http.Server, context.Context) {
func InitServer() {
	//Creating the routers
	log.Infof("Starting server")

	myRouter := mux.NewRouter().StrictSlash(true)
	log.Infof("Calling handler")
	myRouter.HandleFunc("/{channelId}/{videoId}/video.ts", GetVideoObjectHandler).Methods("GET")
	myRouter.HandleFunc("/video/{channelId}", UploadVideoHandler).Methods("POST")
	myRouter.HandleFunc("/video", ListOfVideos).Methods("GET")
	myRouter.HandleFunc("/", ListOfAllVideos).Methods("GET")
	myRouter.HandleFunc("/trending", Trending).Methods("GET")
	myRouter.HandleFunc("/recommended", Recommended).Methods("GET")
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
	myRouter.HandleFunc("/video/addRecommended/{channelId}", AddRecommended).Methods("POST")
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

	// log.Fatal(http.ListenAndServe(":8000", myRouter))
	port := os.Getenv("PORT")
	log.Infof("Listening on port: %v", port)
	address := fmt.Sprintf("%s:%s", "0.0.0.0", port)

	log.Fatal(http.ListenAndServe(address, myRouter))
	log.Info("Server started")

	//log.Info("Server stopped")
}

func main() {
	InitServer()

	log.Infof("Application closed")
}
