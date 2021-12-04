package gateway

import (
	"github.com/Miyagawa-Ryohei/mkmicro/entity"
)

type QueueProxy struct {
	session entity.QueueSessionUpdater
	driver entity.QueueDriver
}

func (q *QueueProxy) Update() {
	d, err := q.session.UpdateQueue()
	if err != nil {
		panic(err)
	}
	q.driver = d
}

func (q *QueueProxy) GetMessage(num int) ([]entity.Message, error){
	resp, err := q.driver.GetMessage(num)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (q *QueueProxy) PutMessage(raw []byte) (error) {
	return q.driver.PutMessage(raw)
}

func (q *QueueProxy) DeleteMessage(msg entity.DeletableMessage) (error){
	return q.driver.DeleteMessage(msg)
}

func (q *QueueProxy) ChangeMessageVisibility(msg entity.ChangeVisibilityMessage) (error){
	return q.driver.ChangeMessageVisibility(msg)
}

func NewQueueProxy (session entity.QueueSessionUpdater) (entity.QueueDriver, error) {
	q,err := session.UpdateQueue()
	if err != nil {
		return nil, err
	}
	return &QueueProxy{
		session: session,
		driver: q,
	}, nil
}