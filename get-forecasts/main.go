package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

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
	ID                   int     `json:"id"`
	ForecastDate         string  `json:"forecast_date"`
	GarbageType          int     `json:"garbage_type"`
	GarbageForecastIndex float64 `json:"garbage_forecast_index"`
	WeatherForecast      int     `json:"weather_forecast"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}

type ResponseJSON struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Forecastsレコードを作成する
// @param db *sql DB接続情報
// @return []Forecast 予報レコードのスライス
func getForecasts(db *sql.DB) ([]Forecast, error) {
	rows, err := db.Query("select id, forecast_date, garbage_type, garbage_forecast_index, weather_forecast, created_at, updated_at from forecasts")
	if err != nil {
		fmt.Println(rows)
		fmt.Println("レコードの取得に失敗しました")
		return nil, err
	}

	// スライスを定義
	var forecasts []Forecast

	// 取得したデータをループで処理
	for rows.Next() {
		var f Forecast
		err := rows.Scan(&f.ID,
			&f.ForecastDate,
			&f.GarbageType,
			&f.GarbageForecastIndex,
			&f.WeatherForecast,
			&f.CreatedAt,
			&f.UpdatedAt)
		if err != nil {
			fmt.Println("Scanエラー")
			return nil, err
		}
		// 取得したデータをスライスに追加
		forecasts = append(forecasts, f)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("rows.Err")
		return nil, err
	}

	return forecasts, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db, err := sql.Open("mysql", "root:password@(db:3306)/training?parseTime=true") // データベースに接続
	if err != nil {
		fmt.Println("DB接続エラー")
	}

	// 関数がリターンする直前に呼び出される
	defer db.Close()

	// Forecastsテーブルを作成
	// 例
	forecasts, err := getForecasts(db)
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

	res, err := json.Marshal(forecasts)
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
