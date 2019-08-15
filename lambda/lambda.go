package main

import (
	"context"
	"heartbeat/pkg/aws"

	"github.com/aws/aws-lambda-go/lambda"
)

// HandleRequest lambda handler request
func HandleRequest(ctx context.Context, name interface{}) (string, error) {
	aws.CheckAliveSQSMessage()
	return "done", nil
}

func main() {
	lambda.Start(HandleRequest)
}
