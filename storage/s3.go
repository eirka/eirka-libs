package storage

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"

	"github.com/eirka/eirka-libs/config"
)

type AmazonS3 struct {
	session *session.Session
}

// create aws session with credentials
func (a *AmazonS3) auth() (err error) {

	// new credentials from settings
	creds := credentials.NewStaticCredentials(config.Settings.Amazon.Id, config.Settings.Amazon.Key, "")

	// create our session
	a.session = session.New(&aws.Config{
		Region:      aws.String(config.Settings.Amazon.Region),
		Credentials: creds,
		MaxRetries:  aws.Int(10),
	})

	return

}

// Upload a file to S3
func (a *AmazonS3) Save(filepath, filename, mime string) (err error) {

	err = a.auth()
	if err != nil {
		return
	}

	file, err := os.Open(filepath)
	if err != nil {
		return errors.New("problem opening file for s3")
	}
	defer file.Close()

	uploader := s3manager.NewUploader(a.session)

	params := &s3manager.UploadInput{
		Bucket:               aws.String(config.Settings.Amazon.Bucket),
		Key:                  aws.String(filename),
		Body:                 file,
		ContentType:          aws.String(mime),
		ServerSideEncryption: aws.String(s3.ServerSideEncryptionAes256),
	}

	_, err = uploader.Upload(params)

	return

}

// Delete a file from S3
func (a *AmazonS3) Delete(key string) (err error) {

	err = a.auth()
	if err != nil {
		return
	}

	svc := s3.New(a.session)

	params := &s3.DeleteObjectInput{
		Bucket: aws.String(config.Settings.Amazon.Bucket),
		Key:    aws.String(key),
	}

	_, err = svc.DeleteObject(params)

	return

}
