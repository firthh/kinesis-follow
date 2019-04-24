# kinesis-follow

A command line tool to print data from an AWS Kinesis stream.

## Usage

```
Usage of kinesis-follow:
  -endpoint string
    	AWS endpoint (useful with localstack)
  -no-verify-ssl
    	No SSL certificate validation
  -region string
    	AWS region the stream exists in (default "eu-west-1")
  -sleep int
    	How long to sleep between polling for new messages in ms (default 1000)
  -stream string
    	Name of the string you want to follow
  -unzip-data
    	Do you want to unzip the data bytes read from Kinesis, for example if the bytes were zipped before publishing to kinesis
```

## Installation

Assuming you have golang installed and setup correctly:
```
go get github.com/uswitch/kinesis-follow
go install github.com/uswitch/kinesis-follow
```

## With LocalStack
The application supports setting the endpoint for example:
```sh
kinesis-follow --stream "output-kinesis-stream" --endpoint "localhost:4568" --no-verify-ssl --unzip-data
```
This is used when localstack is run in a docker container with SSL enabled

## Permissions

It is assumed that you will have credentials for AWS setup as you would to ordinarily use the AWS SDK - https://github.com/aws/aws-sdk-go#configuring-credentials

## TODO
- Get a shard iterator with options other than `LATEST`
- Don't quit on keypress
