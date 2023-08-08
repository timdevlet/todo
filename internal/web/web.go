package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	validator "github.com/go-playground/validator/v10"
	"github.com/timdevlet/todo/internal/configs"
	logs "github.com/timdevlet/todo/internal/log"
	"github.com/timdevlet/todo/internal/todo"
	"github.com/timdevlet/todo/pkg/postgres"
)

type Web struct {
	Port   int
	Router *echo.Echo

	Options *configs.Configs

	TodoService todo.ITodo
	LogService  logs.ILogService
}

func NewWeb(conf *configs.Configs) *Web {
	w := &Web{
		Port:    conf.PORT,
		Options: conf,
	}

	// Todo service
	pdb, _ := postgres.NewPDB(postgres.Config{
		Host:     conf.DB_HOST,
		Port:     conf.DB_PORT,
		User:     conf.DB_USER,
		Password: conf.DB_PASSWORD,
		DBName:   conf.DB_NAME,
		SSLMode:  conf.DB_SSL,
	})
	todoRepository := todo.NewTodoRepository(pdb)
	todoService := todo.NewTodoService(todoRepository)

	// Log service

	logsRepository := logs.NewLogRepository(pdb)
	logService := logs.NewLogService(logsRepository)

	//

	w.TodoService = todoService
	w.LogService = logService

	return w
}

func (a *Web) Init() *echo.Echo {
	e := echo.New()

	//

	e.Use(SetAdminUserMiddleware)
	e.Use(middleware.BodyLimit("1M"))
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.BodyDump(LogMiddleware(a)))

	//

	e.Validator = &CustomValidator{validator: validator.New()}

	//

	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 15, ExpiresIn: 1 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}

	e.Use(middleware.RateLimiterWithConfig(config))

	//

	if a.Options.LOG_FORMAT != "plain" {
		e.Use(middleware.Logger())
	} else {
		e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:      true,
			LogStatus:   true,
			LogRemoteIP: true,
			LogError:    true,
			LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
				log.WithFields(log.Fields{
					"URI":     values.URI,
					"status":  values.Status,
					"latency": values.Latency,
					"ip":      values.RemoteIP,
				}).Info(values.Error)

				return nil
			},
		}))
	}

	// Routers

	initTodosRoutes(a, e)
	initHealthRoutes(a, e)
	initMetricsRoutes(a, e)

	//

	a.Router = e

	return e
}

func (a *Web) Run() {
	// Start server
	go func() {
		if err := a.Router.Start(fmt.Sprintf(":%d", a.Port)); err != nil && err != http.ErrServerClosed {
			a.Router.Logger.Fatal("🙏 shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.Router.Shutdown(ctx); err != nil {
		a.Router.Logger.Fatal(err)
	}
}
