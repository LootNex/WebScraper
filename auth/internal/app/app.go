package app

import (
	"log/slog"
	"time"

	authapp "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/internal/app/grpc"
	auth "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/internal/services"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/internal/storage/psql"
)

type App struct {
	GRPCSrv     *authapp.App
	PsqlAuthSrv *psql.AuthStorage
}

type Storage struct {
	Host           string
	Password       string
	Port           int
	DBname         string
	User           string
	SSLMode        string
	ServiceName    string
	MigrationsPath string
}

func New(log *slog.Logger, grpcPort int, storageCredentials Storage, tokenTTL time.Duration) *App {
	db, err := psql.NewPostgresConnection(psql.ConnectionInfo{
		Host:     storageCredentials.Host,
		Port:     storageCredentials.Port,
		Password: storageCredentials.Password,
		DBName:   storageCredentials.DBname,
		SSLMode:  storageCredentials.SSLMode,
		Username: storageCredentials.User,
	})
	if err != nil {
		panic("no connection to postgres")
	}

	psqlClient := psql.New(db)

	authService := auth.New(log, psqlClient, psqlClient, tokenTTL)

	authApp := authapp.New(log, grpcPort, authService)

	return &App{
		GRPCSrv:     authApp,
		PsqlAuthSrv: psqlClient,
	}
}
