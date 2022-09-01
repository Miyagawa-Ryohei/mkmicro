package queue

import (
	"github.com/Miyagawa-Ryohei/mkmicro/types"
)

type SQSMockDriver struct {
	url           string
	config        *types.QueueConfig
	dummyMessage  *SQSMessage
	messageLength []string
}

func (d *SQSMockDriver) SetDummyMessage(m *SQSMessage) {
	d.dummyMessage = m
}

func (d *SQSMockDriver) GetConfig() *types.QueueConfig {
	return d.config
}

func (d *SQSMockDriver) PutMessage(raw []byte, delay int32) error {
	return nil
}

func (d *SQSMockDriver) SetMessageLength(lens []string) {
	d.messageLength = lens
}
func (d *SQSMockDriver) GetMessageLength() ([]string, error) {
	return d.messageLength, nil
}

func (d *SQSMockDriver) GetMessage(num int) ([]types.Message, error) {
	return []types.Message{d.dummyMessage}, nil
}

func (d *SQSMockDriver) DeleteMessage(msg types.DeletableMessage) error {
	return nil
}

func (d *SQSMockDriver) ChangeMessageVisibility(msg types.ChangeVisibilityMessage, second int32) error {
	return nil
}

func NewSQSMockDriver(config *types.QueueConfig) *SQSMockDriver {
	return &SQSMockDriver{
		url:    config.URL,
		config: config,
	}
}
