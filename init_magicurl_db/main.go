package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

var svc = dynamodb.New(sess, &aws.Config{
	Region: aws.String("us-east-1"),
})

var magicURLTable = "magicUrl"

func main() {
	if !dynamoDbTableExists(magicURLTable) {
		fmt.Println("Table magicURL has not been created")
		os.Exit(1)
	}

	if !base10CounterInitialized() {
		fmt.Printf("Counter not initialized. Initializing...")
		err := initializeBase10Counter()
		if err != nil {
			fmt.Println("Could not initialize the counter")
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("done.")
	}
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

func base10CounterInitialized() bool {
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(magicURLTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Slug": {
				S: aws.String("A"),
			},
		},
	})
	if err != nil {
		return false
	}

	return len(result.Item) > 0
}

func dynamoDbTableExists(tableName string) bool {
	input := &dynamodb.ListTablesInput{
		Limit: aws.Int64(1),
	}

	res, err := svc.ListTables(input)
	if err != nil {
		log.Fatal(err)
	}
	if len(res.TableNames) == 0 || *res.TableNames[0] != tableName {
		return false
	}

	return true
}
