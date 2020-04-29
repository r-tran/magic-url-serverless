# MagicURL Serverless
MagicURL is simple URL shortener API built on top of AWS with [Lambda](https://docs.aws.amazon.com/lambda/), [ApiGateway](https://docs.aws.amazon.com/apigateway/), and [DynamoDB](https://docs.aws.amazon.com/dynamodb/).

The project uses the [Serverless](https://www.serverless.com/framework/docs/) Framework. This framework improves the serverless development experience by reducing the amount of boilerplate needed to deploy an application. The steps for defining deployments of functions and other provisioning steps are all captured in `serverless.yml`.

## Requirements

1. You will need to install `serverless` on your machine. Please refer to the installation section in the [Getting Start Guide](https://www.serverless.com/framework/docs/getting-started/). 
2. If you do not have an AWS account configured for your machine, you'll need to create one and also get your credentials set up. The Serverless docs provide a good walkthrough [here](https://www.serverless.com/framework/docs/providers/aws/guide/credentials/).
3. You will need to have the `Go` programming language and `make` build tool installed. I recommend going through the [Hello World Go](https://www.serverless.com/framework/docs/providers/aws/examples/hello-world/go/) quickstart.

## Run the Code

After the above requirements above have been satisfied, you can clone the clone the repository, set it as the working directory, and run: `make deploy`. 

This will deploy the serverless functions `create_magic_url`, `delete_magic_url`, `get_magicurl` in AWS.

## Resources
Blog post/video with a walkthrough coming soon!