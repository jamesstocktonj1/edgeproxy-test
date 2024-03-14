package app

import (
	"fmt"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func (s *Server) InitTracing() error {
	fmt.Println("Initializing tracing")
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	defcfg := config.Configuration{
		ServiceName: hostname,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
		Headers: (*jaeger.HeadersConfig).ApplyDefaults(&jaeger.HeadersConfig{}),
	}
	cfg, err := defcfg.FromEnv()
	if err != nil {
		return err
	}

	s.tracer, s.closer, err = cfg.NewTracer()
	if err != nil {
		return err
	}

	opentracing.SetGlobalTracer(s.tracer)

	return nil
}

func (s *Server) TracingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var sp opentracing.Span
		ctx, err := s.tracer.Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request().Header),
		)
		if err != nil {
			sp = s.tracer.StartSpan(c.Request().URL.Path)
			sp.SetTag("tracing.error", err.Error())
		} else {
			sp = s.tracer.StartSpan(c.Request().URL.Path, opentracing.ChildOf(ctx))
		}
		defer sp.Finish()

		sp.SetTag("http.client_ip", c.RealIP())
		sp.SetTag("http.method", c.Request().Method)
		sp.SetTag("http.url", c.Request().URL.Path)
		sp.SetTag("http.headers", c.Request().Header)

		reqSpan := c.Request().WithContext(opentracing.ContextWithSpan(c.Request().Context(), sp))
		c.SetRequest(reqSpan)

		s.tracer.Inject(
			sp.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request().Header),
		)

		err = next(c)

		sp.SetTag("http.status_code", c.Response().Status)
		if err != nil {
			sp.LogKV("error.message", err)
			sp.SetTag("error", true)
		}

		return err
	}
}

func CreateChildSpan(c echo.Context, name string) opentracing.Span {
	span := opentracing.SpanFromContext(c.Request().Context())
	return opentracing.StartSpan(
		name,
		opentracing.ChildOf(span.Context()),
	)
}
