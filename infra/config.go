package infra

import (
	"github.com/Miyagawa-Ryohei/gode_conf"
	"log"
	"os"
)


type ConfigLoader struct {
}

func (c ConfigLoader) Load(conf interface{}){
	name := os.Getenv("MK_CONFIG_NAME")
	dir := os.Getenv("MK_CONFIG_DIR")
	if err := gode_conf.LoadTo(conf, &gode_conf.ConfigOption{
		FileName:  name,
		Directory: dir,
	}); err != nil {
		log.Fatal(err)
	}
}
