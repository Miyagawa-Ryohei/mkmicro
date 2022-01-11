package container

import "github.com/Miyagawa-Ryohei/mkmicro/types"

type HandlerContainer struct {
	handlers []types.Handler
}

var container *HandlerContainer = nil

func (c *HandlerContainer) Add(h types.Handler) {
	c.handlers = append(c.handlers, h)
}

func (c *HandlerContainer) Get() []types.Handler {
	return c.handlers
}

func GetHandlerContainer() types.HandlerContainer {
	if container == nil {
		container = &HandlerContainer{
			handlers: []types.Handler{},
		}
	}
	return container
}
