
AWS_REGION = eu-west-2
FUNCTION_NAME = heartbeatLambda
LAMBDA_ROLE = <ROLE>
LAMBDA_RUNTIME = go1.x
LAMBDA_TIMEOUT = 5
LAMBDA_MEMORY_SIZE = 128
TELEGRAM_TOKEN = <TOKEN>
TELEGRAM_CHANNEL = <CHANNEL>

install: build_send build_check_lambda
zip: lambda_zip
deploy: lambda_zip lambda_deploy
create: lambda_create

build_send:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -extldflags "-static"' -o heartbeat

build_check_lambda:
	cd ./lambda && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../heartbeatLambda

lambda_zip:
	zip $(FUNCTION_NAME).zip $(FUNCTION_NAME)

lambda_deploy:
	aws lambda update-function-code \
		--region $(AWS_REGION) \
		--function-name $(FUNCTION_NAME) \
		--environment Variables="{TELEGRAM_TOKEN=$(TELEGRAM_TOKEN),TELEGRAM_CHANNEL=$(TELEGRAM_CHANNEL)}" \
		--zip-file fileb://$(FUNCTION_NAME).zip \
		--publish

lambda_create:
	aws lambda create-function \
		--region $(AWS_REGION) \
		--function-name $(FUNCTION_NAME) \
		--zip-file fileb://$(FUNCTION_NAME).zip \
		--role $(LAMBDA_ROLE) \
		--handler $(FUNCTION_NAME) \
		--runtime $(LAMBDA_RUNTIME) \
		--timeout $(LAMBDA_TIMEOUT) \
		--environment Variables="{TELEGRAM_TOKEN=$(TELEGRAM_TOKEN),TELEGRAM_CHANNEL=$(TELEGRAM_CHANNEL)}" \
		--memory-size $(LAMBDA_MEMORY_SIZE)