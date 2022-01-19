package types

type QueueDriver interface {
	GetConfig() *QueueConfig
	GetMessage(num int) ([]Message, error)
	PutMessage(raw []byte, delay int32) error
	DeleteMessage(msg DeletableMessage) error
	GetMessageLength() ([]string, error)
	ChangeMessageVisibility(msg ChangeVisibilityMessage, second int32) error
}
