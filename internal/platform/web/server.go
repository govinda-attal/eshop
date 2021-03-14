package web

import (
	"context"
	"net/http"
	"strconv"
	"time"

	// oapigenmware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	echo          *echo.Echo
	mux           *http.ServeMux
	httpSrv       *http.Server
	shutDownGrace time.Duration
	basePath      string
}

type ServerOption func(*Server)

func ServerWithGraceShutdown(d time.Duration) ServerOption {
	return func(ws *Server) {
		ws.shutDownGrace = d
	}
}

func ServerWithHealthCheck(path string, f http.HandlerFunc) ServerOption {
	return func(ws *Server) {
		ws.mux.Handle(path, f)
	}
}

func ServerWithBasePath(path string) ServerOption {
	return func(ws *Server) {
		ws.basePath = path
	}
}

func ServerWithPrometheusHandler(path string) ServerOption {
	return func(ws *Server) {
		ws.mux.Handle(path, promhttp.Handler())
	}
}

func ServerWithApiSpec(swagger *openapi3.Swagger) ServerOption {
	return func(ws *Server) {
		// swagger.Servers = nil
		// ws.echo.Use(oapigenmware.OapiRequestValidator(swagger))
		ws.echo.GET(ws.basePath, func(ec echo.Context) error {
			return ec.JSON(http.StatusOK, swagger)
		})
	}
}

func (ws *Server) Echo() *echo.Echo {
	return ws.echo
}

func (ws *Server) Start(ctx context.Context) error {
	if err := ws.httpSrv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
	return nil
}

func (ws *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ws.shutDownGrace)
	defer cancel()

	if err := ws.httpSrv.Shutdown(ctx); err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func NewServer(ctx context.Context, port int, opts ...ServerOption) *Server {
	var (
		e   = echo.New()
		mux = http.NewServeMux()
		srv = &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: mux,
		}
		ws = &Server{
			echo:          e,
			mux:           mux,
			httpSrv:       srv,
			shutDownGrace: time.Second,
			basePath:      "/",
		}
	)
	mux.Handle(ws.basePath, e)
	for _, o := range opts {
		o(ws)
	}
	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = ErrorHandler
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(CtxMware())
	defCors := middleware.DefaultCORSConfig
	defCors.AllowMethods = append(defCors.AllowMethods, http.MethodOptions)
	e.Use(middleware.CORSWithConfig(defCors))
	return ws
}
