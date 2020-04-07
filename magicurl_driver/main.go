package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/personal_projects/magic-url-serverless/magicurl"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var svc = dynamodb.New(sess, &aws.Config{
	Endpoint: aws.String("http://localhost:8000"),
})

var magicURLTable = "magicUrl"

func main() {
	if !dynamoDbTableExists(magicURLTable) {
		fmt.Println("Creating magicURL table...")
		err := createMagicURLTable()
		if err != nil {
			fmt.Println("Error while creating magicURL Table")
			os.Exit(1)
		}
	} else {
		fmt.Printf("magicURL table already exists\n")
	}

	if !base10CounterInitialized() {
		fmt.Println("Counter not initialized. Initializing...")
		err := initializeBase10Counter()
		if err != nil {
			fmt.Println("Could not initialize the counter")
			os.Exit(1)
		}
	}

	//	err := deleteBase10Counter()
	//	if err != nil {
	//		fmt.Println("Could not delete the counter")
	//		os.Exit(1)
	//	}
	//	fmt.Println("Deleted counter")

	/* 	for i := 0; i < 10; i++ {
		count, err := magicurl.IncrementBase10Counter(svc)
		if err != nil {
			fmt.Println("Could not increment the counter")
			os.Exit(1)
		}
		fmt.Printf("Value of base 10 counter: %v\n", count)
	} */

	// Test create slug
	slug, err := magicurl.Create("https://pythonsandpenguins.dev", svc)
	if err != nil {
		fmt.Println("Creating slug had an error")
		os.Exit(1)
	}

	fmt.Printf("Result is %v\n", slug)
}

//TODO: Extract into provisioning step
func initializeBase10Counter() error {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Slug": {
				S: aws.String("0"),
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
				S: aws.String("0"),
			},
		},
	})
	if err != nil {
		return false
	}

	return len(result.Item) > 0
}

func getBase10Count() (int, error) {
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(magicURLTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Slug": {
				S: aws.String("0"),
			},
		},
	})
	if err != nil {
		return -1, err
	}

	var base10Count int
	dynamodbattribute.Unmarshal(result.Item["IdBase10"], &base10Count)
	return base10Count, nil
}

func deleteBase10Counter() error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Slug": {
				S: aws.String("0"),
			},
		},
		TableName: aws.String(magicURLTable),
	}
	_, err := svc.DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
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

func createMagicURLTable() error {
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
		TableName: aws.String(magicURLTable),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		fmt.Println("Error creating the MagicUrl Table")
		return err
	}

	return nil
}
