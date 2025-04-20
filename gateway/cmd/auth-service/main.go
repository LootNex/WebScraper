package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/gateway/internal/db"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/gateway/proto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type AuthServer struct {
	db *sql.DB
	proto.UnimplementedAuthServiceServer
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db := db.NewDB()
	defer db.Close()

	// Инициализация таблиц
	if err := initAuthTables(db); err != nil {
		log.Fatal("Failed to initialize database tables:", err)
	}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, &AuthServer{db: db})

	log.Println("Auth service running on :50051")
	if err := server.Serve(listener); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}

func initAuthTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login VARCHAR(50) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func (s *AuthServer) RegisterUser(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec("INSERT INTO users (login, password_hash) VALUES ($1, $2)", req.Login, string(hashedPassword))
	if err != nil {
		return &proto.RegisterResponse{
			Success: false,
			Message: "User already exists",
		}, nil
	}

	return &proto.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
	}, nil
}

func (s *AuthServer) AuthenticateUser(ctx context.Context, req *proto.AuthRequest) (*proto.AuthResponse, error) {
	var userID int64
	var hashedPassword string
	err := s.db.QueryRow("SELECT id, password_hash FROM users WHERE login = $1", req.Login).
		Scan(&userID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return &proto.AuthResponse{Success: false, Message: "User not found"}, nil
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return &proto.AuthResponse{Success: false, Message: "Invalid password"}, nil
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	return &proto.AuthResponse{
		Success: true,
		Token:   tokenString,
		Message: "Authenticated successfully",
	}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET not set")
	}

	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return &proto.ValidateTokenResponse{
			Valid:   false,
			Message: "Invalid token: " + err.Error(),
		}, nil
	}

	if !token.Valid {
		return &proto.ValidateTokenResponse{
			Valid:   false,
			Message: "Invalid token",
		}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &proto.ValidateTokenResponse{
			Valid:   false,
			Message: "Invalid token claims",
		}, nil
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return &proto.ValidateTokenResponse{
			Valid:   false,
			Message: "Invalid user ID in token",
		}, nil
	}

	return &proto.ValidateTokenResponse{
		Valid:  true,
		UserId: int64(userID),
	}, nil
}