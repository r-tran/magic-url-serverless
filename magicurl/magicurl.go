package magicurl

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Create shortened URL to the database.
func Create(originalURL string, client *dynamodb.DynamoDB) (string, error) {
	//TODO: perform input validation the url
	err := validateURL(originalURL)
	if err != nil {
		return "", err
	}

	//get id from db, atomic increment value in the db
	//id, err := updateMagicURLId(client)
	_, err = updateMagicURLId(client)
	if err != nil {
		return "", err
	}

	//create slug as base64-encoded id

	//create MagicUrlItem, insert in the db

	//on successful add,  return created slug
	return "raytran_slug", nil
}

func updateMagicURLId(client *dynamodb.DynamoDB) (int, error) {
	_, err := client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("magicUrl"),
		Key: map[string]*dynamodb.AttributeValue{
			"Slug": {
				N: aws.String("0"),
			},
		},
	})
	if err != nil {
		return -1, err
	}

	return 0, nil
}

func validateURL(url string) error {
	return nil
}
