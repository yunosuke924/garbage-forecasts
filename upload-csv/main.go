package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Data struct {
	File     string `json:"file"`
	FileName string `json:"file_name"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// S3の接続を初期化
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String("http://host.docker.internal:9000"),
		Region:           aws.String("ap-northeast-1"),
		Credentials:      credentials.NewStaticCredentials("root", "password1", ""),
		S3ForcePathStyle: aws.Bool(true),
	})

	if err != nil {
		panic(err)
	}

	var data Data
	// リクエストボディを構造体に変換してdataに格納
	err = json.Unmarshal([]byte(request.Body), &data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(data)

	// base64でエンコードされたデータをデコード
	dec, err := base64.StdEncoding.DecodeString(data.File)
	log.Println(string(dec))

	// S3のアップローダーを作成
	uploader := s3manager.NewUploader(sess)
	r := strings.NewReader(string(dec))

	// データをS3へアップロード
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("tokyo-garbage-data"),
		Body:   r,
		Key:    aws.String("/" + data.FileName),
	})
	if err != nil {
		panic(err)
	}
	return events.APIGatewayProxyResponse{
		Body:       "OK",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
