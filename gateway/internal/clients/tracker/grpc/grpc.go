package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	trackerpb "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/gateway/gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price-monitoring"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/gateway/internal/domain/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api trackerpb.ScraperClient
	log *slog.Logger
}

func New(log *slog.Logger, addr string, timeout time.Duration, retriesCount int) (*Client, error) {
	const op = "grpc.tracker.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	grpcClient := trackerpb.NewScraperClient(cc)

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

func (c *Client) GetItem(ctx context.Context, userID, link string) (models.Item, error) {
	const op = "grpc.tracker.GetItem"

	resp, err := c.api.GetItem(ctx, &trackerpb.GetItemRequest{
		Link:   link,
		UserId: userID,
	})

	if err != nil {
		return models.Item{}, fmt.Errorf("%s: %w", op, err)
	}

	item := models.Item{
		Name:            resp.GetItem().GetName(),
		StartPrice:      resp.GetItem().GetStartPrice(),
		CurrentPrice:    resp.GetItem().GetCurrentPrice(),
		DifferencePrice: resp.GetItem().GetDiffPrice(),
	}

	return item, nil
}

func (c *Client) GetAllItems(ctx context.Context, userID string) ([]*models.Item, error) {
	const op = "grpc.tracker.GetAllItems"

	resp, err := c.api.GetAllItems(ctx, &trackerpb.GetAllItemsRequest{
		UserId: userID,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	items := make([]*models.Item, len(resp.GetItems()))
	for i, item := range resp.GetItems() {
		items[i] = &models.Item{
			Name:            item.GetName(),
			StartPrice:      item.GetStartPrice(),
			CurrentPrice:    item.GetCurrentPrice(),
			DifferencePrice: item.GetDiffPrice(),
		}
	}

	return items, nil
}
