package app

import (
	"log/slog"
	"net/http"

	"github.com/longwavee/effective-mobile-test/internal/api/rest"
	"github.com/longwavee/effective-mobile-test/internal/api/rest/middlewares"
	"github.com/longwavee/effective-mobile-test/internal/config"
)

type (
	DIContainer struct {
		cfg *config.Config
		log *slog.Logger

		httpServer *rest.HTTPServer
		httpRouter http.Handler
	}
)

func NewDIContainer(cfg *config.Config, log *slog.Logger) *DIContainer {
	return &DIContainer{
		cfg: cfg,
		log: log,
	}
}

func (c *DIContainer) HTTPServer() *rest.HTTPServer {
	if c.httpServer == nil {
		c.httpServer = rest.NewHTTPServer(
			&c.cfg.HTTPServer,
			c.HTTPRouter(),
		)
		c.log.Debug("http server initialized")
	}
	return c.httpServer
}

func (c *DIContainer) HTTPRouter() http.Handler {
	if c.httpRouter == nil {
		c.httpRouter = rest.NewHTTPRouter(
			middlewares.Logging(c.log),
			middlewares.Recovery(c.log),
		)
		c.log.Debug("http router initialized")
	}
	return c.httpRouter
}
