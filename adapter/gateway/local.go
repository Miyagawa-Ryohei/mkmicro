package gateway

import (
	"io"
	"os"
	"path"
)

type LocalFileSystem struct {

}

func(r *LocalFileSystem) tempDir() string {
	return 	os.TempDir()
}


func(r *LocalFileSystem) Read(name string) (*os.File, error) {
	abs := path.Join(r.tempDir(),name)
	return os.Open(abs)
}

func(r *LocalFileSystem) Write(name string, obj io.Reader) error {
	abs := path.Join(r.tempDir(),name)
	buf := []byte{}
	n, err := obj.Read(buf)
	if err != nil {
		return err
	}
	if n == 0 {
		return nil
	}
	if err := os.WriteFile(abs, buf,755); err != nil {
		return err
	}
	return nil
}

