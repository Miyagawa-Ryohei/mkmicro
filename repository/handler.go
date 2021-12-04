package repository

import "github.com/Miyagawa-Ryohei/mkmicro/entity"

type Handler interface {
	Exec(msg entity.Message, dist entity.SessionManager) bool
	GetResultQueueConfig() *entity.QueueConfig
	GetResultSessionConfig() *entity.SessionConfig
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