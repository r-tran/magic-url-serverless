package magicurl

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var magicURLTable = "magicUrl"

// Create shortened URL to the database.
func Create(originalURL string, client *dynamodb.DynamoDB) (string, error) {
	//TODO: perform input validation the url
	err := validateURL(originalURL)
	if err != nil {
		return "", err
	}

	id, err := IncrementBase10Counter(client)
	if err != nil {
		return "", err
	}

	idString := strconv.Itoa(id)
	slug, err := CreateMagicURLItem(originalURL, idString, client)
	if err != nil {
		return "", err
	}

	return slug, nil
}

//CreateMagicURLItem creates an entry in DynamoDb for the MagicUrl
func CreateMagicURLItem(originalURL string, id string, client *dynamodb.DynamoDB) (string, error) {
	slug, err := EncodeToBase62(id)
	if err != nil {
		return "", err
	}

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Slug": {
				S: aws.String(slug),
			},
			"OriginalUrl": {
				S: aws.String(originalURL),
			},
		},
		TableName:           aws.String(magicURLTable),
		ConditionExpression: aws.String("attribute_not_exists(Slug)"),
	}

	_, err = client.PutItem(input)
	if err != nil {
		return "", err
	}

	return slug, nil
}

//IncrementBase10Counter is used for hashing the URL slug
func IncrementBase10Counter(client *dynamodb.DynamoDB) (int, error) {
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(magicURLTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Slug": {
				S: aws.String("0"),
			},
		},
		UpdateExpression: aws.String("SET Base10Counter = Base10Counter + :incr"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":incr": {
				N: aws.String("1"),
			},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}

	res, err := client.UpdateItem(updateInput)
	if err != nil {
		return -1, err
	}

	var counterValue int
	dynamodbattribute.Unmarshal(res.Attributes["Base10Counter"], &counterValue)
	return counterValue, nil
}

func updateMagicURLId(client *dynamodb.DynamoDB) (int, error) {
	res, err := client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("magicUrl"),
		Key: map[string]*dynamodb.AttributeValue{
			"Slug": {
				S: aws.String("0"),
			},
		},
	})

	//Check if the key is existing
	if err != nil {
		return -1, err
	}
	if len(res.Item) == 0 {
		fmt.Println("Item not found")
	}

	return 0, nil
}

func validateURL(url string) error {
	return nil
}
