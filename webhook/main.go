package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	_ "github.com/joho/godotenv/autoload"
	"io/ioutil"
	"log"
	"main/model"
	"net/http"
	"os"
)

type Response struct {
	Message string `json:"message"`
	At      string `json:"at"`
}

var cfg aws.Config
var err error
var endpoint string
var region string
var zones string
var fsKey string
var CCEmail []string

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	region = os.Getenv("region")
	zones = os.Getenv("zones")
	endpoint = os.Getenv("fs_domain")
	fsKey = os.Getenv("fsKey")

	fmt.Println("env region:", region)
	fmt.Println("env zones:", zones)
	fmt.Println("env fs_endpoint:", endpoint)
	fmt.Println("env fsKey:", fsKey)

	//fmt.Printf("req: %#v\n", req.Body)

	msgId := req.QueryStringParameters["msg_id"]
	taskId := req.QueryStringParameters["task_id"]
	taskType := req.QueryStringParameters["task_type"]
	faultTime := req.QueryStringParameters["fault_time"]
	taskStatus := req.QueryStringParameters["task_status"]
	taskSummary := req.QueryStringParameters["task_summary"]
	taskName := req.QueryStringParameters["taskName"]
	serverId := req.QueryStringParameters["server_id"]
	content := req.QueryStringParameters["content"]
	token := req.QueryStringParameters["token"]
	urlToken := ""

	fmt.Println(msgId, taskId, taskType, faultTime, taskStatus, taskSummary, taskName, serverId, content, token, urlToken)

	str := fmt.Sprintf("%s%s%s%s", msgId, taskId, faultTime, urlToken)
	data := []byte(str)
	has := md5.Sum(data)
	md5Token := fmt.Sprintf("%x", has) //将[]byte转成16进制

	fmt.Println(md5Token, token)

	// 搞清楚 token的判斷後, 再填入
	if md5Token == token {
		subject := fmt.Sprintf("監控 - 系統監控 - [%s]", msgId)
		desc := fmt.Sprintf(content)
		CCEmail = append(CCEmail, "develop@test.com")
		email := "support@test.com"
		postData := model.FS{
			Description: desc,
			Subject:     subject,
			Email:       email,
			Priority:    4,
			Status:      2,
			CcEmail:     CCEmail,
			Source:      7,
		}
		fs := createFS(postData)

		body, _ := json.Marshal(fs)
		res := events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers:    map[string]string{"Content-Type": "text/json; charset=utf-8"},
			Body:       string(body),
		}
		return res, nil
	}

	res := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "text/json; charset=utf-8"},
		Body:       string("外來的不要亂搓！！！！"),
	}

	return res, nil
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
