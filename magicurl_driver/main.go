package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

// var svc = dynamodb.New(sess, &aws.Config{
// 	Endpoint: aws.String("http://localhost:8000"),
// })
var svc = dynamodb.New(sess, &aws.Config{
	Region: aws.String("us-east-1"),
})

var magicURLTable = "magicUrl"

func main() {
	if !dynamoDbTableExists(magicURLTable) {
		fmt.Println("Creating magicURL table...")
		err := createMagicURLTable()
		if err != nil {
			fmt.Println("Error while creating magicURL Table")
			fmt.Println(err)
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
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// // Test create slug
	// originalURL := "https://pythonsandpenguins.dev"
	// magicURL, err := magicurl.Create(originalURL, svc)
	// if err != nil {
	// 	fmt.Println("Creating slug had an error")
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("Created Magic URL: %v\n", magicURL)

	// // //Retrieve slug
	// fmt.Println("Retrieving the slug...")
	// magicURLItem, err := magicurl.Get(magicURL.Slug, svc)
	// if err != nil {
	// 	fmt.Println("GET slug had an error")
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("Result: %v\n", magicURLItem)

	// //Delete slug
	// fmt.Println("Deleting the slug...")
	// result, err := magicurl.Delete(magicURLItem.Slug, svc)
	// if err != nil {
	// 	fmt.Println("Delete slug had an error")
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("Result: %v\n", result)
}

//TODO: Extract into provisioning step
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
	if len(res.TableNames) == 0 || *res.TableNames[0] != tableName {
		return false
	}

	return true
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
