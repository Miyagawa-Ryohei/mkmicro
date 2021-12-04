package app

import (
	"github.com/labstack/echo/v4"
	"mkmicro/infra"
)

type Server struct {
	echo *echo.Echo
	group *echo.Group
}

func (s *Server) Listen() {
	s.echo.Start(":5678")
}

func (s *Server) Init() {
	repo := infra.GetRouteRepository()
	handlerMap := repo.Get()
	for path, methods := range handlerMap {
		for method, handler := range methods {
			s.group.Add(method,path,handler)
		}
	}
}

func NewServer () *Server {
	e := echo.New()
	return &Server {
		echo: e,
		group: e.Group("/api"),
	}
}