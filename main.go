package main

import (
	"flag"
	"fmt"
	"time"

	"bytes"
	"compress/gzip"
	"crypto/tls"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"io"

	"net/http"
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

func gUnzipData(data []byte) (resData []byte, err error) {
	b := bytes.NewBuffer(data)

	var r io.Reader
	r, err = gzip.NewReader(b)
	if err != nil {
		return
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return
	}

	resData = resB.Bytes()

	return
}

func followShard(svc *kinesis.Kinesis, sleepBetweenPolling int, shardIterator string, unzipData bool) {
	for true {
		time.Sleep(time.Duration(sleepBetweenPolling) * time.Millisecond)
		getRecordsOut, err := getRecords(svc, shardIterator)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		shardIterator = *getRecordsOut.NextShardIterator
		for _, el := range getRecordsOut.Records {
			if unzipData {
				unzippedBytes, err := gUnzipData(el.Data)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				s := string(unzippedBytes)
				fmt.Println(s)
			} else {
				s := string(el.Data)
				fmt.Println(s)
			}
		}
	}
}

func main() {
	streamName := flag.String("stream", "", "Name of the string you want to follow")
	region := flag.String("region", "eu-west-1", "AWS region the stream exists in")
	endpoint := flag.String("endpoint", "", "AWS endpoint (useful with localstack)")
	sleepBetweenPolling := flag.Int("sleep", 1000, "How long to sleep between polling for new messages in ms")
	noVerifySSL := flag.Bool("no-verify-ssl", false, "No SSL certificate validation")
	unzipData := flag.Bool("unzip-data", false, "Do you want to unzip the data bytes read from Kinesis")

	flag.Parse()

	if *streamName == "" {
		fmt.Println("Must provide and argument stream")
		return
	}

	if *noVerifySSL {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	sess := session.New(&aws.Config{
		Region:   aws.String(*region),
		Endpoint: aws.String(*endpoint),
	})
	svc := kinesis.New(sess)

	stream, err := getStreamDescription(svc, *streamName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, shard := range stream.Shards {
		shardIterator, err2 := getShardIterator(svc, *streamName, *shard.ShardId)
		if err2 != nil {
			fmt.Println(err2.Error())
			return
		}
		go followShard(svc, *sleepBetweenPolling, shardIterator, *unzipData)
	}

	var input string
	fmt.Scanln(&input)
}
