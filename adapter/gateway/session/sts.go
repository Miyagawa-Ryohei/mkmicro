package session

import (
	"context"
	"fmt"
	"github.com/Miyagawa-Ryohei/mkmicro/adapter/gateway"
	"github.com/Miyagawa-Ryohei/mkmicro/adapter/gateway/driver/queue"
	"github.com/Miyagawa-Ryohei/mkmicro/adapter/gateway/driver/storage"
	"github.com/Miyagawa-Ryohei/mkmicro/types"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type CustomEndpointResolver struct {
	cfg types.AWSConfig
}

func (c CustomEndpointResolver) ResolveEndpoint(service, region string, options ...interface{}) (aws.Endpoint, error) {
	return aws.Endpoint{
		PartitionID:   "aws",
		URL:           c.cfg.Endpoint.URL,
		SigningRegion: c.cfg.Endpoint.Region,
		SigningMethod: "s3v4",
	}, nil
}

type CustomCredentialProvider struct {
	src types.AWSConfig
}

func (p CustomCredentialProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     p.src.Credential.AccessKey,
		SecretAccessKey: p.src.Credential.AccessKeySecret,
		SessionToken:    "",
		Source:          "",
		CanExpire:       false,
		Expires:         time.Time{},
	}, nil
}

type STSManager struct {
	queue         *types.QueueDriver
	queueConfig   *types.QueueConfig
	storage       *types.StorageDriver
	storageConfig *types.StorageConfig
}

func (s *STSManager) UpdateSession() {
}

func (s *STSManager) GetQueue() (types.QueueDriver, error) {
	if s.queue == nil {
		q, err := s.CreateQueueWithConfig(*s.queueConfig)
		if err != nil {
			return nil, err
		}
		s.queue = &q
	}
	return *s.queue, nil
}

func (s *STSManager) CreateQueueWithConfig(customConfig types.QueueConfig) (types.QueueDriver, error) {
	queueResolver := getResolvers(customConfig.GetAWSConfig())
	queueCfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		queueResolver...,
	)
	if err != nil {
		return nil, err
	}

	client := sqs.NewFromConfig(queueCfg)
	driver := queue.NewSQSDriver(client, &customConfig)
	proxy := gateway.NewQueueProxyWithDriverInstance(s, driver)

	return proxy, nil
}

func (s *STSManager) UpdateQueue(cfg *types.QueueConfig) (types.QueueDriver, error) {
	var client *sqs.Client
	if cfg != nil {
		queueResolver := getResolvers(cfg.GetAWSConfig())
		queueCfg, err := awsConfig.LoadDefaultConfig(
			context.TODO(),
			queueResolver...,
		)
		if err != nil {
			return nil, err
		}
		client = sqs.NewFromConfig(queueCfg)
	}
	driver := queue.NewSQSDriver(client, cfg)
	return driver, nil
}

func (s *STSManager) GetStorage() (types.StorageDriver, error) {
	if s.storage == nil {
		st, err := s.CreateStorageWithConfig(*s.storageConfig)
		if err != nil {
			return nil, err
		}
		s.storage = &st
	}
	return *s.storage, nil
}

func (s *STSManager) CreateStorageWithConfig(customConfig types.StorageConfig) (types.StorageDriver, error) {
	resolver := getResolvers(customConfig.GetAWSConfig())
	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		resolver...,
	)

	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	driver := storage.NewS3Driver(client, &customConfig)
	proxy := gateway.NewStorageProxyWithDriverInstance(s, driver)

	return proxy, nil
}

func (s *STSManager) UpdateStorage(cfg *types.StorageConfig) (types.StorageDriver, error) {
	var client *s3.Client
	if cfg == nil {
		return nil, fmt.Errorf("update session error( config is nil )")
	}
	storageResolver := getResolvers(cfg.GetAWSConfig())
	storageCfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		storageResolver...,
	)
	if err != nil {
		return nil, err
	}
	client = s3.NewFromConfig(storageCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	driver := storage.NewS3Driver(client, cfg)
	return driver, nil
}

func getResolvers(config types.AWSConfig) []func(*awsConfig.LoadOptions) error {
	r := CustomEndpointResolver{
		cfg: config,
	}
	resolvers := []func(*awsConfig.LoadOptions) error{}
	if config.Profile != nil {
		resolvers = append(resolvers, awsConfig.WithSharedConfigProfile(config.Profile.Name))
	} else if config.Credential != nil {
		p := CustomCredentialProvider{
			src: config,
		}
		resolvers = append(resolvers, awsConfig.WithCredentialsProvider(p), awsConfig())
	}
	if config.Endpoint != nil {
		resolvers = append(resolvers, awsConfig.WithEndpointResolverWithOptions(r))
	}
	return resolvers
}

func newSTSSTSManager(queueConfig types.QueueConfig, storageConfig types.StorageConfig) (*STSManager, error) {
	return &STSManager{
		queueConfig:   &queueConfig,
		storageConfig: &storageConfig,
	}, nil
}

type STSManagerFactory struct {
	queue   types.QueueConfig
	storage types.StorageConfig
}

func (f STSManagerFactory) Create() (types.SessionManager, error) {
	return newSTSSTSManager(f.queue, f.storage)
}

func (f STSManagerFactory) CreateWithConfig(queue types.QueueConfig, storage types.StorageConfig) (types.SessionManager, error) {
	return newSTSSTSManager(queue, storage)
}

func NewSTSManagerFactory(queue types.QueueConfig, storage types.StorageConfig) STSManagerFactory {
	return STSManagerFactory{
		queue:   queue,
		storage: storage,
	}
}
