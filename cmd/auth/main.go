package main

import (
	"github.com/devvdark0/auth-service/internal/app"
	"github.com/devvdark0/auth-service/internal/config"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := configureLogger(cfg.Env)
	defer log.Sync()

	log.Info("starting application...")

	application := app.New(log, cfg.StoragePath, cfg.GRPC.Port)
	go application.GRPCServer.MustRun()

	//TODO: start grpc-server

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("stopping application", zap.String("signal", sign.String()))

	application.GRPCServer.Stop()

	log.Info("application stopped")
}

func configureLogger(env string) *zap.Logger {
	var (
		log *zap.Logger
		err error
	)
	switch env {
	case envLocal:
		log, err = zap.NewDevelopment()
	case envProd:
		log, err = zap.NewProduction()
	}
	if err != nil {
		panic("failed to set up logger: " + err.Error())
	}
	return log
}
