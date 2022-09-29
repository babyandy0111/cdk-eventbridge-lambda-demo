package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"os"
	"time"
)

var cfg aws.Config
var err error
var detail EventBridge

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(_ context.Context, event events.CloudWatchEvent) error {
	region := os.Getenv("region")
	zones := os.Getenv("zones")
	iam_endpoint := os.Getenv("iam_endpoint")

	fmt.Println("env region:", region)
	fmt.Println("env zones:", zones)
	fmt.Println("env iam_endpoint:", iam_endpoint)

	fmt.Printf("event: %#v\n", event)
	return nil
}
