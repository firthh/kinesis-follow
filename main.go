package main

import (
	// "bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

func getStreamDescription(svc *kinesis.Kinesis, name string) (kinesis.StreamDescription, error) {
	// describe the stream\
	params := &kinesis.DescribeStreamInput{
		StreamName: aws.String(name),
		Limit:      aws.Int64(1),
	}
	resp, err := svc.DescribeStream(params)
	if err != nil {
		return kinesis.StreamDescription{}, err
	}
	return *resp.StreamDescription, nil
}

func getShardIterator(svc *kinesis.Kinesis, streamName string, shardId string) (string, error) {
	params := &kinesis.GetShardIteratorInput{
		ShardId:           aws.String(shardId),
		ShardIteratorType: aws.String("LATEST"),
		StreamName:        aws.String(streamName),
	}
	resp, err := svc.GetShardIterator(params)

	if err != nil {
		return "", err
	}

	return *resp.ShardIterator, nil
}

func getRecords(svc *kinesis.Kinesis, shardIterator string) (kinesis.GetRecordsOutput, error) {

	params := &kinesis.GetRecordsInput{
		ShardIterator: aws.String(shardIterator), // Required
		Limit:         aws.Int64(100),
	}
	resp, err := svc.GetRecords(params)

	if err != nil {
		return kinesis.GetRecordsOutput{}, err
	}
	return *resp, nil
}

func main() {
	streamName := "banking-clicks-production"

	sess := session.New(&aws.Config{Region: aws.String("eu-west-1")})
	svc := kinesis.New(sess)

	stream, err := getStreamDescription(svc, streamName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(*stream.Shards[0].ShardId)

	shardIterator, err2 := getShardIterator(svc, streamName, *stream.Shards[0].ShardId)

	if err2 != nil {
		fmt.Println(err2.Error())
		return
	}
	// Pretty-print the response data.
	fmt.Println(shardIterator)

	for true {
		time.Sleep(1000 * time.Millisecond)
		getRecordsOut, err3 := getRecords(svc, shardIterator)
		if err3 != nil {
			fmt.Println(err3.Error())
			return
		}
		shardIterator = *getRecordsOut.NextShardIterator
		for _, el := range getRecordsOut.Records {
			s := string(el.Data[:])
			fmt.Println(s)
		}

	}
}
