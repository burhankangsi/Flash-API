package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func ListObjectsFromS3History(channelId string) ([]string, error) {
	os.Setenv("AWS_ACCESS_KEY", "AKIAVX37IPHMG6BBDYWJ")
	os.Setenv("AWS_SECRET_KEY", "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")

	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	bucket := "images-yc"

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("thumbnails/history/" + channelId),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Infof("Failed to list s3 objects %v", err)
		return nil, err
	}
	var images []string
	var count int
	for _, key := range resp.Contents {
		fmt.Println(*key.Key)
		imageUrl := fmt.Sprintf("https://%v.s3.amazonaws.com/%v", bucket, *key.Key)
		if count > 0 {
			images = append(images, imageUrl)
		}
		count++
	}
	return images, nil

}

func ListObjectsFromS3Subscriptions(channelId string) ([]string, error) {
	os.Setenv("AWS_ACCESS_KEY", "AKIAVX37IPHMG6BBDYWJ")
	os.Setenv("AWS_SECRET_KEY", "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")

	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	bucket := "images-yc"

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("thumbnails/history/" + channelId),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Infof("Failed to list s3 objects %v", err)
		return nil, err
	}
	var images []string
	var count int
	for _, key := range resp.Contents {
		fmt.Println(*key.Key)
		imageUrl := fmt.Sprintf("https://%v.s3.amazonaws.com/%v", bucket, *key.Key)
		if count > 0 {
			images = append(images, imageUrl)
		}
		count++
	}
	return images, nil

}

func ListObjectsFromS3Liked(channelId string) ([]string, error) {
	os.Setenv("AWS_ACCESS_KEY", "AKIAVX37IPHMG6BBDYWJ")
	os.Setenv("AWS_SECRET_KEY", "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")

	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	bucket := "images-yc"

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("thumbnails/history/" + channelId),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Infof("Failed to list s3 objects %v", err)
		return nil, err
	}
	var images []string
	var count int
	for _, key := range resp.Contents {
		fmt.Println(*key.Key)
		imageUrl := fmt.Sprintf("https://%v.s3.amazonaws.com/%v", bucket, *key.Key)
		if count > 0 {
			images = append(images, imageUrl)
		}
		count++
	}
	return images, nil

}

func ListObjectsFromS3Trending(channelId string) ([]string, error) {
	os.Setenv("AWS_ACCESS_KEY", "AKIAVX37IPHMG6BBDYWJ")
	os.Setenv("AWS_SECRET_KEY", "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")

	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	bucket := "images-yc"

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("thumbnails/history/" + channelId),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Infof("Failed to list s3 objects %v", err)
		return nil, err
	}
	var images []string
	var count int
	for _, key := range resp.Contents {
		fmt.Println(*key.Key)
		imageUrl := fmt.Sprintf("https://%v.s3.amazonaws.com/%v", bucket, *key.Key)
		if count > 0 {
			images = append(images, imageUrl)
		}
		count++
	}
	return images, nil

}

func ListObjectsFromS3WatchLater(channelId string) ([]string, error) {
	os.Setenv("AWS_ACCESS_KEY", "AKIAVX37IPHMG6BBDYWJ")
	os.Setenv("AWS_SECRET_KEY", "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")

	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	bucket := "images-yc"

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("thumbnails/history/" + channelId),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Infof("Failed to list s3 objects %v", err)
		return nil, err
	}
	var images []string
	var count int
	for _, key := range resp.Contents {
		fmt.Println(*key.Key)
		imageUrl := fmt.Sprintf("https://%v.s3.amazonaws.com/%v", bucket, *key.Key)
		if count > 0 {
			images = append(images, imageUrl)
		}
		count++
	}
	return images, nil

}

func ListObjectsFromS3ListChannel(channelId string) ([]string, error) {
	os.Setenv("AWS_ACCESS_KEY", "AKIAVX37IPHMG6BBDYWJ")
	os.Setenv("AWS_SECRET_KEY", "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")

	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	bucket := "images-yc"

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("thumbnails/history/" + channelId),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Infof("Failed to list s3 objects %v", err)
		return nil, err
	}
	var images []string
	var count int
	for _, key := range resp.Contents {
		fmt.Println(*key.Key)
		imageUrl := fmt.Sprintf("https://%v.s3.amazonaws.com/%v", bucket, *key.Key)
		if count > 0 {
			images = append(images, imageUrl)
		}
		count++
	}
	return images, nil

}

func ListObjectsFromS3ChannelList(channelId string) ([]string, error) {
	os.Setenv("AWS_ACCESS_KEY", "AKIAVX37IPHMG6BBDYWJ")
	os.Setenv("AWS_SECRET_KEY", "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")

	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	bucket := "images-yc"

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("thumbnails/history/" + channelId),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Infof("Failed to list s3 objects %v", err)
		return nil, err
	}
	var images []string
	var count int
	for _, key := range resp.Contents {
		fmt.Println(*key.Key)
		imageUrl := fmt.Sprintf("https://%v.s3.amazonaws.com/%v", bucket, *key.Key)
		if count > 0 {
			images = append(images, imageUrl)
		}
		count++
	}
	return images, nil

}

func ListObjectsFromS3SubscriptionList(channelId string) ([]string, error) {
	os.Setenv("AWS_ACCESS_KEY", "AKIAVX37IPHMG6BBDYWJ")
	os.Setenv("AWS_SECRET_KEY", "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")

	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	bucket := "images-yc"

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("thumbnails/history/" + channelId),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Infof("Failed to list s3 objects %v", err)
		return nil, err
	}
	var images []string
	var count int
	for _, key := range resp.Contents {
		fmt.Println(*key.Key)
		imageUrl := fmt.Sprintf("https://%v.s3.amazonaws.com/%v", bucket, *key.Key)
		if count > 0 {
			images = append(images, imageUrl)
		}
		count++
	}
	return images, nil

}

func ListObjectsFromS3YourVideos(channelId string) ([]string, error) {
	os.Setenv("AWS_ACCESS_KEY", "AKIAVX37IPHMG6BBDYWJ")
	os.Setenv("AWS_SECRET_KEY", "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")

	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	bucket := "images-yc"

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("thumbnails/history/" + channelId),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Infof("Failed to list s3 objects %v", err)
		return nil, err
	}
	var images []string
	var count int
	for _, key := range resp.Contents {
		fmt.Println(*key.Key)
		imageUrl := fmt.Sprintf("https://%v.s3.amazonaws.com/%v", bucket, *key.Key)
		if count > 0 {
			images = append(images, imageUrl)
		}
		count++
	}
	return images, nil

}

func ListObjectsFromS3Playlist(channelId string) ([]string, error) {
	os.Setenv("AWS_ACCESS_KEY", "AKIAVX37IPHMG6BBDYWJ")
	os.Setenv("AWS_SECRET_KEY", "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")

	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	bucket := "images-yc"

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("thumbnails/history/" + channelId),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Infof("Failed to list s3 objects %v", err)
		return nil, err
	}
	var images []string
	var count int
	for _, key := range resp.Contents {
		fmt.Println(*key.Key)
		imageUrl := fmt.Sprintf("https://%v.s3.amazonaws.com/%v", bucket, *key.Key)
		if count > 0 {
			images = append(images, imageUrl)
		}
		count++
	}
	return images, nil

}

func ListObjectsFromS3Shorts(channelId string) ([]string, error) {
	os.Setenv("AWS_ACCESS_KEY", "AKIAVX37IPHMG6BBDYWJ")
	os.Setenv("AWS_SECRET_KEY", "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")

	svc := s3.New(session.New(), &aws.Config{
		Region: aws.String("us-east-1"),
	})
	bucket := "images-yc"

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String("thumbnails/history/" + channelId),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Infof("Failed to list s3 objects %v", err)
		return nil, err
	}
	var images []string
	var count int
	for _, key := range resp.Contents {
		fmt.Println(*key.Key)
		imageUrl := fmt.Sprintf("https://%v.s3.amazonaws.com/%v", bucket, *key.Key)
		if count > 0 {
			images = append(images, imageUrl)
		}
		count++
	}
	return images, nil

}
