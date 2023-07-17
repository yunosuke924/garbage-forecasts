package main

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func downloadCSV(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	csvData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return csvData, nil
}

type Record struct {
	Column1 string
	Column2 string
	// 別のカラムがある場合は、必要に応じてフィールドを追加する
}

func parseCSV(csvData []byte) ([]Record, error) {
	reader := csv.NewReader(strings.NewReader(string(csvData)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var result []Record
	for _, r := range records {
		record := Record{
			Column1: r[0], // カラムのインデックスを指定
			Column2: r[1], // 別のカラムがある場合は、必要に応じてフィールドを追加する
		}
		result = append(result, record)
	}

	return result, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db, err := sql.Open("mysql", "root:password@(db:3306)/training?parseTime=true") // データベースに接続
	if err != nil {
		fmt.Println("DB接続エラー")
	}
	defer db.Close()

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// IPアドレスを取得
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	url := "https://www.opendata.metro.tokyo.lg.jp/taitou/chiikibetu_youbi_ichiran.csv"

	csvData, err := downloadCSV(url)
	if err != nil {
		log.Fatal("CSVのダウンロードエラー:", err)
	}

	records, err := parseCSV(csvData)
	if err != nil {
		log.Fatal("CSVのパースエラー:", err)
	}

	// 構造体のスライスを使って操作を行う
	for _, record := range records {
		fmt.Println(record.Column1, record.Column2)
	}

	return events.APIGatewayProxyResponse{
		Body:       "データをセットしました",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
