package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/longwavee/effective-mobile-test/internal/api/rest"
	"github.com/longwavee/effective-mobile-test/internal/api/rest/handlers"
	"github.com/longwavee/effective-mobile-test/internal/api/rest/middlewares"
	"github.com/longwavee/effective-mobile-test/internal/config"
	"github.com/longwavee/effective-mobile-test/internal/pkg/postgres"
	"github.com/longwavee/effective-mobile-test/internal/repository"
	"github.com/longwavee/effective-mobile-test/internal/service"
)

type (
	DIContainer struct {
		cfg *config.Config
		log *slog.Logger

		httpServer *rest.HTTPServer
		httpRouter http.Handler

		postgres *postgres.Client

		healthHandler *handlers.HealthHandler
		subsHandler   *handlers.SubscriptionHandler

		healthService *service.HealthService
		subsService   *service.SubscriptionService

		subscriptionRepo *repository.SubscriptionRepo
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
			c.HealthHandler(),
			c.SubscriptionHandler(),

			middlewares.Logging(c.log),
			middlewares.Recovery(c.log),
		)
		c.log.Debug("http router initialized")
	}
	return c.httpRouter
}

func (c *DIContainer) Postgres() *postgres.Client {
	if c.postgres == nil {
		client, err := postgres.New(
			context.TODO(),
			c.cfg.Postgres.ConnString(),
		)
		if err != nil {
			log.Println(fmt.Errorf("failed to init postgres: %w", err))
			os.Exit(1)
		}

		c.postgres = client
		c.log.Debug("postgres initialized")
	}
	return c.postgres
}

func (c *DIContainer) HealthHandler() *handlers.HealthHandler {
	if c.healthHandler == nil {
		c.healthHandler = handlers.NewHealthHandler(
			&c.cfg.HTTPServer,
			c.HealthService(),
			c.log,
		)
		c.log.Debug("health handler initialized")
	}
	return c.healthHandler
}

func (c *DIContainer) SubscriptionHandler() *handlers.SubscriptionHandler {
	if c.subsHandler == nil {
		c.subsHandler = handlers.NewSubscriptionHandler(
			c.SubscriptionService(),
			c.log,
		)
		c.log.Debug("subscription handler initialized")
	}
	return c.subsHandler
}

func (c *DIContainer) HealthService() *service.HealthService {
	if c.healthService == nil {
		c.healthService = service.NewHealthService(
			c.Postgres(),
		)
		c.log.Debug("health service initialized")
	}
	return c.healthService
}

func (c *DIContainer) SubscriptionService() *service.SubscriptionService {
	if c.subsService == nil {
		c.subsService = service.NewSubscriptionService(
			c.SubscriptionsRepo(),
		)
		c.log.Debug("subscription service initialized")
	}
	return c.subsService
}

func (c *DIContainer) SubscriptionsRepo() *repository.SubscriptionRepo {
	if c.subscriptionRepo == nil {
		c.subscriptionRepo = repository.NewSubscriptionRepo(
			c.Postgres().Pool,
		)
		c.log.Debug("subscription repo initialized")
	}
	return c.subscriptionRepo
}
