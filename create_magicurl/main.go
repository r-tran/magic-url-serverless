package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/personal_projects/magic-url-serverless/magicurl"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// MagicURLRequest contains the original URL
type MagicURLRequest struct {
	OriginalURL string `json:"url`
}

// Validate checks that the MagicURL request URL is a valid URL format
func (*MagicURLRequest) Validate() error {
	return nil
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var magicURLRequest MagicURLRequest
	json.Unmarshal([]byte(request.Body), &magicURLRequest)
	originalURL := magicURLRequest.OriginalURL

	err := magicURLRequest.Validate()
	if err != nil {
		return Response{Body: "Error", StatusCode: 400}, err
	}

	slug, err := magicurl.Create(originalURL)
	if err != nil {
		return Response{Body: "Error", StatusCode: 400}, err
	}

	var buf bytes.Buffer
	body, err := json.Marshal(map[string]interface{}{
		"slug": slug,
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "create-magicurl-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
