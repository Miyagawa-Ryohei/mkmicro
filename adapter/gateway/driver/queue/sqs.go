package queue

import (
	"bytes"
	"context"
	"github.com/Miyagawa-Ryohei/mkmicro/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSConfig struct {
	url string
}

type SQSDriver struct {
	url    string
	queue  *sqs.Client
	config *types.QueueConfig
}

type SQSMessage struct {
	raw     *awsTypes.Message
	deleted bool
}

func (m *SQSMessage) GetBody() []byte {
	return bytes.NewBufferString(*m.raw.Body).Bytes()
}

func (m *SQSMessage) GetChangeVisibilityID() string {
	return *m.raw.ReceiptHandle
}

func (m *SQSMessage) GetDeleteID() string {
	return *m.raw.ReceiptHandle
}

func (m *SQSMessage) SetDeleted(deleted bool) {
	m.deleted = deleted
}

func (m *SQSMessage) IsDeleted() bool {
	return m.deleted
}

func (d *SQSDriver) GetConfig() *types.QueueConfig {
	return d.config
}
func (d *SQSDriver) PutMessage(raw []byte) error {

	params := &sqs.SendMessageInput{
		MessageBody:  aws.String(string(raw)),
		QueueUrl:     aws.String(d.url),
		DelaySeconds: 1,
	}

	if _, err := d.queue.SendMessage(context.TODO(), params); err != nil {
		return err
	}
	return nil
}

func (d *SQSDriver) parseMessage(msgs []awsTypes.Message) []types.Message {
	ret := []types.Message{}
	for _, m := range msgs {
		ret = append(ret, &SQSMessage{
			raw:     &m,
			deleted: false,
		})
	}
	return ret
}

func (d *SQSDriver) GetMessage(num int) ([]types.Message, error) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(d.url),
		MaxNumberOfMessages: int32(num),
		WaitTimeSeconds:     20,
		VisibilityTimeout:   60,
	}
	resp, err := d.queue.ReceiveMessage(context.TODO(), params)

	if err != nil {
		return nil, err
	}

	if len(resp.Messages) == 0 {
		return nil, nil
	}

	return d.parseMessage(resp.Messages), nil
}

func (d *SQSDriver) DeleteMessage(msg types.DeletableMessage) error {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(d.url),
		ReceiptHandle: aws.String(msg.GetDeleteID()),
	}
	_, err := d.queue.DeleteMessage(context.TODO(), params)
	msg.SetDeleted(true)

	if err != nil {
		return err
	}
	return nil
}

func (d *SQSDriver) ChangeMessageVisibility(msg types.ChangeVisibilityMessage) error {
	params := &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(d.url),
		ReceiptHandle:     aws.String(msg.GetChangeVisibilityID()),
		VisibilityTimeout: 60,
	}
	_, err := d.queue.ChangeMessageVisibility(context.TODO(), params)

	if err != nil {
		return err
	}
	return nil
}

func NewSQSDriver(q *sqs.Client, config *types.QueueConfig) *SQSDriver {
	return &SQSDriver{
		queue:  q,
		url:    config.URL,
		config: config,
	}
}
