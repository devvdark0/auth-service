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

func main() {
	var cfgPath = flag.String("config_path", "./config/config.yaml", "set path to config file")

	cfg, err := config.MustLoad(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	api := api.New(cfg)

	log.Print("starting server...", "port=", cfg.App.Port)
	go func() {

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

	api.Stop(ctx)

	log.Println("server stopped gracefully")
}
