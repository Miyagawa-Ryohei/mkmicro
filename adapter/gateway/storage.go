package gateway

import "github.com/Miyagawa-Ryohei/mkmicro/entity"

type StorageProxy struct {
	session entity.StorageSessionUpdater
	driver  entity.StorageDriver
}

func (q *StorageProxy) Update() {
	d, err := q.session.UpdateStorage()
	if err != nil {
		panic(err)
	}
	q.driver = d
}

func (d *StorageProxy) Get(root string, path string) ([]byte, error){
	return d.driver.Get(root,path)
}
func (d *StorageProxy) Put(root string, path string, raw []byte) error{
	return d.driver.Put(root,path,raw)
}

func NewStorageProxy(session entity.StorageSessionUpdater) (entity.StorageDriver, error) {
	d, err := session.UpdateStorage()
	if err != nil {
		return nil, err
	}
	return &StorageProxy{
		session: session,
		driver:  d,
	}, nil
}
