package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/devvdark0/auth-service/internal/config"
)

var cfgPath = flag.String("config_path", "config/config.yaml", "define config path")

func main() {
	flag.Parse()

	cfg, err := config.MustLoad(*cfgPath)
	if err != nil {
		panic(err)
	}

	//TODO: init logger

	//TODO: init storage

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port),
		ReadTimeout:  cfg.App.Timeout,
		WriteTimeout: cfg.App.Timeout,
		IdleTimeout:  cfg.App.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

}
