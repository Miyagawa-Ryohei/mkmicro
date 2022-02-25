package queue

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
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

func (m *SQSMessage) GetDeduplicationID() string {
	body := m.GetBody()
	hash := sha256.New()
	hash.Write(body)
	v := hash.Sum(nil)
	return string(v)
}

func (m *SQSMessage) GetBody() []byte {
	return bytes.NewBufferString(*m.raw.Body).Bytes()
}

func (m *SQSMessage) GetID() string {
	return *m.raw.MessageId
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
func (d *SQSDriver) PutMessage(raw []byte, delay int32) error {

	params := &sqs.SendMessageInput{
		MessageBody:  aws.String(string(raw)),
		QueueUrl:     aws.String(d.url),
		DelaySeconds: delay,
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

func (d *SQSDriver) GetMessageLength() ([]string, error) {
	params := &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(d.url),
		AttributeNames: []awsTypes.QueueAttributeName{
			"ApproximateNumberOfMessages",
			"ApproximateNumberOfMessagesNotVisible",
			"ApproximateNumberOfMessagesDelayed",
		},
	}
	resp, err := d.queue.GetQueueAttributes(context.TODO(), params)

	if err != nil {
		return nil, err
	}

	if resp.Attributes == nil {
		return nil, nil
	}

	return []string{
		resp.Attributes["ApproximateNumberOfMessages"],
		resp.Attributes["ApproximateNumberOfMessagesNotVisible"],
		resp.Attributes["ApproximateNumberOfMessagesDelayed"],
	}, nil
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
	res, err := d.queue.DeleteMessage(context.TODO(), params)

	if err != nil {
		fmt.Printf("%+v", res)
		return err
	}

	msg.SetDeleted(true)
	return nil
}

func (d *SQSDriver) ChangeMessageVisibility(msg types.ChangeVisibilityMessage, second int32) error {
	params := &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(d.url),
		ReceiptHandle:     aws.String(msg.GetChangeVisibilityID()),
		VisibilityTimeout: second,
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
