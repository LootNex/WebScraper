package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/internal/app"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/internal/config"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/internal/storage/psql/migrator"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger()

	log.Info("starting auth service")

	application := app.New(log, cfg.GRPC.Port, app.Storage{
		Host:     cfg.Storage.Host,
		Password: cfg.Storage.Password,
		Port:     cfg.Storage.Port,
		SSLMode:  cfg.Storage.SSLmode,
		User:     cfg.Storage.User,
		DBname:   cfg.Storage.DBname,
	}, app.TokensStorage{
		Addr:cfg.TokensStorage.Addr,
		Password: cfg.TokensStorage.Password,
	} ,cfg.TokenTTL)

	migrator.Migrate(migrator.MigrationsConnectionInfo{
		Host:           cfg.Storage.Host,
		Port:           cfg.Storage.Port,
		Password:       cfg.Storage.Password,
		Username:       cfg.Storage.User,
		DBName:         cfg.Storage.DBname,
		SSLMode:        cfg.Storage.SSLmode,
		ServiceName:    cfg.Storage.ServiceName,
		MigrationsPath: cfg.Storage.MigrationsPath,
	})

	go func() {
		application.GRPCSrv.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGALRM)

	sign := <-stop

	log.Info("stopping auth API", slog.String("signal", sign.String()))

	application.GRPCSrv.Stop()
	application.PsqlAuthSrv.Close()

	log.Info("auth API stopped")
}

func setupLogger() *slog.Logger {
	var log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	return log
}
