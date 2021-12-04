package entity

import "github.com/Miyagawa-Ryohei/mkmicro/infra"

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
	CreateWithConfig(queue infra.QueueConfig, session infra.SessionConfig) SessionManager
}