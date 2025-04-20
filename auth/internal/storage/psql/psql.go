package psql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/internal/domain/models"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/internal/storage"
)

type ConnectionInfo struct {
	Host     string
	Port     int
	Username string
	DBName   string
	SSLMode  string
	Password string
}

type AuthStorage struct {
	db *sql.DB
}

func NewPostgresConnection(info ConnectionInfo) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=%v",
		info.Username,
		info.Password,
		info.Host,
		info.Port,
		info.DBName,
		info.SSLMode))

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func New(db *sql.DB) *AuthStorage {
	return &AuthStorage{
		db: db,
	}
}

func (a *AuthStorage) Close() {
	a.db.Close()
}

func (a *AuthStorage) SaveUser(ctx context.Context, user models.User) error {
	const op = "storage.psql.SaveUser"

	if err := a.findLogins(ctx, user.TelegramLogin, user.Login); err == nil {
		return fmt.Errorf("%s: %w", op, storage.ErrUserExists)
	}

	_, err := a.db.ExecContext(ctx,
		"INSERT INTO auth.users (id, telegram_login, login, pass_hash) VALUES ($1, $2, $3, $4)",
		user.ID, user.TelegramLogin, user.Login, user.PassHash)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *AuthStorage) UserLoginsByTelegram(ctx context.Context, telegramLogin string) ([]models.User, error) {
	const op = "storage.psql.UserLoginsByTelegram"

	if err := a.findTelegramLogin(ctx, telegramLogin); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	query := `SELECT id, login, telegram_login, pass_hash FROM auth.users WHERE telegram_login = $1`

	rows, err := a.db.QueryContext(ctx, query, telegramLogin)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.TelegramLogin,
			&user.PassHash,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (a *AuthStorage) findLogins(ctx context.Context, telegramLogin, login string) error {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM auth.users WHERE telegram_login = $1 AND login = $2)",
		telegramLogin, login).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return storage.ErrUserNotFound
	}

	return nil
}

func (a *AuthStorage) findTelegramLogin(ctx context.Context, telegramLogin string) error {
	var exists bool
	err := a.db.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM auth.users WHERE telegram_login = $1)", telegramLogin).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return storage.ErrUserNotFound
	}
	return nil
}
