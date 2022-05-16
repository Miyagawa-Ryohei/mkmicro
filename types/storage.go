package types

type StorageDriver interface {
	GetConfig() *StorageConfig
	Get(root string, path string) ([]byte, error)
	Download(bucket string, key string, dist string) error
	Put(root string, path string, raw []byte) error
}
