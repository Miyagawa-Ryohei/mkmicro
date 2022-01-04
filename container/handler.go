package container

import "github.com/Miyagawa-Ryohei/mkmicro/types"

type handlerContainer struct {
	handlers []types.Handler
}

var container *handlerContainer = nil

func (c *handlerContainer) Add(h types.Handler) {
	c.handlers = append(c.handlers, h)
}

func (r *handlerContainer) Get() []types.Handler {
	return r.handlers
}

func GetHandlerContainer() types.HandlerContainer {
	if container == nil {
		container = &handlerContainer{
			handlers: []types.Handler{},
		}
	}
	return container
}
