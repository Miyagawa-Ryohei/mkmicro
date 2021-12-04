package entity

type QueueDriver interface {
	GetMessage(num int) ([]Message, error)
	PutMessage(raw []byte) (error)
	DeleteMessage(msg DeletableMessage) error
	ChangeMessageVisibility(msg ChangeVisibilityMessage) (error)
}