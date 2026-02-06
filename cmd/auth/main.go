package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/devvdark0/auth-service/internal/config"
	"github.com/devvdark0/auth-service/internal/db"
	"github.com/devvdark0/auth-service/internal/migrations"
	"golang.org/x/sync/errgroup"
)

var cfgPath = flag.String("config_path", "config/config.yaml", "define config path")

func main() {
	flag.Parse()

	cfg, err := config.MustLoad(*cfgPath)
	if err != nil {
		panic(err)
	}

	log := configureLogger(cfg.App.Env)

	db, err := db.InitDB(cfg.DB)
	if err != nil {
		log.Error("failed to initialize database", "err", err)
		os.Exit(1)
	}

	if err := migrations.RunMigrations(cfg.DB.URL, "./migrations"); err != nil {
		log.Error("migrations are failed", "err", err)
	}
	log.Info("migrations are succesfully applied")
	_ = db

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port),
		ReadTimeout:  cfg.App.Timeout,
		WriteTimeout: cfg.App.Timeout,
		IdleTimeout:  cfg.App.IdleTimeout,
	}

	erg, _ := errgroup.WithContext(context.Background())

	log.Info("starting httpserver...")

	erg.Go(func() error {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to run server", "err", err)
			return err
		}
		log.Info("server was succesfully started...", "port", cfg.App.Port)
		return nil
	})

	erg.Go(func() error {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

		<-signalCh

		ctx, cancel := context.WithTimeout(context.Background(), cfg.App.ShutdownTimeout)
		defer cancel()

		log.Info("shutting down gracefully...")

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("failed to stop the server", "err", err)
			return err
		}
		return nil
	})

	if err := erg.Wait(); err != nil && err != context.Canceled {
		log.Error("server stopped with err", "err", err)
	} else {
		log.Info("server stopped")
	}

}

func configureLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "dev":
		log = slog.New(slog.NewTextHandler(os.Stdout, nil))
	case "prod":
		log = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	return log
}
