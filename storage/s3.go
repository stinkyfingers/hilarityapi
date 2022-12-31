package storage

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const bucket = "hilarity-jds"

// TODO
type S3 struct {
}

func (s *S3) GetQuestions() ([]string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1"),
	})
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)
	resp, err := svc.GetObject(&s3.GetObjectInput{
		Key:    aws.String(questionKey),
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, err
	}
	var output []string
	err = json.NewDecoder(resp.Body).Decode(&output)
	return output, err
}
