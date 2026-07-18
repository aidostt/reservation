package app

import (
	"dip/internal/config"
	"dip/internal/delivery"
	"fmt"
	"log"
	"net"

	"dip/internal/repository"
	"dip/internal/server"
	"dip/internal/service"

	"context"
	"dip/pkg/logger"
	"dip/pkg/tracing"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	// "github.com/go-git/go-git/v5/plumbing/transport/server"
)

func Run(configPath, envPath string) {
	cfg, err := config.Init(configPath, envPath)
	if err != nil {
		logger.Error(err)

		return
	}

	shutdownTracing, err := tracing.Init(context.Background(), "reservation-service")
	if err != nil {
		logger.Errorf("tracing init: %s", err.Error())
	} else {
		defer func() { _ = shutdownTracing(context.Background()) }()
	}
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName))
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}

	log.Println("Connected to db")
	defer pool.Close()

	repos := repository.NewRepository(pool)
	services := service.NewService(service.Dependencies{
		Environment:  cfg.Environment,
		Domain:       cfg.GRPC.Host,
		Repos:        repos,
		TurnDuration: cfg.TurnDuration,
	})
	handlers := delivery.NewHandler(services)

	// gRPC Server
	srv := server.NewServer()
	srv.RegisterServers(handlers)

	l, err := net.Listen("tcp", fmt.Sprintf("%v:%v", cfg.GRPC.Host, cfg.GRPC.Port))
	if err != nil {
		logger.Errorf("error occurred while getting listener for the server: %s\n", err.Error())
		return
	}
	go func() {
		if err = srv.Run(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running grpc server: %s\n", err.Error())
		}
	}()

	metricsSrv := &http.Server{Addr: ":" + metricsPort(), Handler: server.MetricsHandler()}
	go func() {
		if err := metricsSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("metrics server error: %s\n", err.Error())
		}
	}()

	logger.Info("Server started at: " + cfg.GRPC.Host + ":" + cfg.GRPC.Port)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	srv.Stop()
	_ = metricsSrv.Shutdown(context.Background())
	logger.Info("Stopping server at: " + cfg.GRPC.Host + ":" + cfg.GRPC.Port)

}

// metricsPort is the port for the Prometheus metrics endpoint; it defaults to
// 9464 and can be overridden with METRICS_PORT.
func metricsPort() string {
	if p := os.Getenv("METRICS_PORT"); p != "" {
		return p
	}
	return "9464"
}
