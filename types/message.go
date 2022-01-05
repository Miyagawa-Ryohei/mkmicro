package types

type Message interface {
	DeletableMessage
	ChangeVisibilityMessage
	GetBody() []byte
}

type DeletableMessage interface {
	GetDeleteID() string
	IsDeleted() bool
	SetDeleted(deleted bool)
}

type ChangeVisibilityMessage interface {
	GetChangeVisibilityID() string
}
