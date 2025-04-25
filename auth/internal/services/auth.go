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
	tokenChanger TokenChanger
	tokenChecker TokenChecker
	TokenTTL     time.Duration
}

type UserChanger interface {
	SaveUser(ctx context.Context, user models.User) error
}

type UserProvider interface {
	UserLoginsByTelegram(ctx context.Context, telegramLogin string) ([]models.User, error)
}

type TokenChanger interface {
	SaveJWT(ctx context.Context, JWT, telegramLogin string, ttl time.Duration) error
	DeleteJWT(ctx context.Context, telegramLogin string) error
}

type TokenChecker interface {
	JWT(ctx context.Context, telegramLogin string) (string, error)
}

var (
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenNotFound      = errors.New("token not found")
	ErrTokenExists        = errors.New("token already exists")
)

func New(log *slog.Logger, userChanger UserChanger, userProvider UserProvider, tokenChanger TokenChanger,
	tokenChecker TokenChecker, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userChanger:  userChanger,
		userProvider: userProvider,
		tokenChanger: tokenChanger,
		tokenChecker: tokenChecker,
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

	if err := a.tokenChanger.SaveJWT(ctx, token, telegramLogin, a.TokenTTL); err != nil {
		if errors.Is(err, storage.ErrTokenExists) {
			log.Warn("user already logged in", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %w", op, ErrTokenExists)
		}
		log.Error("failed to log user in", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully logged the user in")

	return token, nil
}

func (a *Auth) IsLogged(ctx context.Context, telegramLogin string) (string, error) {
	const op = "auth.IsLogged"

	log := a.log.With(
		slog.String("op", op),
		slog.String("telegram login", telegramLogin),
	)

	log.Info("checking if this user is logged in")

	token, err := a.tokenChecker.JWT(ctx, telegramLogin)
	if err != nil {
		if errors.Is(err, storage.ErrTokenNotFound) {
			log.Warn("token not found", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %w", op, ErrTokenNotFound)
		}
		log.Error("failed to get token", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully got the info about the user status")

	return token, nil
}

func (a *Auth) Logout(ctx context.Context, telegramLogin string) error {
	const op = "auth.Logout"

	log := a.log.With(
		slog.String("op", op),
		slog.String("telegram login", telegramLogin),
	)

	log.Info("loggin the user out")

	err := a.tokenChanger.DeleteJWT(ctx, telegramLogin)
	if err != nil {
		log.Error("failed to log the user out", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully logged the user out")

	return nil
}
