package gateway

import (
	"github.com/Miyagawa-Ryohei/mkmicro/types"
	"io"
)

type StorageProxy struct {
	session types.StorageSessionUpdater
	driver  types.StorageDriver
}

func (d *StorageProxy) GetConfig() *types.StorageConfig {
	return d.driver.GetConfig()
}
func (d *StorageProxy) Update() {
	driver, err := d.session.UpdateStorage(d.driver.GetConfig())
	if err != nil {
		panic(err)
	}
	d.driver = driver
}

func (d *StorageProxy) Get(root string, path string) ([]byte, error) {
	return d.driver.Get(root, path)
}
func (d *StorageProxy) Put(root string, path string, raw []byte) error {
	return d.driver.Put(root, path, raw)
}

func (d *StorageProxy) GetByStream(bucket string, key string) (io.Reader, error) {
	return d.driver.GetByStream(bucket, key)
}

func (d *StorageProxy) Download(bucket string, key string, dist string) error {
	return d.driver.Download(bucket, key, dist)
}

func NewStorageProxy(session types.StorageSessionUpdater) (types.StorageDriver, error) {
	d, err := session.UpdateStorage(nil)
	if err != nil {
		return nil, err
	}
	return &StorageProxy{
		session: session,
		driver:  d,
	}, nil
}

func NewStorageProxyWithDriverInstance(session types.StorageSessionUpdater, s types.StorageDriver) types.StorageDriver {
	return &StorageProxy{
		session: session,
		driver:  s,
	}
}
