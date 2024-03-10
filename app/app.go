package app

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Echo *echo.Echo
}

func (s *Server) Init() error {
	s.Echo = echo.New()

	s.Echo.Use(middleware.Logger())
	s.Echo.Use(middleware.Recover())

	s.Echo.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	group := s.Echo.Group("/v0")
	group.GET("/who", func(c echo.Context) error {
		hostname, err := os.Hostname()
		if err != nil {
			return c.String(500, "error")
		}
		return c.String(200, hostname)
	})

	return nil
}

func (s *Server) Start() error {
	return s.Echo.Start(":6060")
}

func (s *Server) Stop() error {
	return s.Echo.Close()
}
