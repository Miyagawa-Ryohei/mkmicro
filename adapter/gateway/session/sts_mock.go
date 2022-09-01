package session

import (
	"github.com/Miyagawa-Ryohei/mkmicro/adapter/gateway/driver/queue"
	"github.com/Miyagawa-Ryohei/mkmicro/adapter/gateway/driver/storage"
	"github.com/Miyagawa-Ryohei/mkmicro/types"
)

type STSMockManager struct {
	queue         types.QueueDriver
	queueConfig   *types.QueueConfig
	storage       types.StorageDriver
	storageConfig *types.StorageConfig
}

func (s *STSMockManager) UpdateSession() {
}

func (s *STSMockManager) GetQueue() (types.QueueDriver, error) {
	return s.queue, nil
}

func (s *STSMockManager) CreateQueueWithConfig(cfg types.QueueConfig) (types.QueueDriver, error) {
	return s.queue, nil
}

func (s *STSMockManager) UpdateQueue(cfg *types.QueueConfig) (types.QueueDriver, error) {
	return s.queue, nil
}

func (s *STSMockManager) GetStorage() (types.StorageDriver, error) {
	return s.storage, nil
}

func (s *STSMockManager) CreateStorageWithConfig(cfg types.StorageConfig) (types.StorageDriver, error) {
	return s.storage, nil
}

func (s *STSMockManager) UpdateStorage(cfg *types.StorageConfig) (types.StorageDriver, error) {
	return s.storage, nil
}

func newSTSMockManager(queueConfig types.QueueConfig, storageConfig types.StorageConfig) (types.SessionManager, error) {
	return &STSMockManager{
		queue:         queue.NewSQSMockDriver(&queueConfig),
		storage:       storage.NewLocalFileDriver(&storageConfig),
		queueConfig:   &queueConfig,
		storageConfig: &storageConfig,
	}, nil
}

type STSMockManagerFactory struct {
	queue   types.QueueConfig
	storage types.StorageConfig
}

func (f STSMockManagerFactory) Create() (types.SessionManager, error) {
	return newSTSMockManager(f.queue, f.storage)
}

func (f STSMockManagerFactory) CreateWithConfig(queue types.QueueConfig, storage types.StorageConfig) (types.SessionManager, error) {
	return newSTSMockManager(queue, storage)
}

func NewSTSMockManagerFactory(queue types.QueueConfig, storage types.StorageConfig) STSManagerFactory {
	return STSManagerFactory{
		queue:   queue,
		storage: storage,
	}
}
