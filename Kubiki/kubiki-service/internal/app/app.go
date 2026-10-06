package app

import (
	"log/slog"

	"kubiki_Service/internal/config"
	grpcapp "kubiki_Service/internal/grpc/app"
	"kubiki_Service/internal/services/kubiki"
	"kubiki_Service/internal/storage/postgres"
)

type App struct {
	GRPCServer    *grpcapp.App
	Khranilishche *postgres.Khranilishche
}

func New(log *slog.Logger, conf *config.Config) *App {
	khranilishche, err := postgres.NovoeKhranilishche(conf.Storage.DSN)
	if err != nil {
		panic(err)
	}

	kubikiServis := kubiki.New(log, khranilishche)

	grpcApp := grpcapp.New(log, kubikiServis, conf.GRPC.Port)

	return &App{
		GRPCServer:    grpcApp,
		Khranilishche: khranilishche,
	}
}