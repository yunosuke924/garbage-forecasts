package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

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


type Forecast struct {
	// ID                    int `json:"id"`
	ForecastDate          string `json:"forecast_date"`
	GarbageType           int `json:"garbage_type"`
	GarbageForecastIndex  float64 `json:"garbage_forecast_index"`
	WeatherForecast       int `json:"weather_forecast"`
	// CreatedAt             string `json:"created_at"`
	// UpdatedAt             string `json:"updated_at"`
}

// Forecastsレコードを作成する
func createForecast(db *sql.DB, f Forecast) {
	fmt.Print("=======createForecast is started===================")
	fmt.Println(f)
	_, err := db.Exec(
		"insert into forecasts (forecast_date, garbage_type, garbage_forecast_index, weather_forecast) values (?, ?, ?, ?)",
		f.ForecastDate, f.GarbageType, f.GarbageForecastIndex, f.WeatherForecast,
	)
	if err != nil {
		fmt.Println("レコード作成エラー")
	}
	fmt.Print("=======createForecast is started===================")
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("=======handler is started===================")
	var payload Forecast
	fmt.Println(request.Body)
	err := json.Unmarshal([]byte(request.Body), &payload)
	if err != nil {
		// JSONデコードが失敗した場合のエラーハンドリング
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error decoding JSON",
		}, nil
	}
	db, err := sql.Open("mysql", os.Getenv("DB_CONNECTION")) // データベースに接続
	if err != nil {
		fmt.Println("DB接続エラー")
	}
	
	// 関数がリターンする直前に呼び出される
	defer func() {
		fmt.Println("=======handler is finished===================")
		db.Close()
	}()

	// Forecastsテーブルを作成
	// 例
	createForecast(db, payload)

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

	resJson := Forecast{
		ForecastDate:          payload.ForecastDate,
		GarbageType:           payload.GarbageType,
		GarbageForecastIndex:  payload.GarbageForecastIndex,
		WeatherForecast:       payload.WeatherForecast,
	}

	fmt.Println("==========================")
	fmt.Println(resJson)
	fmt.Println("==========================")
	res, err := json.Marshal(resJson)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(res),
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
