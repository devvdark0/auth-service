package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/devvdark0/auth-service/internal/config"
	"github.com/devvdark0/auth-service/internal/logger"
)

var configPath = flag.String("config-path", "./config/config.yaml", "use to set config path")

func main() {

	cfg, err := config.MustLoad(*configPath)
	if err != nil {
		panic(err)
	}

	log := logger.NewLogger(cfg.App.Env)

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port),
		WriteTimeout: cfg.App.Timeout,
		ReadTimeout:  cfg.App.Timeout,
		IdleTimeout:  cfg.App.IdleTimeout,
	}

	log.Info("starting server...", "port", cfg.App.Port)

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to run the server", "err", err)
		panic(err)
	}

}
