package app

import (
	"dip/handler"
	"dip/internal/config"
	"fmt"
	"log"

	"dip/internal/server"
	"dip/repository"
	"dip/service"

	"context"
	"dip/internal/logger"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	// "github.com/go-git/go-git/v5/plumbing/transport/server"
)

func Run(configPath, envPath string) {
	cfg, err := config.Init(configPath, envPath)
	if err != nil {
		logger.Error(err)

		return
	}

	pool, err := pgxpool.New(context.Background(), cfg.Postgres.URI)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}

	fmt.Println(cfg.Postgres.URI)
	log.Println("Connected to db")
	defer pool.Close()

	repos := repository.NewRepository(pool)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init())

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}

}
