package types

type Application interface {
	PushResultMessage(result []byte)
	PutResultFile(name string, root string, data []byte)
}

type Subscriber interface {
	Listen()
}
