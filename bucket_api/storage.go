package bucket_api

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type awsCreds struct {
	session *session.Session
}

// func saveToBucket() {

// 	ret := retrieve{}
// 	buffer := new(bytes.Buffer)
// 	buffer.ReadFrom(ret.pipeReader)
// 	res := buffer.Bytes()

// 	// Upload Files
// 	err = uploadFile(session, res)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func createAWSSession() {
	var currAws awsCreds
	session, err := session.NewSession(&aws.Config{Region: aws.String("ap-south-1")})
	if err != nil {
		log.Fatal(err)
	}
	currAws.session = session
}

func uploadFile(video []byte) error {

	var currAws awsCreds
	vid := string(video)
	upFile, err := os.Open(vid)
	if err != nil {
		return err
	}
	defer upFile.Close()

	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	// _, err = s3.New(session).PutObject(&s3.PutObjectInput{
	// 	Bucket:               aws.String(AWS_S3_BUCKET),
	// 	Key:                  aws.String(video),
	// 	ACL:                  aws.String("private"),
	// 	Body:                 bytes.NewReader(fileBuffer),
	// 	ContentLength:        aws.Int64(fileSize),
	// 	ContentType:          aws.String(http.DetectContentType(fileBuffer)),
	// 	ContentDisposition:   aws.String("attachment"),
	// 	ServerSideEncryption: aws.String("AES256"),
	// })
	// return err

	uploader := s3manager.NewUploader(currAws.session)
	AWS_S3_BUCKET := ""

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(AWS_S3_BUCKET), // Bucket to be used
		Key:    aws.String(vid),           // Name of the file to be saved
		Body:   upFile,                    // File
	})
	if err != nil {
		// Do your error handling here
		return err
	}
	return nil
}
