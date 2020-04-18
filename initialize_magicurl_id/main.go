package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var magicURLTable = "magicUrl"
var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var svc = dynamodb.New(sess)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) error {
	return initializeBase10Counter()
}

func initializeBase10Counter() error {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Slug": {
				S: aws.String("A"),
			},
			"Base10Counter": {
				N: aws.String("0"),
			},
		},
		TableName: aws.String(magicURLTable),
	}

	_, err := svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
