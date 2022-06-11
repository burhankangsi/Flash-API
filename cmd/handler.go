package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

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
