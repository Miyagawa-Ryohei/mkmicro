package queue

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"mkmicro/entity"
	"mkmicro/infra"
)

type SQSConfig struct {
	url string
}

type SQSDriver struct {
	url  string
	queue *sqs.SQS
}

type SQSMessage struct {
	raw *sqs.Message
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

func (d *SQSDriver) PutMessage(raw []byte) (error) {

	params := &sqs.SendMessageInput{
		MessageBody:  aws.String(string(raw)),
		QueueUrl:     aws.String(d.url),
		DelaySeconds: aws.Int64(1),
	}

	if _ , err := d.queue.SendMessage(params); err != nil {
		return err
	}
	return nil
}

func (d *SQSDriver) parseMessage(msgs []*sqs.Message) []entity.Message {
	ret := []entity.Message{}
	for _, m := range msgs {
		ret = append(ret, &SQSMessage{raw: m})
	}
	return ret
}

func (d *SQSDriver) GetMessage(num int) ([]entity.Message, error) {
	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(d.url),
		// 一度に取得する最大メッセージ数。最大でも10まで。
		MaxNumberOfMessages: aws.Int64(int64(num)),
		// これでキューが空の場合はロングポーリング(20秒間繋ぎっぱなし)になる。
		WaitTimeSeconds: aws.Int64(20),

		VisibilityTimeout: aws.Int64(60),
	}
	resp, err := d.queue.ReceiveMessage(params)

	if err != nil {
		return nil, err
	}

	if len(resp.Messages) == 0 {
		return nil, nil
	}

	return d.parseMessage(resp.Messages), nil
}

func (d *SQSDriver) DeleteMessage(msg entity.DeletableMessage) (error) {
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(d.url),
		ReceiptHandle: aws.String(msg.GetDeleteID()),
	}
	_, err := d.queue.DeleteMessage(params)

	if err != nil {
		return err
	}
	return nil
}

func (d *SQSDriver) ChangeMessageVisibility(msg entity.ChangeVisibilityMessage) (error) {
	params := &sqs.ChangeMessageVisibilityInput{
		QueueUrl:      aws.String(d.url),
		ReceiptHandle: aws.String(msg.GetChangeVisibilityID()),
		VisibilityTimeout : aws.Int64(60),
	}
	_, err := d.queue.ChangeMessageVisibility(params)

	if err != nil {
		return err
	}
	return nil
}

func NewSQSDriver (q *sqs.SQS, config infra.QueueConfig) *SQSDriver {
	return &SQSDriver{
		queue: q,
		url : config.URL,
	}
}