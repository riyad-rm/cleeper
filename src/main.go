package main


import "./awsHandler"

import (
        "context"
        "github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, event awsHandler.LambdaTrigger) (string, error) {
    // We need to return errors
    awsHandler.Action(event)
    return "ok", nil
}

func main(){
	lambda.Start(HandleRequest)
}