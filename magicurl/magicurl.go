package magicurl

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var magicURLTable = "magicUrl"

//Get returns the Magic URL given the shortened URL slug.
//If the Magic URL is not found, throws GetMagicURLSlugError
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
		return nil, &GetMagicURLSlugError{message, err}
	}

	magicURLItem := MagicURL{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &magicURLItem)
	if err != nil {
		message := fmt.Sprintf("Failed to create MagicURL result for slug: %s", urlSlug)
		return nil, &GetMagicURLSlugError{message, err}
	}

	return &magicURLItem, nil
}

// Create returns a MagicURL containing the shortened URL slug for the originalURL.
func Create(originalURL string, client *dynamodb.DynamoDB) (*MagicURL, error) {
	sanitizedURL, err := validateURL(originalURL)
	if err != nil {
		message := fmt.Sprintf("Invalid URL format: %s", originalURL)
		return nil, &CreateMagicURLSlugError{message, err}
	}

	id, err := incrementBase10Counter(client)
	if err != nil {
		message := fmt.Sprintf("Could not update base counter for slug")
		return nil, &CreateMagicURLSlugError{message, err}
	}

	idString := strconv.Itoa(id)
	magicURLSlug, err := createMagicURLItem(sanitizedURL, idString, client)
	if err != nil {
		message := fmt.Sprintf("Failed to create MagicURL item for slug")
		return nil, &CreateMagicURLSlugError{message, err}
	}

	return Get(magicURLSlug, client)
}

//Delete removes MagicURl slug from datastore, returns the deleted slug.
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
		return nil, &DeleteMagicURLSlugError{message, err}
	}

	magicURLItem := MagicURL{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &magicURLItem)
	if err != nil {
		message := fmt.Sprintf("Failed to create MagicURL result for slug: %s", urlSlug)
		return nil, &DeleteMagicURLSlugError{message, err}
	}

	return &magicURLItem, nil
}

func validateURL(urlTarget string) (string, error) {
	parsedURL, err := url.ParseRequestURI(urlTarget)
	if err != nil {
		return "", err
	}
	return parsedURL.String(), err
}

func createMagicURLItem(originalURL string, id string, client *dynamodb.DynamoDB) (string, error) {
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

	return slug, err
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
