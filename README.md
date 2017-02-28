# kinesis-follow

A command line tool to print data from an AWS Kinesis stream.

## Usage

```
▶ kinesis-follow --help
Usage of kinesis-follow:
  -region string
    	AWS region the stream exists in (default "eu-west-1")
  -sleep int
    	How long to sleep between polling for new messages in ms (default 1000)
  -stream string
    	Name of the string you want to follow
```

## Installation

Assuming you have golang installed and setup correctly:
```
go install github.com/uswitch/kinesis-follow
```

## TODO
- Get a shard iterator with options other than `LATEST`
- Don't quit on keypress