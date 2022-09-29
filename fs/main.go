package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/support"
	"github.com/aws/aws-sdk-go-v2/service/support/types"
	"io/ioutil"
	"log"
	"main/model"
	"net/http"
	"os"
)

var cfg aws.Config
var err error
var detail model.Detail
var endpoint string
var region string
var zones string
var fsKey string

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(_ context.Context, event events.CloudWatchEvent) error {
	cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	region = os.Getenv("region")
	zones = os.Getenv("zones")
	endpoint = os.Getenv("fs_domain")
	fsKey = os.Getenv("fsKey")

	//j, _ := json.Marshal(&event.Detail)
	//fmt.Println(string(j))
	_ = json.Unmarshal(event.Detail, &detail)
	CaseID := fmt.Sprintf("%v", detail.CaseID)
	EventName := fmt.Sprintf("%v", detail.EventName)
	DisplayID := fmt.Sprintf("%v", detail.DisplayID)
	CommunicationID := fmt.Sprintf("%v", detail.CommunicationID)
	Origin := fmt.Sprintf("%v", detail.Origin)
	Subject := ""
	SubmittedBy := ""
	Status := ""
	CaseBody := ""
	var CCEmail []string

	fmt.Printf("CaseID: %#v\n", CaseID)
	fmt.Printf("EventName: %#v\n", EventName)
	fmt.Printf("DisplayID: %#v\n", DisplayID)
	fmt.Printf("CommunicationID: %#v\n", CommunicationID)
	fmt.Printf("Origin: %#v\n", Origin)
	fmt.Printf("CaseBody: %#v\n", CaseBody)

	// 先取api資料, 整理數據
	data := getCaseByCaseID(CaseID)
	// fmt.Printf("data: %#v\n", data)
	for _, value := range *data {
		Subject = fmt.Sprintf("%v", *value.Subject)
		SubmittedBy = fmt.Sprintf("%v", *value.SubmittedBy)
		CCEmail = value.CcEmailAddresses
		Status = fmt.Sprintf("%v", *value.Status)
		CaseBody = fmt.Sprintf("%v", *value.RecentCommunications.Communications[0].Body)
		fmt.Printf("Subject: %#v\n", Subject)
		fmt.Printf("SubmittedBy: %#v\n", SubmittedBy)
		fmt.Printf("CcEmailAddresses: %#v\n", CCEmail)
		fmt.Printf("Status: %#v\n", Status)
		fmt.Printf("Body: %#v\n", CaseBody)
	}

	if EventName == "CreateCase" {
		subject := fmt.Sprintf("AWS Support Case was created [%v] - %v", DisplayID, Subject)
		desc := fmt.Sprintf("You can access the case at the following url: https://console.aws.amazon.com/support/home#/case/?displayId=%v \n %s", DisplayID, CaseBody)
		CCEmail = append(CCEmail, "your@gmail.com")
		email := "your@gmail.com"
		postData := model.FS{
			Description: desc,
			Subject:     subject,
			Email:       email,
			Priority:    4,
			Status:      2,
			CcEmail:     CCEmail,
			Source:      7,
			CustomFields: model.CustomFields{
				CaseId:    CaseID,
				DisplayId: DisplayID,
			},
		}
		createFS(postData)
	}

	if EventName == "AddCommunicationToCase" {
		// todo
	}

	if EventName == "ResolveCase" {
		// todo
	}

	return nil
}

func myHttp(method, url, postStr string, target interface{}) error {
	client := &http.Client{}
	var jsonStr = []byte(postStr)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println(url)
		log.Println(err.Error())
		return nil
	}
	req.SetBasicAuth(fsKey, "X")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("fs error:" + err.Error())
		return nil
	}
	defer resp.Body.Close()

	if target != nil {
		return json.NewDecoder(resp.Body).Decode(target)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	fmt.Println(string([]byte(body)))
	return nil
}

// https://api.freshservice.com/#create_ticket
func createFS(postData model.FS) *model.FSResponse {
	b, _ := json.Marshal(postData)
	url := fmt.Sprintf("https://%s/api/v2/tickets/", endpoint)
	res := new(model.FSResponse)
	_ = myHttp(http.MethodPost, url, string(b), res)
	return res
}

// https://api.freshservice.com/#create_a_reply
func replyFS(fsTickerID string, postData model.ReplyFS) {
	b, _ := json.Marshal(postData)
	url := fmt.Sprintf("https://%s/api/v2/tickets/%s/reply", endpoint, fsTickerID)
	_ = myHttp(http.MethodPost, url, string(b), nil)
}

func getCaseByCaseID(caseId string) *[]types.CaseDetails {
	svc := support.NewFromConfig(cfg)
	caseID := []string{caseId}
	params := &support.DescribeCasesInput{CaseIdList: caseID}
	resp, err := svc.DescribeCases(context.TODO(), params)
	if err != nil {
		log.Fatalf("resp, %v", err)
	}

	cases := &resp.Cases

	return cases
}

func getCaseByDisplayID(DisplayId string) *[]types.CaseDetails {
	svc := support.NewFromConfig(cfg)
	params := &support.DescribeCasesInput{DisplayId: aws.String(DisplayId)}
	resp, err := svc.DescribeCases(context.TODO(), params)
	if err != nil {
		log.Fatalf("resp, %v", err)
	}

	cases := &resp.Cases

	return cases
}

func createCase() {
	svc := support.NewFromConfig(cfg)

	params := &support.CreateCaseInput{
		CommunicationBody: aws.String("CommunicationBody"), // Required
		Subject:           aws.String("Subject"),           // Required
		//AttachmentSetId:   aws.String("AttachmentSetId"),
		CategoryCode: aws.String("charge-inquiry"),
		CcEmailAddresses: []string{
			"andy@test.com",
		},
		IssueType:    aws.String("customer-service"),
		Language:     aws.String("en"),
		ServiceCode:  aws.String("billing"),
		SeverityCode: aws.String("low"),
	}
	resp, err := svc.CreateCase(context.TODO(), params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%+v\n", *resp)
}

func getAttachment(AttachmentId, AttachmentFileName string) {
	svc := support.NewFromConfig(cfg)
	// caseID := []string{"case-185271018684-muen-2022-957fef9ef36a2a40", "case-185271018684-muen-2022-b808f7058d891ad"}
	params := &support.DescribeAttachmentInput{AttachmentId: aws.String(AttachmentId)}
	resp, err := svc.DescribeAttachment(context.TODO(), params)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := &resp.Attachment.Data

	f, err := os.Create(AttachmentFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.Write(*data)

	if err2 != nil {
		log.Fatal(err2)
	}
}
