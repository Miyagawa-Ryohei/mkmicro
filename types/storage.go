package types

type StorageDriver interface {
	GetConfig() *StorageConfig
	Get(root string, path string) ([]byte, error)
	Put(root string, path string, raw []byte) error
}
