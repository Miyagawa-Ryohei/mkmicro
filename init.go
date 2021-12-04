package main

import (
	"github.com/Miyagawa-Ryohei/mkmicro/handler"
	"github.com/Miyagawa-Ryohei/mkmicro/repository"
	"github.com/Miyagawa-Ryohei/mkmicro/router"
	"net/http"
)

func Initialize () {
	routeRepo := repository.GetRouteRepository()
	routeRepo.Add("/check",http.MethodGet,router.HealthCheck)
	routeRepo.Add("/sample", http.MethodPost, router.SampleRouter)
	handlers := repository.GetHandlerRepository()
	handlers.Add(handler.SampleHandler{})
}