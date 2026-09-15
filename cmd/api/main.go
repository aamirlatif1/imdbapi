package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aamirlatif1/imdbapi/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	config config.Config
	logger *slog.Logger
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error(err.Error())
	}
	ctx := context.Background()

	pool, err := dbConnectionPool(ctx, cfg.Database.ConnectionString())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("database connection established")
	defer pool.Close()

	app := &application{
		config: cfg,
		logger: logger,
	}

	err = app.serve(pool)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}

func (app *application) serve(pool *pgxpool.Pool) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", app.config.HTTPPort),
		Handler:      app.routes(pool),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 10 * time.Minute,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Info(fmt.Sprintf("received signal %s", s.String()))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.Env)

	err := srv.ListenAndServe()
	app.logger.Error(err.Error())

	err = <-shutdown
	if err != nil {
		return err
	}
	app.logger.Info("stopped server")
	return nil
}

func dbConnectionPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}
	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = time.Hour
	poolCfg.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
