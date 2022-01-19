package session

import (
	"bytes"
	"github.com/Miyagawa-Ryohei/mkmicro/types"
	"os"
	"testing"
	"time"
)

func Test_DefaultProfilePutAndGet(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	object := bytes.NewBufferString(t.Name() + now).Bytes()
	cfg := types.StorageConfig{
		Endpoint: &types.EndPoint{
			Region: "ap-northeast-1",
			URL:    "http://localhost:9000",
		},
	}
	f := NewSTSManagerFactory(types.QueueConfig{}, cfg)
	s, err := f.Create()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("PutAndGetObject", func(t *testing.T) {
		storage, err := s.GetStorage()
		if err != nil {
			t.Fatal(err)
		}
		if err := storage.Put("sample", t.Name(), object); err != nil {
			t.Fatal(err)
		}

		buf, err := storage.Get("sample", t.Name())
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(buf, object) != 0 {
			t.Fatalf("getObject is mismatch putObject: put -> %s, get -> %s", object, buf)
		}
	})
}

func Test_AssumeRolePutAndGet(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	object := bytes.NewBufferString(t.Name() + now).Bytes()
	cfg := types.StorageConfig{
		Profile: &types.Profile{
			AssumeRoleArn: "arn:aws:iam::845799411254:role/S3UserRole",
		},
	}
	f := NewSTSManagerFactory(types.QueueConfig{}, cfg)
	s, err := f.Create()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("PutAndGetObjectWithAssumeRole", func(t *testing.T) {
		storage, err := s.GetStorage()
		if err != nil {
			t.Fatal(err)
		}
		if err := storage.Put("mkmicrotest", t.Name(), object); err != nil {
			t.Fatal(err)
		}

		buf, err := storage.Get("mkmicrotest", t.Name())
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(buf, object) != 0 {
			t.Fatalf("getObject is mismatch putObject: put -> %s, get -> %s", object, buf)
		}
	})
}

func Test_Queue(t *testing.T) {
	cfg := types.QueueConfig{
		URL: os.Getenv("AWS_QUEUE_URL"),
		Credential: &types.Credential{
			AccessKey:       os.Getenv("AWS_ACCESS_KEY_ID"),
			AccessKeySecret: os.Getenv("AWS_ACCESS_KEY_SECRET"),
		},
	}
	f := NewSTSManagerFactory(cfg, types.StorageConfig{})
	s, err := f.Create()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("PushMessage", func(t *testing.T) {
		q, err := s.GetQueue()
		if err != nil {
			t.Fatal(err)
		}
		if err := q.PutMessage(bytes.NewBufferString("HelloQueue").Bytes()); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("GetMessage", func(t *testing.T) {
		q, err := s.GetQueue()
		if err != nil {
			t.Fatal(err)
		}
		if err := q.PutMessage(bytes.NewBufferString("HelloQueue").Bytes()); err != nil {
			t.Fatal(err)
		}
		messages, err := q.GetMessage(10)
		if err != nil {
			t.Fatal(err)
		}
		if len(messages) == 0 {
			t.Fatal("message empty")
		}
		msg := string(messages[0].GetBody())
		if  msg != "HelloQueue" {
			t.Fatal("message expected : HelloQueue, but actual : " + msg)
		}
		for _, m := range messages {
			if err := q.DeleteMessage(m); err != nil {
				t.Fatal(m)
			}
		}
	})
}
