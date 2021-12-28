package infra

import (
	"github.com/Miyagawa-Ryohei/gode_conf"
	"log"
)


type ConfigLoader struct {
	Name *string
	Dir *string
}

func (c ConfigLoader) Load(conf interface{}){
	if err := gode_conf.LoadTo(conf, &gode_conf.ConfigOption{
		FileName:  *c.Name,
		Directory: *c.Dir,
	}); err != nil {
		log.Fatal(err)
	}
}
