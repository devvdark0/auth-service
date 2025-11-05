package app

import (
	grpcapp "github.com/devvdark0/auth-service/internal/app/grpc"
	"go.uber.org/zap"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *zap.Logger, storagePath string, grpcPort int) *App {

	//TODO: init storage

	//TODO: init auth service

	gRPCApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCServer: gRPCApp,
	}
}
