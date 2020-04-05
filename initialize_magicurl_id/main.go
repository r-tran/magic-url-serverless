package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var svc = dynamodb.New(sess)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) error {
	addItem("magicUrl")
}

func main() {
	lambda.Start(Handler)
}

func addItem(tableName string) error {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodbAttributmap[string]*dynamodb.AttributeValue{
			"Slug"	: {
				S: aws.String("0")
			}
			"Value" : {
				N: aws.Int64(0)
			}
		},
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		return err
	}
	return nil
}
