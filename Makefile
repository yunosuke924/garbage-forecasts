.PHONY: build

build:
	sam build
deploy:
	sam deploy --profile go-academy
local:
	sam local start-api --docker-network go-academy
