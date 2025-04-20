package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/gateway/internal/db"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/gateway/internal/scrapping"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/gateway/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type PriceMonitoringServer struct {
	db *sql.DB
	proto.UnimplementedPriceMonitoringServiceServer
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db := db.NewDB()
	defer db.Close()

	if err := createTables(db); err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	server := &PriceMonitoringServer{db: db}
	grpcServer := grpc.NewServer()
	proto.RegisterPriceMonitoringServiceServer(grpcServer, server)

	go server.startPriceChecker()

	log.Println("Price monitoring service running on :50052")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id SERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			link TEXT NOT NULL,
			current_price DECIMAL(10, 2),
			last_checked TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, link)
		);
		CREATE TABLE IF NOT EXISTS price_history (
			id SERIAL PRIMARY KEY,
			item_id INTEGER REFERENCES items(id) ON DELETE CASCADE,
			price DECIMAL(10, 2) NOT NULL,
			checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func (s *PriceMonitoringServer) AddItem(ctx context.Context, req *proto.AddItemRequest) (*proto.AddItemResponse, error) {
	_, err := s.db.Exec(
		"INSERT INTO items (user_id, link) VALUES ($1, $2) ON CONFLICT (user_id, link) DO NOTHING",
		req.UserId, req.Link,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add item: %w", err)
	}

	return &proto.AddItemResponse{
		Success: true,
		Message: "Item added successfully",
	}, nil
}

func (s *PriceMonitoringServer) CheckItem(ctx context.Context, req *proto.CheckItemRequest) (*proto.CheckItemResponse, error) {
	var itemID int
	var currentPrice sql.NullFloat64
	var lastChecked time.Time
	var link string

	err := s.db.QueryRow(
		"SELECT id, link, current_price, last_checked FROM items WHERE user_id = $1 AND link = $2",
		req.UserId, req.Link,
	).Scan(&itemID, &link, &currentPrice, &lastChecked)
	if err != nil {
		if err == sql.ErrNoRows {
			return &proto.CheckItemResponse{Success: false, Message: "Item not found"}, nil
		}
		return nil, fmt.Errorf("failed to query item: %w", err)
	}

	price := 0.0
	if currentPrice.Valid {
		price = currentPrice.Float64
	}

	return &proto.CheckItemResponse{
		Success:      true,
		Link:         link,
		CurrentPrice: price,
		LastChecked:  lastChecked.Format(time.RFC3339),
	}, nil
}

func (s *PriceMonitoringServer) GetAllItems(ctx context.Context, req *proto.GetAllItemsRequest) (*proto.GetAllItemsResponse, error) {
	rows, err := s.db.Query(
		"SELECT link, current_price, last_checked FROM items WHERE user_id = $1 ORDER BY created_at DESC",
		req.UserId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	var items []*proto.Item
	for rows.Next() {
		var link string
		var currentPrice sql.NullFloat64
		var lastChecked time.Time
		if err := rows.Scan(&link, &currentPrice, &lastChecked); err != nil {
			log.Printf("Error scanning item: %v", err)
			continue
		}

		price := 0.0
		if currentPrice.Valid {
			price = currentPrice.Float64
		}

		items = append(items, &proto.Item{
			Link:         link,
			CurrentPrice: price,
			LastChecked:  lastChecked.Format(time.RFC3339),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return &proto.GetAllItemsResponse{
		Success: true,
		Items:   items,
	}, nil
}

func (s *PriceMonitoringServer) GetPriceHistory(ctx context.Context, req *proto.PriceHistoryRequest) (*proto.PriceHistoryResponse, error) {
	var itemID int
	err := s.db.QueryRow(
		"SELECT id FROM items WHERE user_id = $1 AND link = $2",
		req.UserId, req.Link,
	).Scan(&itemID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &proto.PriceHistoryResponse{Success: false, Message: "Item not found"}, nil
		}
		return nil, fmt.Errorf("failed to get item ID: %w", err)
	}

	rows, err := s.db.Query(
		"SELECT price, checked_at FROM price_history WHERE item_id = $1 ORDER BY checked_at DESC LIMIT 30",
		itemID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query price history: %w", err)
	}
	defer rows.Close()

	var history []*proto.PriceEntry
	for rows.Next() {
		var price float64
		var checkedAt time.Time
		if err := rows.Scan(&price, &checkedAt); err != nil {
			log.Printf("Error scanning price entry: %v", err)
			continue
		}
		history = append(history, &proto.PriceEntry{
			Price:     price,
			CheckedAt: checkedAt.Format(time.RFC3339),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return &proto.PriceHistoryResponse{
		Success: true,
		History: history,
	}, nil
}

func (s *PriceMonitoringServer) CheckAll(ctx context.Context, req *proto.CheckAllRequest) (*proto.CheckAllResponse, error) {
	rows, err := s.db.Query(
		"SELECT id, user_id, link, current_price FROM items WHERE last_checked < NOW() - INTERVAL '1 hour' OR last_checked IS NULL",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query items for checking: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, userID int64
		var link string
		var currentPrice sql.NullFloat64
		if err := rows.Scan(&id, &userID, &link, &currentPrice); err != nil {
			log.Printf("Error scanning item: %v", err)
			continue
		}

		newPrice := scraping.GetCurrentPrice(link)
	
		tx, err := s.db.Begin()
		if err != nil {
			log.Printf("Error starting transaction: %v", err)
			continue
		}

		if _, err := tx.Exec(
			"INSERT INTO price_history (item_id, price) VALUES ($1, $2)",
			id, newPrice,
		); err != nil {
			tx.Rollback()
			log.Printf("Error saving price history: %v", err)
			continue
		}

		if _, err := tx.Exec(
			"UPDATE items SET current_price = $1, last_checked = NOW() WHERE id = $2",
			newPrice, id,
		); err != nil {
			tx.Rollback()
			log.Printf("Error updating price: %v", err)
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Error committing transaction: %v", err)
			continue
		}

		if currentPrice.Valid && newPrice < currentPrice.Float64 {
			go s.notifyPriceDrop(userID, link, newPrice, currentPrice.Float64)
		}
	}

	return &proto.CheckAllResponse{Success: true}, nil
}

func (s *PriceMonitoringServer) notifyPriceDrop(userID int64, link string, newPrice, oldPrice float64) {
	log.Printf("Price dropped for user %d, item %s: %.2f -> %.2f (%.2f%%)", 
		userID, link, oldPrice, newPrice, (oldPrice-newPrice)/oldPrice*100)
	// TODO: Реализовать отправку уведомления через API Gateway
}

func (s *PriceMonitoringServer) startPriceChecker() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		if _, err := s.CheckAll(context.Background(), &proto.CheckAllRequest{}); err != nil {
			log.Printf("Error checking prices: %v", err)
		}
	}
}