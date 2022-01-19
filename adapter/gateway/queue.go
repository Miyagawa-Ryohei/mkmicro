package gateway

import (
	"github.com/Miyagawa-Ryohei/mkmicro/types"
)

type QueueProxy struct {
	session types.QueueSessionUpdater
	driver  types.QueueDriver
}

func (q *QueueProxy) GetConfig() *types.QueueConfig {
	return q.driver.GetConfig()
}
func (q *QueueProxy) Update() {
	d, err := q.session.UpdateQueue(q.driver.GetConfig())
	if err != nil {
		panic(err)
	}
	q.driver = d
}

func (q *QueueProxy) GetMessage(num int) ([]types.Message, error) {
	resp, err := q.driver.GetMessage(num)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (q *QueueProxy) GetMessageLength() ([]string, error) {
	return q.driver.GetMessageLength()
}

func (q *QueueProxy) PutMessage(raw []byte, delay int32) error {
	return q.driver.PutMessage(raw, delay)
}

func (q *QueueProxy) DeleteMessage(msg types.DeletableMessage) error {
	return q.driver.DeleteMessage(msg)
}

func (q *QueueProxy) ChangeMessageVisibility(msg types.ChangeVisibilityMessage) error {
	return q.driver.ChangeMessageVisibility(msg)
}

func NewQueueProxy(session types.QueueSessionUpdater) (types.QueueDriver, error) {
	q, err := session.UpdateQueue(nil)
	if err != nil {
		return nil, err
	}
	return &QueueProxy{
		session: session,
		driver:  q,
	}, nil
}

func NewQueueProxyWithDriverInstance(session types.QueueSessionUpdater, q types.QueueDriver) types.QueueDriver {
	return &QueueProxy{
		session: session,
		driver:  q,
	}
}
