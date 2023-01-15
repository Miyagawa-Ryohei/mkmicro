package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Miyagawa-Ryohei/mkmicro/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"testing"
)

var cli *s3.Client

type TestCredentialProvider struct{}

func (p TestCredentialProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     "dummy",
		SecretAccessKey: "dummy",
	}, nil
}

type TestEndpointResolver struct{}

func (r TestEndpointResolver) ResolveEndpoint(service, region string, options ...interface{}) (aws.Endpoint, error) {
	return aws.Endpoint{
		URL: "http://localhost:4566",
	}, nil
}

const (
	testBucket     = "test-bucket"
	testDataPrefix = "ListTest"
)

func init() {

	c, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(TestCredentialProvider{}),
		config.WithEndpointResolverWithOptions(TestEndpointResolver{}),
	)
	if err != nil {
		log.Fatalln(err)
	}

	cli = s3.NewFromConfig(c, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	if _, err := cli.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(testBucket),
	}); err != nil {
		log.Fatalln(err)
	}

	for i := 1; i <= 2100; i++ {
		buf := bytes.NewBufferString("a")
		fname := aws.String(fmt.Sprintf("%s/test_data_%d", testDataPrefix, i))
		if _, err := cli.PutObject(
			context.TODO(),
			&s3.PutObjectInput{
				Bucket:        aws.String(testBucket),
				Key:           fname,
				Body:          buf,
				ContentLength: int64(buf.Len()),
			},
			s3.WithAPIOptions(
				v4.SwapComputePayloadSHA256ForUnsignedPayloadMiddleware, // 2.
			),
		); err != nil {
			log.Fatalln(err)
		}
	}
}

func TestS3Driver_List(t *testing.T) {
	driver := NewS3Driver(cli, &types.StorageConfig{})
	list, err, next := driver.List(testBucket, testDataPrefix)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) > 0 {
		for {
			l, err := next()
			if err != nil {
				t.Fatal(err)
			}
			if len(l) == 0 {
				break
			}
			list = append(list, l...)
		}
	}
	if len(list) != 2100 {
		t.Logf("expected 3500, but actual %d", len(list))
		t.Fail()
	}
}
