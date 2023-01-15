package storage

import (
	"bytes"
	"context"
	"github.com/Miyagawa-Ryohei/mkmicro/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"io"
	"os"
	"path"
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

func (d *S3Driver) GetByStream(bucket string, key string) (io.ReadCloser, error) {
	param := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	resp, err := d.s3.GetObject(context.TODO(), param)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (d *S3Driver) Download(bucket string, key string, dist string) error {
	param := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	resp, err := d.s3.GetObject(context.TODO(), param)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	dir, _ := path.Split(dist)
	if dir != "" {
		if fi, err := os.Stat(dir); os.IsNotExist(err) || !fi.IsDir() {
			if e := os.MkdirAll(dir, 0755); e != nil {
				return e
			}
		}
	}

	f, err := os.OpenFile(dist, os.O_CREATE|os.O_WRONLY, 0755)
	defer f.Close()
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

func (d *S3Driver) List(bucket string, prefix string) (list []string, err error, next func() ([]string, error)) {
	t := uuid.New()

	parseToStringArray := func(contents []s3Types.Object) (ret []string) {
		ret = []string{}
		for _, object := range contents {
			ret = append(ret, *object.Key)
		}
		return
	}
	param := &s3.ListObjectsV2Input{
		Bucket:            aws.String(bucket),
		Prefix:            aws.String(prefix),
		ContinuationToken: aws.String(t.String()),
	}

	resp, err := d.s3.ListObjectsV2(context.TODO(), param)
	if err != nil {
		return nil, err, nil
	}

	return parseToStringArray(resp.Contents), nil, func() ([]string, error) {
		if resp.IsTruncated != true {
			return []string{}, nil
		}
		param := &s3.ListObjectsV2Input{
			Bucket:            aws.String(bucket),
			Prefix:            aws.String(prefix),
			ContinuationToken: resp.NextContinuationToken,
			StartAfter:        resp.Contents[len(resp.Contents)-1].Key,
		}
		resp, err = d.s3.ListObjectsV2(context.TODO(), param)
		if err != nil {
			return nil, err
		}
		return parseToStringArray(resp.Contents), err
	}
}

func NewS3Driver(s *s3.Client, config *types.StorageConfig) *S3Driver {
	return &S3Driver{
		s3:     s,
		config: config,
	}
}
