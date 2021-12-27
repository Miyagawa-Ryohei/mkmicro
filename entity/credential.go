package entity

type SessionManager interface {
	UpdateSession()
	GetQueue() (QueueDriver, error)
	GetStorage() (StorageDriver, error)
	QueueSessionUpdater
	StorageSessionUpdater
}

type QueueSessionUpdater interface {
	UpdateQueue(cfg *QueueConfig) (QueueDriver, error)
}

type StorageSessionUpdater interface {
	UpdateStorage(cfg *StorageConfig) (StorageDriver, error)
}

type SessionManagerFactory interface {
	Create() (SessionManager, error)
	CreateWithConfig(queue QueueConfig, session StorageConfig) (SessionManager, error)
}