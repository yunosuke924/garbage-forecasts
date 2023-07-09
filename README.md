# ca-sam-app-08-garbage-forecast

This is a sample template for ca-sam-app-08-garbage-forecast - Below is a brief explanation of what we have generated for you:

```bash
.
├── Makefile                    <-- Make to automate build
├── README.md                   <-- This instructions file
├── hello-world                 <-- Source code for a lambda function
│   ├── main.go                 <-- Lambda function code
│   └── main_test.go            <-- Unit tests
└── template.yaml
```

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Setup process

### Installing dependencies & building the target 

In this example we use the built-in `sam build` to automatically download all the dependencies and package our build target.   
Read more about [SAM Build here](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-build.html) 

The `sam build` command is wrapped inside of the `Makefile`. To execute this simply run
 
```shell
make
```

### Local development

**Invoking function locally through local API Gateway**

```bash
sam local start-api
```

If the previous command ran successfully you should now be able to hit the following local endpoint to invoke your function `http://localhost:3000/hello`

**SAM CLI** is used to emulate both Lambda and API Gateway locally and uses our `template.yaml` to understand how to bootstrap this environment (runtime, where the source code is, etc.) - The following excerpt is what the CLI will read in order to initialize an API and its routes:

```yaml
...
Events:
    HelloWorld:
        Type: Api # More info about API Event Source: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#api
        Properties:
            Path: /hello
            Method: get
```

## Packaging and deployment

AWS Lambda Golang runtime requires a flat folder with the executable generated on build step. SAM will use `CodeUri` property to know where to look up for the application:

```yaml
...
    FirstFunction:
        Type: AWS::Serverless::Function
        Properties:
            CodeUri: hello_world/
            ...
```

To deploy your application for the first time, run the following in your shell:

```bash
sam deploy --guided
```

The command will package and deploy your application to AWS, with a series of prompts:

* **Stack Name**: The name of the stack to deploy to CloudFormation. This should be unique to your account and region, and a good starting point would be something matching your project name.
* **AWS Region**: The AWS region you want to deploy your app to.
* **Confirm changes before deploy**: If set to yes, any change sets will be shown to you before execution for manual review. If set to no, the AWS SAM CLI will automatically deploy application changes.
* **Allow SAM CLI IAM role creation**: Many AWS SAM templates, including this example, create AWS IAM roles required for the AWS Lambda function(s) included to access AWS services. By default, these are scoped down to minimum required permissions. To deploy an AWS CloudFormation stack which creates or modifies IAM roles, the `CAPABILITY_IAM` value for `capabilities` must be provided. If permission isn't provided through this prompt, to deploy this example you must explicitly pass `--capabilities CAPABILITY_IAM` to the `sam deploy` command.
* **Save arguments to samconfig.toml**: If set to yes, your choices will be saved to a configuration file inside the project, so that in the future you can just re-run `sam deploy` without parameters to deploy changes to your application.

You can find your API Gateway Endpoint URL in the output values displayed after deployment.

### Testing

We use `testing` package that is built-in in Golang and you can simply run the following command to run our tests:

```shell
go test -v ./hello-world/
```
# Appendix

### Golang installation

Please ensure Go 1.x (where 'x' is the latest version) is installed as per the instructions on the official golang website: https://golang.org/doc/install

A quickstart way would be to use Homebrew, chocolatey or your linux package manager.

#### Homebrew (Mac)

Issue the following command from the terminal:

```shell
brew install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
brew update
brew upgrade golang
```

#### Chocolatey (Windows)

Issue the following command from the powershell:

```shell
choco install golang
```

If it's already installed, run the following command to ensure it's the latest version:

```shell
choco upgrade golang
```

## Bringing to the next level

Here are a few ideas that you can use to get more acquainted as to how this overall process works:

* Create an additional API resource (e.g. /hello/{proxy+}) and return the name requested through this new path
* Update unit test to capture that
* Package & Deploy

Next, you can use the following resources to know more about beyond hello world samples and how others structure their Serverless applications:

* [AWS Serverless Application Repository](https://aws.amazon.com/serverless/serverlessrepo/)

## 概要

- 地域のゴミの収集日情報と気象庁が公開している天気予報情報をもとに、ゴミ出しをするべきか否かを判断することができるサービス

## **解決する問題**

- 大雪、台風、豪雨、大風など、ゴミ出しの日と悪天候が重なることで発生する以下の問題の解決の一助を目指す。
    - 家庭ゴミなどを適切なタイミングで廃棄しないことで、不衛生な状況を引き起こす可能性がある。
    - ゴミが溜まることで居住スペースを圧迫してしまう可能性がある。
    - 強風によってゴミが飛ばされたりして環境汚染やゴミ回収業者など関係者の方に迷惑をかけてしまう。

## 機能要件

優先度：高

- **ゴミの種類**から、**ゴミ出し日指数**を表示
    - 指定した種類のゴミの二週間分の**ゴミ出し日指数**を表示。
- **ゴミ出し日**から、**ゴミ出し日和指数**を表示。
    - 指定した日の全種類の**ゴミ出し日指数**を表示。
- **ゴミ種別**と**ゴミ出し日**から、**ゴミ出し日指数**を表示。
    - 指定した日と種類の**ゴミ出し日指数**を表示。

優先度：中

- リマインダー機能：ゴミの種類を登録しておいて、次にそのゴミを出せる日が、強風や大雨などで適切ではない場合にslack チャンネルにアラートを送信。

## 開発要件

- SAM アプリケーションの作成。
- マイグレーションファイルの作成。
    
    ```yaml
    id: ID
    date: 予報日
    garbage_type: ゴミの分類
    garbage_forecast_index: ゴミ出し日和指数
    weather_forecast: 天気予報
    created_at: レコード作成日
    updated_at: レコード更新日
    ```
    
- dockerコンテナを起動してローカルで動作確認。
- 一旦仮デプロイ、URL確認
- 本番DBにデータを注入、DB確認。
- 天気予報APIから天気予報情報を取得、レスポンス確認。
- ゴミの種類ごとの収集日を取得、レスポンスの確認。
- 上記二つをもとに次の日のforecastsレコードを作成、更新する関数を作成。
- 上記タスクの定期実行バッチ（Scheduler）の作成。
- **指定した種類**のゴミ出し日指数を表示するエンドポイントの作成
    - GET `/api/v1/garbage_forecasts/search?garbage_type=burnable&days=14`
    
    ```go
    {
        "garbage_type": "burnable",
        "forecast": [
            {
                "date": "2023-07-10",
                "garbage_forecast_index": 0.6,
                "weather_forecast": "Sunny"
            },
            {
                "date": "2023-07-17",
                "garbage_forecast_index": 0.9,
                "weather_forecast": "Cloudy"
            },
            // ... (他の日付の情報)
        ]
    }
    ```
    
- **指定した日**の全種類のゴミ出し日指数を表示するエンドポイントの作成
    - GET `/api/v1/trash_forecasts/search?datetime=2023-7-09`
    
    ```go
    {
    	"date": "2023-07-10",
      "forecast": [
          {
              // "date": "2023-07-10",
    					"trash_type": "burnable",
              "garbage_forecast_index": 0.6,
              "weather_forecast": "Sunny"
          },
          {
              // "date": "2023-07-17",
    					"trash_type": "non_burnable",
              "garbage_forecast_index": 0.9,
              "weather_forecast": "Cloudy"
          },
      ]
    }
    ```
    
- **ゴミ種別**と**ゴミ出し日**から、**ゴミ出し日指数**を表示するエンドポイントの作成
    - GET `/api/v1/garabage_forecasts/search?garbage_type=burnable&datetime=2023-7-09`
    
    ```go
    {
        "garbage_type": "burnable",
        "date": "2023-07-10",
        "garbage_forecast_index": 0.7,
        "weather_forecast": "Sunny"
    }
    ```
