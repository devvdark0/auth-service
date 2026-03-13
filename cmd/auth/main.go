package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devvdark0/auth-service/internal/api"
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

	api := api.New(cfg)

	done := make(chan struct{})

	log.Print("starting server...", "port=", cfg.App.Port)
	go func() {
		defer close(done)

		if err := api.Run(); err != nil {
			log.Panic(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	<-sigCh

	log.Print("received signal, initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	<-done

	api.Stop(ctx)

	log.Println("server stopped gracefully")
}
