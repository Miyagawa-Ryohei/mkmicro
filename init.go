package main

import (
	"mkmicro/handler"
	"mkmicro/repository"
	"mkmicro/router"
	"net/http"
)

func Initialize () {
	routeRepo := repository.GetRouteRepository()
	routeRepo.Add("/check",http.MethodGet,router.HealthCheck)
	routeRepo.Add("/sample", http.MethodPost, router.SampleRouter)
	handlers := repository.GetHandlerRepository()
	handlers.Add(handler.SampleHandler{})
}