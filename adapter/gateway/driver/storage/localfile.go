package storage

import (
	"github.com/Miyagawa-Ryohei/mkmicro/types"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type LocalFileDriver struct {
	config *types.StorageConfig
}

func (d *LocalFileDriver) GetConfig() *types.StorageConfig {
	return d.config
}

func (d *LocalFileDriver) Put(bucket string, key string, data []byte) error {
	fpath := path.Join(bucket, key)
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}

func (d *LocalFileDriver) Get(bucket string, key string) ([]byte, error) {
	fpath := path.Join(bucket, key)
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func (d *LocalFileDriver) GetByStream(bucket string, key string) (io.ReadCloser, error) {
	fpath := path.Join(bucket, key)
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (d *LocalFileDriver) Download(bucket string, key string, dist string) error {
	fpath := path.Join(bucket, key)
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	tpath := path.Join(dist, key)
	t, err := os.OpenFile(tpath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	if _, err := io.Copy(t, f); err != nil {
		return err
	}
	return nil
}

func NewLocalFileDriver(config *types.StorageConfig) *LocalFileDriver {
	return &LocalFileDriver{
		config: config,
	}
}
