package infra

import(
	"github.com/Miyagawa-Ryohei/gode_conf"
	"os"
)

type QueueConfig struct {
	URL string
}

type SessionConfig struct {
	Endpoint string
	Region string
	RoleArn string
}

type Config struct {
	Queue QueueConfig
	Session SessionConfig
}

type ConfigLoader struct {
	Name *string
	Dir *string
}

func (c ConfigLoader) Load(conf interface{}){
	gode_conf.LoadTo(conf, &gode_conf.ConfigOption{
		FileName:  os.Getenv("MK_MSC_ENV"),
	})
}
