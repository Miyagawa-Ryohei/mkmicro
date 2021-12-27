package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Miyagawa-Ryohei/mkmicro/entity"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Driver struct {
	s3 *s3.Client
	config *entity.StorageConfig
}
func (d *S3Driver) GetConfig( ) *entity.StorageConfig {
	return d.config
}

func (d *S3Driver) Put(bucket string, key string, data []byte) (error) {
	param := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key : aws.String(key),
		Body : bytes.NewReader(data),
		ContentLength: int64(len(data)),
		ContentType:   aws.String("plain/text"),
	}

	fmt.Printf("%+v",*param)
	lsParam := &s3.ListBucketsInput{}
	resps, err := d.s3.ListBuckets(context.TODO(), lsParam)
	for _,b := range resps.Buckets {
		fmt.Printf("%s\n", *b.Name)
	}
	if err != nil {
		return err
	}

	if _, err := d.s3.PutObject(context.TODO(), param); err != nil {
		return err
	}
	return nil
}

func (d *S3Driver) Get(bucket string, key string) ([]byte, error) {
	param := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key : aws.String(key),
	}
	resp, err := d.s3.GetObject(context.TODO(), param)
	if err != nil {
		return nil, err
	}

	buf := []byte{}
	if _, err := resp.Body.Read(buf); err != nil {
		return nil, err
	}
	return buf, nil
}


func NewS3Driver (s *s3.Client, config *entity.StorageConfig) *S3Driver {
	return &S3Driver{
		s3: s,
		config: config,
	}
}