package entity

type StorageDriver interface {
	Get(root string, path string) ([]byte, error)
	Put(root string, path string, raw []byte) error
}