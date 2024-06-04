package app

import (
	"dip/handler"
	"dip/internal/config"
	"fmt"
	"log"
	"net"

	"dip/internal/server"
	"dip/repository"
	"dip/service"

	"context"
	"dip/internal/logger"
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
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName))
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}

	log.Println("Connected to db")
	defer pool.Close()

	repos := repository.NewRepository(pool)
	services := service.NewService(service.Dependencies{
		Environment: cfg.Environment,
		Domain:      cfg.GRPC.Host,
		Repos:       repos,
	})
	handlers := handler.NewHandler(services)

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

	logger.Info("Server started at: " + cfg.GRPC.Host + ":" + cfg.GRPC.Port)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	srv.Stop()
	logger.Info("Stopping server at: " + cfg.GRPC.Host + ":" + cfg.GRPC.Port)

}
