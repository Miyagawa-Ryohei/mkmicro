package repository

import (
	"github.com/labstack/echo/v4"
	"github.com/Miyagawa-Ryohei/mkmicro/entity"
)

type RouteHandler func(msg entity.Message) []byte

type RouteRepository interface {
	Add(path string, method string, handler echo.HandlerFunc)
	Get()  map[string]map[string]echo.HandlerFunc
}

type routeRepository struct {
	router map[string]map[string]echo.HandlerFunc
}

var router *routeRepository= nil

func (r *routeRepository) Add (path string, method string, h echo.HandlerFunc) {
	if r.router[path] == nil {
		r.router[path] = map[string]echo.HandlerFunc{}
	}
	r.router[path][method] = h
}

func (r *routeRepository) Get () map[string]map[string]echo.HandlerFunc {
	return r.router
}

func GetRouteRepository () RouteRepository{
	if router == nil {
		router = &routeRepository{
			router: map[string]map[string]echo.HandlerFunc{},
		}
	}
	return router
}