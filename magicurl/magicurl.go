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
func Create(originalURL string, client *dynamodb.DynamoDB) (*MagicURL, error) {
	//TODO: perform input validation the url
	err := validateURL(originalURL)
	if err != nil {
		message := fmt.Sprintf("Invalid URL format: %s", originalURL)
		return nil, &CreateMagicURLSlugError{message}
	}

	id, err := incrementBase10Counter(client)
	if err != nil {
		message := fmt.Sprintf("Could not update base counter for slug")
		return nil, &CreateMagicURLSlugError{message}
	}

	idString := strconv.Itoa(id)
	magicURLItem, err := createMagicURLItem(originalURL, idString, client)
	if err != nil {
		message := fmt.Sprintf("Failed to create MagicURL item for slug")
		return nil, &CreateMagicURLSlugError{message}
	}

	return magicURLItem, nil
}

//Get retrieves the Slug from the database, if not found returns exception
func Get(urlSlug string, client *dynamodb.DynamoDB) (*MagicURL, error) {
	result, err := client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(magicURLTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Slug": {
				S: aws.String(urlSlug),
			},
		},
	})
	if err != nil {
		message := fmt.Sprintf("Could not locate MagicURL for slug: %s", urlSlug)
		return nil, &GetMagicURLSlugError{message}
	}

	magicURLItem := MagicURL{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &magicURLItem)
	if err != nil {
		message := fmt.Sprintf("Failed to create MagicURL result for slug: %s", urlSlug)
		return nil, &GetMagicURLSlugError{message}
	}

	return &magicURLItem, nil
}

//Delete retrieves the slug stored in the data store and removes it.
func Delete(urlSlug string, client *dynamodb.DynamoDB) (*MagicURL, error) {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Slug": {
				S: aws.String(urlSlug),
			},
		},
		ConditionExpression: aws.String("attribute_not_exists(Base10Counter)"),
		TableName:           aws.String(magicURLTable),
		ReturnValues:        aws.String("ALL_OLD"),
	}

	result, err := client.DeleteItem(input)
	if err != nil {
		message := fmt.Sprintf("Could not locate MagicURL for slug: %s", urlSlug)
		return nil, &DeleteMagicURLSlugError{message}
	}

	magicURLItem := MagicURL{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &magicURLItem)
	if err != nil {
		message := fmt.Sprintf("Failed to create MagicURL result for slug: %s", urlSlug)
		return nil, &DeleteMagicURLSlugError{message}
	}

	return &magicURLItem, nil
}

func createMagicURLItem(originalURL string, id string, client *dynamodb.DynamoDB) (*MagicURL, error) {
	slug, err := EncodeToBase62(id)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return Get(slug, client)
}

func incrementBase10Counter(client *dynamodb.DynamoDB) (int, error) {
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(magicURLTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Slug": {
				S: aws.String("A"),
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
