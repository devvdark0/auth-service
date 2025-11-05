package grpcapp

import (
	"fmt"
	authgrpc "github.com/devvdark0/auth-service/internal/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	log        *zap.Logger
	gRPCServer *grpc.Server
	gRPCPort   int
}

func New(log *zap.Logger, port int) *App {

	gRPCSrv := grpc.NewServer()
	authgrpc.Register(gRPCSrv)

	return &App{
		log:        log,
		gRPCServer: gRPCSrv,
		gRPCPort:   port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"
	log := a.log.With(zap.Field{Key: "op", String: op})

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.gRPCPort))
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	log.Info("running grpc server is running: ", zap.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(zap.Field{Key: "op", String: op}).Info("stopping grpc server: ", zap.Int("port", a.gRPCPort))

	a.gRPCServer.GracefulStop()
}
