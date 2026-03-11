package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/devvdark0/auth-service/internal/config"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	var cfgPath = flag.String("config_path", "./config/config.yaml", "set path to config file")

	cfg, err := config.MustLoad(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	log := configureLogger(cfg.App.Env)

	log.Debug("initializing db....")
	db, err := config.InitDb(&cfg.Db)
	if err != nil {
		log.Error("failed to init databaase", "err", err)
	}
	_ = db

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port),
		Handler:      nil,
		ReadTimeout:  cfg.App.Timeout,
		WriteTimeout: cfg.App.Timeout,
		IdleTimeout:  cfg.App.IdleTimeout,
	}

	log.Info("starting server...", "port", cfg.App.Port)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func configureLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return log
}
