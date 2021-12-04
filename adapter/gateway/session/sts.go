package session

import (
	"github.com/Miyagawa-Ryohei/mkmicro/adapter/gateway"
	"github.com/Miyagawa-Ryohei/mkmicro/adapter/gateway/driver/queue"
	"github.com/Miyagawa-Ryohei/mkmicro/adapter/gateway/driver/storage"
	"github.com/Miyagawa-Ryohei/mkmicro/entity"
	"github.com/aws/aws-sdk-go/service/sqs"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type STSManager struct{
	queueConfig entity.QueueConfig
	sessionConfig entity.SessionConfig
	sess *session.Session
	sSess *session.Session
}

func (s *STSManager) UpdateSession() {
	if s.sessionConfig.RoleArn == "" {
		s.sSess = s.sess
		return
	}
	optionProvider := func(p *stscreds.AssumeRoleProvider) {
		p.Duration = time.Duration(12) * time.Hour
	}
	sCreds := stscreds.NewCredentials(s.sess, s.sessionConfig.RoleArn, optionProvider)

	sConfig := &aws.Config{
		Region: s.sess.Config.Region,
		Credentials: sCreds,
	}
	sSess := session.Must(session.NewSession(sConfig))

	s.sSess = sSess
}

func (s *STSManager) GetQueue() (entity.QueueDriver, error) {
	return gateway.NewQueueProxy(s)
}

func (s *STSManager) UpdateQueue() (entity.QueueDriver, error) {
	d := sqs.New(s.sess)
	return queue.NewSQSDriver(d, s.queueConfig), nil
}

func (s *STSManager) GetStorage() (entity.StorageDriver, error) {
	return gateway.NewStorageProxy(s)
}

func (s *STSManager) UpdateStorage() (entity.StorageDriver, error) {
	d := s3.New(s.sSess)
	return storage.NewS3Driver(d), nil
}

func newSTSSTSManager (queueConfig entity.QueueConfig, sessionConfig entity.SessionConfig) *STSManager {

	awsConfig := &aws.Config{}
	if sessionConfig.Endpoint != "" {
		awsConfig.Endpoint = aws.String(sessionConfig.Endpoint)
	}

	if sessionConfig.Region == "" {
		awsConfig.Region = aws.String("ap-northeast-1")
	} else {
		awsConfig.Region = aws.String(sessionConfig.Region)
	}


	sess, err := session.NewSession(awsConfig)
	if err != nil {
		panic(err)
	}

	s := &STSManager{
		queueConfig : queueConfig,
		sessionConfig : sessionConfig,
		sess: sess,
	}
	s.UpdateSession()
	return s
}

type STSManagerFactory struct {
	queue entity.QueueConfig
	sess entity.SessionConfig
}

func(f STSManagerFactory) Create () entity.SessionManager {
	return newSTSSTSManager(f.queue, f.sess)
}

func(f STSManagerFactory) CreateWithConfig(queue entity.QueueConfig, sess entity.SessionConfig) entity.SessionManager {
	return newSTSSTSManager(queue, sess)
}

func NewSTSManagerFactory(queue entity.QueueConfig, sess entity.SessionConfig) STSManagerFactory {
	return STSManagerFactory{
		queue: queue,
		sess: sess,
	}
}