package migrator

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationsConnectionInfo struct {
	Host           string
	Port           int
	Username       string
	DBName         string
	SSLMode        string
	Password       string
	ServiceName    string
	MigrationsPath string
}

func Migrate(info MigrationsConnectionInfo) {
	URL := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=%v&search_path=%s&x-migrations-table=%s_migrations",
		info.Username,
		info.Password,
		info.Host,
		info.Port,
		info.DBName,
		info.SSLMode,
		info.ServiceName,
		info.ServiceName)

	m, err := migrate.New(info.MigrationsPath, URL)
	if err != nil {
		panic(err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}
}
