package storage

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/trevorsibanda/myhustlezw/api/model"
)

var (
	awsS3Endpoint string
	awsS3Bucket   string

	awsS3Config *aws.Config
	awsSession  *session.Session
	awsS3Client *s3.S3

	imgproxySecret   string
	imgproxySalt     string
	imgproxyEndpoint string
	//LocalStorageDir storage directory
	LocalStorageDir string
)

//LocalStoragePath returns the default upload path
func LocalStoragePath(path string) string {
	return fmt.Sprintf("%s/%s", LocalStorageDir, path)
}

//InitStorageEngine initializes the storage engine
func InitStorageEngine() {
	awsS3Endpoint = os.Getenv("S3_ENDPOINT")
	awsS3Bucket = os.Getenv("S3_BUCKET")
	awsS3Key := os.Getenv("S3_KEY")
	awsS3Secret := os.Getenv("S3_SECRET")

	LocalStorageDir = os.Getenv("UPLOADS_STORAGE_DIR")
	if len(LocalStorageDir) == 0 {
		//lets hope we dont crash
		LocalStorageDir = "/tmp"
	}

	awsS3Config = &aws.Config{
		Credentials: credentials.NewStaticCredentials(awsS3Key, awsS3Secret, ""),
		Endpoint:    aws.String(awsS3Endpoint),
		Region:      aws.String("us-east-1"),
	}

	awsSession = session.New(awsS3Config)
	awsS3Client = s3.New(awsSession)

	imgproxySecret = os.Getenv("IMGPROXY_SECRET")
	imgproxySalt = os.Getenv("IMGPROXY_SALT")
	imgproxyEndpoint = os.Getenv("IMGPROXY_ENDPOINT")

	model.RegisterPreviewURLGenerator(FilePreviewImageURL)

}
