package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	authpb "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/gateway/pkg/pb/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api authpb.Auth_V1Client
	log *slog.Logger
}

func New(ctx context.Context, log *slog.Logger, addr string, timeout time.Duration, retriesCount int) (*Client, error) {
	const op = "grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	grpcClient := authpb.NewAuth_V1Client(cc)

	return &Client{
		api: grpcClient,
		log: log,
	}, nil
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (c *Client) Register(ctx context.Context, telegramLogin, login, password string) (string, error) {
	const op = "grpc.Register"

	resp, err := c.api.Register(ctx, &authpb.RegisterRequest{
		TelegramLogin: telegramLogin,
		Login:         login,
		Password:      password,
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetUserId(), nil
}

func (c *Client) Login(ctx context.Context, telegramLogin, login, password string) (string, error) {
	const op = "grpc.Login"

	resp, err := c.api.Login(ctx, &authpb.LoginRequest{
		TelegramLogin: telegramLogin,
		Login:         login,
		Password:      password,
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetToken(), nil
}

func (c *Client) IsLogged(ctx context.Context, telegramLogin string) (bool, error) {
	const op = "grpc.IsLogged"

	resp, err := c.api.IsLogged(ctx, &authpb.IsLoggedRequest{
		TelegramLogin: telegramLogin,
	})

	if err != nil {
		return resp.GetIsLogged(), fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetIsLogged(), nil
}

func (c *Client) Logout(ctx context.Context, telegramLogin string) error {
	const op = "grpc.Login"

	_, err := c.api.Logout(ctx, &authpb.LogoutRequest{
		TelegramLogin: telegramLogin,
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
