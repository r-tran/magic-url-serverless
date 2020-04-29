# MagicURL Serverless
MagicURL is simple URL shortener API built on top of AWS with [Lambda](https://docs.aws.amazon.com/lambda/), [ApiGateway](https://docs.aws.amazon.com/apigateway/), and [DynamoDb](https://docs.aws.amazon.com/dynamodb/).

The project uses the [Serverless](https://www.serverless.com/framework/docs/) Framework. This framework improves the serverless app development experience and get rid of all the boilerplate. The functions deployment and provisioning steps are all captured in `serveless.yml`.

## Requirements

1. You will need to install `serverless` on your machine. Find instructions [here](https://www.serverless.com/framework/docs/getting-started/). 
2. If you do not have an AWS account configure on the machine, you'll need to create one and get your credentials set up. The Serverless docs provide a walkthrough [here](https://www.serverless.com/framework/docs/providers/aws/guide/credentials/)
3. You will need to have the `Go` programming language and `make` build tool installed. I recommend going through the [Hello World Go](https://www.serverless.com/framework/docs/providers/aws/examples/hello-world/go/) quickstart.

## Run the Code

After the requirements above have been satisfied, you can clone the code and run `make deploy`. 
This will deploy the serverless functions `create_magic_url`, `delete_magic_url`, `get_magicurl` in AWS.

## Resources
Blog post/video with a walkthrough coming soon!