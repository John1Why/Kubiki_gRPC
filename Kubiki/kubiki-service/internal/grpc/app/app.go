package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	kubikigrpc "kubiki_Service/internal/grpc/kubiki"
	"kubiki_Service/internal/services/kubiki"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, servis *kubiki.Kubiki, port int) *App {
	gRPCServer := grpc.NewServer()

	kubikigrpc.Register(gRPCServer, servis)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server zapushchen", slog.Int("port", a.port))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("ostanovka grpc servera")

	a.gRPCServer.GracefulStop()
}