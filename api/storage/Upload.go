package storage

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/trevorsibanda/myhustlezw/api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//UploadConfig stores information about the file upload
type UploadConfig struct {
	Owner         primitive.ObjectID
	ID            primitive.ObjectID
	Ext           string
	AllowedMimes  []string
	DeleteFirst   bool
	MaxUploadTime time.Duration
}

//HandleUploadToS3 performs a file upload to an aws s3 bucket then updates
//the file in the database
func HandleUploadToS3(creator model.User, file model.CreatorFile) (newFilename string, eTag *string, err error) {

	reader, err := os.Open(LocalStoragePath(file.Filename))
	if err != nil {
		return
	}
	newFilename = fmt.Sprintf("%s/%s%s", creator.Username, file.ID.Hex(), filepath.Ext(file.Filename))
	object := s3.PutObjectInput{
		Bucket: aws.String(awsS3Bucket),
		Key:    aws.String(newFilename),
		Body:   reader,
		ACL:    aws.String("public-read"),
	}
	var output *s3.PutObjectOutput

	output, err = awsS3Client.PutObject(&object)
	if err != nil {
		fmt.Println("error uploading ", err.Error())
		return
	}

	eTag = output.ETag

	return
}

//JustUploadFile just uploads a file to s3 bucket
func JustUploadFile(owner string, id primitive.ObjectID, source string) (newFilename string, err error) {
	reader, err := os.Open(source)
	if err != nil {
		return
	}
	newFilename = fmt.Sprintf("%s/%s%s", owner, id.Hex(), filepath.Ext(source))
	fmt.Printf("\nuploading %s to %s", source, newFilename)
	object := s3.PutObjectInput{
		Bucket: aws.String(awsS3Bucket),
		Key:    aws.String(newFilename),
		Body:   reader,
		ACL:    aws.String("public-read"),
	}
	var output *s3.PutObjectOutput

	output, err = awsS3Client.PutObject(&object)
	if err != nil {
		fmt.Println("error uploading ", err.Error(), output)
		return
	}
	return
}

//GenerateUploadToS3URL generates a presigned url to allow a file to be uploaded directly to s3 bucket
func GenerateUploadToS3URL(uploadConfig UploadConfig) (uploadURL string, headers map[string]string, err error) {
	input := &s3.ListObjectsInput{
		Bucket: aws.String("myhustlezw"),
	}

	objects, err := awsS3Client.ListObjects(input)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, obj := range objects.Contents {
		fmt.Println(aws.StringValue(obj.Key))
	}

	req, _ := awsS3Client.PutObjectRequest(&s3.PutObjectInput{
		ACL:    aws.String("public-read"),
		Bucket: aws.String(awsS3Bucket),
		Key:    aws.String(fmt.Sprintf("content/%s/%s.%s", uploadConfig.Owner.Hex(), uploadConfig.ID.Hex(), uploadConfig.Ext)),
	})
	uploadURL, err = req.Presign(uploadConfig.MaxUploadTime)
	if err != nil {
		return
	}
	urlObj, err := url.Parse(uploadURL)
	headers = map[string]string{}
	for key, fields := range urlObj.Query() {
		if len(fields) > 0 {
			headers[key] = fields[0]
		}
	}
	return
}
