package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/devvdark0/auth-service/internal/config"
)

func main() {
	var cfgPath = flag.String("config_path", "./config/config.yaml", "set path to config file")

	cfg, err := config.MustLoad(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port),
		Handler:      nil,
		ReadTimeout:  cfg.App.Timeout,
		WriteTimeout: cfg.App.Timeout,
		IdleTimeout:  cfg.App.IdleTimeout,
	}

	log.Print("starting server....", "port=", cfg.App.Port)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
