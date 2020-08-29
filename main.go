package main

import (
	"fmt"
	"github.com/docker/go-plugins-helpers/volume"
	"log"
	"os"
)

func main() {
	// On start always check for arguments and update s3fs credentials
	s3Bucket := os.Getenv("S3_BUCKET")
	s3AccessKey := os.Getenv("S3_ACCESS_KEY")
	s3SecretKey := os.Getenv("S3_SECRET_KEY")
	//s3Endpoint := os.Getenv("S3_ENDPOINT")
	//s3fsOptions := os.Getenv("S3VOL_OPTIONS")

	if s3Bucket == "" || s3AccessKey == "" || s3SecretKey == "" {
		log.Fatal("S3_BUCKET, S3_ACCESS_KEY and S3_SECRET_KEY are required")
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	credFile := fmt.Sprintf("%s/.s3backer_passwd", homedir)

	fd, err := os.OpenFile(credFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fd.WriteString(fmt.Sprintf("%s:%s\n", s3AccessKey, s3SecretKey))
	if err != nil {
		log.Fatal(err)
	}
	err = fd.Sync()
	if err != nil {
		log.Fatal(err)
	}
	err = fd.Close()
	if err != nil {
		log.Fatal(err)
	}

	driver := NewS3Driver()
	handler := volume.NewHandler(driver)
	if err := handler.ServeUnix("s3vol", 0); err != nil {
		log.Fatal(err)
	}
}
