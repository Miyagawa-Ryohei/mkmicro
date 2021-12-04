package entity

type Message interface{
	DeletableMessage
	ChangeVisibilityMessage
	GetBody() []byte
}

type DeletableMessage interface {
	GetDeleteID() string
}

type ChangeVisibilityMessage interface {
	GetChangeVisibilityID() string
}