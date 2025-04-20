package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/internal/domain/models"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/internal/lib/jwt"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userChanger  UserChanger
	userProvider UserProvider
	TokenTTL     time.Duration
}

type UserChanger interface {
	SaveUser(ctx context.Context, user models.User) error
}

type UserProvider interface {
	UserLoginsByTelegram(ctx context.Context, telegramLogin string) ([]models.User, error)
}

var (
	ErrUserExists         = errors.New("user already exist")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func New(log *slog.Logger, userChanger UserChanger, userProvider UserProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userChanger:  userChanger,
		userProvider: userProvider,
		TokenTTL:     tokenTTL,
	}
}

func (a *Auth) Register(ctx context.Context, telegramLogin, login, password string) (string, error) {
	const op = "auth.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("login", login),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	id := uuid.New().String()

	user := models.User{
		ID:            id,
		Login:         login,
		TelegramLogin: telegramLogin,
		PassHash:      passHash,
	}

	if err := a.userChanger.SaveUser(ctx, user); err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to save user", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	return id, nil
}

func (a *Auth) Login(ctx context.Context, telegramLogin, login, password string) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("login", login),
	)

	log.Info("attemting to log user in")

	users, err := a.userProvider.UserLoginsByTelegram(ctx, telegramLogin)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		log.Error("failed to get user", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	foundTelegram := false
	var requiredUser models.User
	for _, user := range users {
		if user.Login == login {
			foundTelegram = true
			requiredUser = user
			break
		}
	}

	if !foundTelegram {
		log.Warn("user with provided telegram login and user login not found")

		return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
	}

	if err := bcrypt.CompareHashAndPassword(requiredUser.PassHash, []byte(password)); err != nil {
		log.Info("Invalid credentials", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(requiredUser, a.TokenTTL)
	if err != nil {
		log.Error("failed to generate token", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
