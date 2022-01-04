package types

type Container interface {
	Add(handler interface{})
	Get() []interface{}
}

type HandlerContainer interface {
	Add(handler Handler)
	Get() []Handler
}

type Handler interface {
	Exec(msg Message, dist SessionManager) error
}
