package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/devvdark0/auth-service/internal/config"
)

var configPath = flag.String("config-path", "./config/config.yaml", "use to set config path")

func main() {

	cfg, err := config.MustLoad(*configPath)
	if err != nil {
		panic(err)
	}

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port),
		WriteTimeout: cfg.App.Timeout,
		ReadTimeout:  cfg.App.Timeout,
		IdleTimeout:  cfg.App.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

}
