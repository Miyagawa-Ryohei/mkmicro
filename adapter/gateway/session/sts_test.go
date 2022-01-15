package session

import (
	"bytes"
	"github.com/Miyagawa-Ryohei/mkmicro/types"
	"testing"
	"time"
)

func Test_DefaultProfilePutAndGet(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	object := bytes.NewBufferString(t.Name()+now).Bytes()
	cfg := types.StorageConfig{
		Endpoint: &types.EndPoint{
			Region: "ap-northeast-1",
			URL:    "http://localhost:9000",
		},
	}
	f := NewSTSManagerFactory(types.QueueConfig{},cfg)
	s, err := f.Create()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("PutAndGetObject", func(t *testing.T){
		storage, err := s.GetStorage()
		if err != nil {
			t.Fatal(err)
		}
		if err := storage.Put("sample",t.Name(),object); err != nil {
			t.Fatal(err)
		}

		buf, err := storage.Get("sample",t.Name())
		if  err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(buf,object) != 0{
			t.Fatalf("getObject is mismatch putObject: put -> %s, get -> %s",object,buf)
		}
	})

}
func Test_AssumeRolePutAndGet(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	object := bytes.NewBufferString(t.Name()+now).Bytes()
	cfg := types.StorageConfig{
		Profile:    &types.Profile{
			AssumeRoleArn: "arn:aws:iam::845799411254:role/S3UserRole",
		},
	}
	f := NewSTSManagerFactory(types.QueueConfig{},cfg)
	s, err := f.Create()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("PutAndGetObjectWithAssumeRole", func(t *testing.T){
		storage, err := s.GetStorage()
		if err != nil {
			t.Fatal(err)
		}
		if err := storage.Put("mkmicrotest",t.Name(),object); err != nil {
			t.Fatal(err)
		}

		buf, err := storage.Get("mkmicrotest",t.Name())
		if  err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(buf,object) != 0{
			t.Fatalf("getObject is mismatch putObject: put -> %s, get -> %s",object,buf)
		}

	})

}