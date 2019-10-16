package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
)

func manageS3(thisSession *session.Session) {
	s3svc := s3.New(thisSession, aws.NewConfig().WithRegion("us-east-1"))

	input := &s3.ListBucketsInput{}

	result, err := s3svc.ListBuckets(input)
	if err != nil {
		log.Error(err.Error())
	}

	log.Info(result)
}
