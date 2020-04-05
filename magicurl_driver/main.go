package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/personal_projects/magic-url-serverless/magicurl"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var svc = dynamodb.New(sess, &aws.Config{
	Endpoint: aws.String("http://localhost:8000"),
})

func main() {
	//Check if table exists, create if not exists
	magicURLTable := "magicUrl"

	if !dynamoDbTableExists(magicURLTable) {
		fmt.Printf("Creating table %s\n", magicURLTable)
		createDynamoDbTable(magicURLTable)
	} else {
		fmt.Printf("%s table already exists\n", magicURLTable)
	}

	// Test create slug
	slug, err := magicurl.Create("original_url", svc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result is %v\n", slug)
}

func dynamoDbTableExists(tableName string) bool {
	input := &dynamodb.ListTablesInput{
		Limit: aws.Int64(1),
	}

	res, err := svc.ListTables(input)
	if err != nil {
		log.Fatal(err)
	}

	tableFound := *res.TableNames[0] == tableName
	return tableFound
}

func createDynamoDbTable(tableName string) {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Slug"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Slug"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		fmt.Println("Got error calling CreateTable:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
