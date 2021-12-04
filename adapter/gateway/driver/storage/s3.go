package storage

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Driver struct {
	s3 *s3.S3
}

func (d *S3Driver) Put(bucket string, key string, data []byte) (error) {
	param := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key : aws.String(key),
		Body : bytes.NewReader(data),
	}
	if _, err := d.s3.PutObject(param); err != nil {
		return err
	}
	return nil
}

func (d *S3Driver) Get(bucket string, key string) ([]byte, error) {
	param := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key : aws.String(key),
	}
	resp, err := d.s3.GetObject(param)
	if err != nil {
		return nil, err
	}

	buf := []byte{}
	if _, err := resp.Body.Read(buf); err != nil {
		return nil, err
	}
	return buf, nil
}


func NewS3Driver (s *s3.S3) *S3Driver {
	return &S3Driver{
		s3: s,
	}
}