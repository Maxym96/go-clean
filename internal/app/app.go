// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"go-clean/config"
	amqprpc "go-clean/internal/controller/amqp_rpc"
	v1 "go-clean/internal/controller/http/v1"
	"go-clean/internal/usecase"
	"go-clean/internal/usecase/repo"
	"go-clean/internal/usecase/webapi"
	"go-clean/pkg/httpserver"
	"go-clean/pkg/logger"
	"go-clean/pkg/postgres"
	"go-clean/pkg/rabbitmq/rmq_rpc/server"

	"github.com/gin-gonic/gin"
)

// Run creates objects via constructors.
func Run(cfg *config.Configuration) {
	l := logger.New(cfg.LogLevel)

	// Repository
	pgPoolMax, err := strconv.Atoi(cfg.PgPoolMax)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - pgPoolMax.Atoi: %w", err))
		return
	}
	pg, err := postgres.New(cfg.PostgreSQLUrl, postgres.MaxPoolSize(pgPoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	//Postgres migration
	err = postgres.Migrate(cfg)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - unable to apply migrations: %q\n", err))
		return
	}

	// Use case
	translationUseCase := usecase.NewTranslationUseCase(
		repo.New(pg),
		webapi.New(),
	)

	// RabbitMQ RPC Server
	rmqRouter := amqprpc.NewRouter(translationUseCase)

	rmqServer, err := server.New(cfg.RabbitMQUrl, cfg.RmqRpcServer, rmqRouter, l)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	}

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, translationUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HttpPort))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-rmqServer.Notify():
		l.Error(fmt.Errorf("app - Run - rmqServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	err = rmqServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	}
}
