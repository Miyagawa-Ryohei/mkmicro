package storage

import (
	"bytes"
	"context"
	"github.com/Miyagawa-Ryohei/mkmicro/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Driver struct {
	s3     *s3.Client
	config *types.StorageConfig
}

func (d *S3Driver) GetConfig() *types.StorageConfig {
	return d.config
}

func (d *S3Driver) Put(bucket string, key string, data []byte) error {
	param := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentLength: int64(len(data)),
	}

	if _, err := d.s3.PutObject(context.TODO(), param); err != nil {
		return err
	}
	return nil
}

func (d *S3Driver) Get(bucket string, key string) ([]byte, error) {
	param := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	resp, err := d.s3.GetObject(context.TODO(), param)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func NewS3Driver(s *s3.Client, config *types.StorageConfig) *S3Driver {
	return &S3Driver{
		s3:     s,
		config: config,
	}
}
