package types

type QueueDriver interface {
	GetConfig() *QueueConfig
	GetMessage(num int) ([]Message, error)
	PutMessage(raw []byte) error
	DeleteMessage(msg DeletableMessage) error
	GetMessageLength() ([]string, error)
	ChangeMessageVisibility(msg ChangeVisibilityMessage) error
}
