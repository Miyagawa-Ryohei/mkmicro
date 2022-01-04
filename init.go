package main

import (
	"github.com/Miyagawa-Ryohei/mkmicro/container"
	"github.com/Miyagawa-Ryohei/mkmicro/handler"
)

func Initialize() {
	handlers := container.GetHandlerContainer()
	handlers.Add(handler.SampleHandler{})
}
