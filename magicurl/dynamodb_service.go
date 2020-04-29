package magicurl

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DynamoDbService is a wrapper around the dynamodb client.
type DynamoDbService struct {
	Client *dynamodb.DynamoDB
}

//New initializes a new instance of the DynamoDbService
func NewDynamoDbService(client *dynamodb.DynamoDB) *DynamoDbService {
	return &DynamoDbService{client}
}

// Get retrieves the MagicURL item from the DynamoDb datastore
func (d *DynamoDbService) Get(slug string) (*MagicURL, error) {
	result, err := d.Client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(magicURLTable),
		Key: map[string]*dynamodb.AttributeValue{
			"Slug": {
				S: aws.String(slug),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	magicURLItem := MagicURL{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &magicURLItem)
	if err != nil {
		return nil, err
	} else if magicURLItem.IsEmpty() {
		return nil, fmt.Errorf("Failed query to find MagicURL with slug %s", slug)
	}

	return &magicURLItem, err
}

// IncrementCounter performs an atomic update of a counter in dynamodb
// Returns the counter value in decimal.
func (d *DynamoDbService) IncrementCounter() (int, error) {
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

	res, err := d.Client.UpdateItem(updateInput)
	if err != nil {
		return -1, err
	}

	var counterValue int
	err = dynamodbattribute.Unmarshal(res.Attributes["Base10Counter"], &counterValue)
	if err != nil {
		return -1, err
	}
	return counterValue, nil
}

// PutMagicURLItem creates MagicURL item containing originalURL and slug in DynamoDB.
func (d *DynamoDbService) PutMagicURLItem(slug, originalURL string) error {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Slug": {
				S: aws.String(slug),
			},
			"OriginalURL": {
				S: aws.String(originalURL),
			},
		},
		TableName:           aws.String(magicURLTable),
		ConditionExpression: aws.String("attribute_not_exists(Slug)"),
	}

	_, err := d.Client.PutItem(input)
	return err
}

// DeleteMagicURLItem removes the MagicURL item containing originalURL and slug in DynamoDB.
func (d *DynamoDbService) DeleteMagicURLItem(slug string) (*MagicURL, error) {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Slug": {
				S: aws.String(slug),
			},
		},
		ConditionExpression: aws.String("attribute_not_exists(Base10Counter)"),
		TableName:           aws.String(magicURLTable),
		ReturnValues:        aws.String("ALL_OLD"),
	}

	result, err := d.Client.DeleteItem(input)
	if err != nil {
		return nil, err
	}

	magicURLItem := MagicURL{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &magicURLItem)
	if err != nil {
		return nil, err
	} else if magicURLItem.IsEmpty() {
		return nil, fmt.Errorf("could not find slug %s in dynamodb table", slug)
	}

	return &magicURLItem, nil
}
