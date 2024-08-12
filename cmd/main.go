// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cs2-server/backend/config"
	_ "github.com/cs2-server/backend/docs"
	"github.com/cs2-server/backend/internal/api"
	"github.com/cs2-server/backend/internal/middleware"
	"github.com/cs2-server/backend/internal/service"
	"github.com/cs2-server/backend/internal/storage"
	"github.com/cs2-server/backend/pkg/jwt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	swagger "github.com/swaggo/http-swagger"
)

func main() {
	if err := run(); err != nil {
		logrus.Fatalln(err)
	}
}

func run() error {
	logger := logrus.New()

	cfg, err := config.Init()
	if err != nil {
		return fmt.Errorf("cfg: %v", err)
	}

	db, err := pgxpool.Connect(context.Background(), cfg.Postgres.DSN)
	if err != nil {
		return fmt.Errorf("db: %v", err)
	}

	jwt := jwt.New(cfg.JWT.Key)
	auth := api.NewAuthAPI(cfg, logger, jwt, service.NewAuthService(storage.NewAuthStorage(db)))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /swagger/*", swagger.Handler(swagger.URL(cfg.Swagger.URL)))

	mux.HandleFunc("GET /auth/login", middleware.Log(auth.Login))
	mux.HandleFunc("GET /auth/process", middleware.Log(auth.ProcessLogin))
	mux.HandleFunc("POST /auth/refresh", jwt.Auth(middleware.Log(auth.RefreshToken)))

	mux.HandleFunc("GET /profile", jwt.Auth(middleware.Log(auth.GetProfile)))

	var (
		sigCh = make(chan os.Signal, 1)
		errCh = make(chan error)
	)

	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	s := &http.Server{
		Addr:         cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler:      mux,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http start: %v", err)
		}
	}()

	logger.Infof("http server is listening on %s:%s", cfg.HTTP.Host, cfg.HTTP.Port)

	select {
	case sig := <-sigCh:
		logger.Infof("got signal: %s. shutting down...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			return fmt.Errorf("http shutdown: %v", err)
		}
	case err := <-errCh:
		return err
	}

	return nil
}
