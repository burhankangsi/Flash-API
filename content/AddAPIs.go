package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

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
