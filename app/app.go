package app

import (
	"io"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/opentracing/opentracing-go"
)

type Server struct {
	Echo   *echo.Echo
	tracer opentracing.Tracer
	closer io.Closer
}

func (s *Server) Init() error {
	s.Echo = echo.New()

	s.Echo.Use(middleware.Logger())
	s.Echo.Use(middleware.Recover())

	err := s.InitTracing()
	if err != nil {
		return err
	}
	s.Echo.Use(s.TracingMiddleware)

	s.Echo.GET("/health", func(c echo.Context) error {
		sp := CreateChildSpan(c, "health")
		defer sp.Finish()

		return c.String(200, "OK")
	})

	group := s.Echo.Group("/v0")
	group.GET("/who", func(c echo.Context) error {
		sp := CreateChildSpan(c, "who")
		defer sp.Finish()

		hostname, err := os.Hostname()
		if err != nil {
			sp.SetTag("error", true)
			sp.SetTag("error.message", err.Error())
			return c.String(500, "error")
		}
		return c.String(200, hostname)
	})
	group.GET("/who-rand", s.WhoRand)

	return nil
}

func (s *Server) Start() error {
	return s.Echo.Start(":6060")
}

func (s *Server) Stop() error {
	s.closer.Close()
	return s.Echo.Close()
}

func (s *Server) WhoRand(c echo.Context) error {
	sp := CreateChildSpan(c, "who-rand")
	defer sp.Finish()

	hostname, err := os.Hostname()
	if err != nil {
		sp.SetTag("error", true)
		sp.SetTag("error.message", err.Error())
		return c.String(500, "error")
	}

	if rand.IntN(10) == 0 {
		return c.String(200, hostname)
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/vr/who-rand", nil)
	if err != nil {
		sp.SetTag("error", true)
		sp.SetTag("error.message", err.Error())
		return c.String(500, "error")
	}
	s.tracer.Inject(
		sp.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		sp.SetTag("error", true)
		sp.SetTag("error.message", err.Error())
		return c.String(500, "error")
	}
	defer resp.Body.Close()

	host, err := io.ReadAll(resp.Body)
	if err != nil {
		sp.SetTag("error", true)
		sp.SetTag("error.message", err.Error())
		return c.String(500, "error")
	}

	return c.String(200, hostname+" -> "+string(host))
}
