package entity

type SessionManager interface {
	UpdateSession()
	GetQueue() (QueueDriver, error)
	GetStorage() (StorageDriver, error)
	QueueSessionUpdater
	StorageSessionUpdater
}

type QueueSessionUpdater interface {
	UpdateQueue() (QueueDriver, error)
}

type StorageSessionUpdater interface {
	UpdateStorage() (StorageDriver, error)
}

type SessionManagerFactory interface {
	Create() SessionManager
	CreateWithConfig(queue QueueConfig, session SessionConfig) SessionManager
}