package infra

import (
	"github.com/Miyagawa-Ryohei/gode_conf"
	"log"
	"os"
)


type ConfigLoader struct {
	Name *string
	Dir *string
}

func (c ConfigLoader) Load(conf interface{}){
	if err := gode_conf.LoadTo(conf, &gode_conf.ConfigOption{
		FileName:  os.Getenv("MK_MSC_ENV"),
	}); err != nil {
		log.Fatal(err)
	}
}
