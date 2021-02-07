package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	apiKey     = "apiKey"
	apiToken   = "apiToken"
	boardID    = "B6ZZ3iHJ"
	url        = "https://api.trello.com/1/boards/%s?actions=all&actions_limit=1000&card_attachment_fields=all&cards=all&lists=all&members=all&member_fields=all&card_attachment_fields=all&checklists=all&fields=all&key=%s&token=%s"
	fileName   = "trelloBoardContents-%s-bquach.json"
	bucketName = "waejay-trello-backup"
	awsID      = "awsID"
	awsAccess  = "awsAcess"
)

func main() {
	// Prepare Trello board's contents and write into a file
	writeTrelloBoardContentsFile()

	// Initialize S3Manager client using our specified const credentials
	creds := credentials.NewStaticCredentials(awsID, awsAccess, "")
	config := aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: creds,
	}
	sess := session.New(&config)
	s3Manager := s3manager.NewUploader(sess)

	// Upload contents to S3
	file, err := os.Open("/tmp/" + getFileNameWithDate())
	checkError(err)

	_, err = s3Manager.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filepath.Base("/tmp/" + getFileNameWithDate())),
		Body:   file,
	})
	checkError(err)

	fmt.Println("Uploaded trello contents to S3.")
}

func writeTrelloBoardContentsFile() {
	// Create endpoint
	endpoint := getAPIPath(url, boardID, apiKey, apiToken)

	// Make GET request
	fmt.Printf("Using endpoint:\n%s\n", endpoint)
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	// Extract the data body from the HTTP response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Since this script is running in a Lambda function and not a running host,
	// we store the file in the /tmp folder
	err = ioutil.WriteFile("/tmp/"+getFileNameWithDate(), body, 0644)
	checkError(err)
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func getAPIPath(url, boardID, apiKey, apiToken string) string {
	return fmt.Sprintf(url, boardID, apiKey, apiToken)
}

func getFileNameWithDate() string {
	return fmt.Sprintf(fileName, time.Now().Format("2020-01-01"))
}
