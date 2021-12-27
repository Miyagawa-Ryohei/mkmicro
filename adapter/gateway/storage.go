package gateway

import (
	"github.com/Miyagawa-Ryohei/mkmicro/entity"
)

type StorageProxy struct {
	session entity.StorageSessionUpdater
	driver  entity.StorageDriver
}

func (d *StorageProxy) GetConfig() *entity.StorageConfig {
	return d.driver.GetConfig()
}
func (d *StorageProxy) Update() {
	driver, err := d.session.UpdateStorage(d.driver.GetConfig())
	if err != nil {
		panic(err)
	}
	d.driver = driver
}

func (d *StorageProxy) Get(root string, path string) ([]byte, error){
	return d.driver.Get(root,path)
}
func (d *StorageProxy) Put(root string, path string, raw []byte) error{
	return d.driver.Put(root,path,raw)
}

func NewStorageProxy(session entity.StorageSessionUpdater) (entity.StorageDriver, error) {
	d, err := session.UpdateStorage(nil)
	if err != nil {
		return nil, err
	}
	return &StorageProxy{
		session: session,
		driver:  d,
	}, nil
}

func NewStorageProxyWithDriverInstance (session entity.StorageSessionUpdater, s entity.StorageDriver) entity.StorageDriver {
	return &StorageProxy{
		session: session,
		driver: s,
	}
}
