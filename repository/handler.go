package repository

import "github.com/Miyagawa-Ryohei/mkmicro/entity"

type Handler interface {
	Exec(msg entity.Message, dist entity.SessionManager) error
	GetResultQueueConfig() *entity.QueueConfig
	GetResultStorageConfig() *entity.StorageConfig
}

type HandlerRepository interface {
	Add(handler Handler)
	Get() []Handler
}

type handlerRepository struct {
	handler []Handler
}

var repo *handlerRepository= nil

func (r *handlerRepository) Add (h Handler) {
	r.handler = append(r.handler, h)
}

func (r *handlerRepository) Get () []Handler {
	return r.handler
}

func GetHandlerRepository () HandlerRepository{
	if repo == nil {
		repo = &handlerRepository{
			handler: []Handler{},
		}
	}
	return repo
}