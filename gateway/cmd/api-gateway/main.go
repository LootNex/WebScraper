package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	authclient "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/gateway/internal/clients/auth/grpc"
	trackerclient "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/gateway/internal/clients/tracker/grpc"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/gateway/internal/lib/jwt"

	"github.com/gorilla/mux"
)

type GatewayServer struct {
	authClient    authclient.Client
	trackerClient trackerclient.Client
}

func main() {
	authAddr := os.Getenv("AUTH_SERVICE_ADDR")
	trackerAddr := os.Getenv("PRICE_SERVICE_ADDR")

	var logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	authClient, err := authclient.New(logger, authAddr, time.Second*5, 1)
	if err != nil {
		panic("failed to connect to auth server")
	}

	trackerClient, err := trackerclient.New(logger, trackerAddr, time.Second*5, 1)
	if err != nil {
		panic("failed to connect to tracker server")
	}

	server := &GatewayServer{
		authClient:    *authClient,
		trackerClient: *trackerClient,
	}

	r := mux.NewRouter()
	r.HandleFunc("/login", server.handleLogin).Methods("POST")
	r.HandleFunc("/register", server.handleRegister).Methods("POST")
	r.HandleFunc("/logout", server.handleLogout).Methods("POST")
	r.HandleFunc("/check_item", server.handleGetItem).Methods("POST")
	r.HandleFunc("/get_all_items", server.handleGetAllItems).Methods("GET")

	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func (s *GatewayServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login         string `json:"login"`
		Password      string `json:"password"`
		TelegramLogin string `json:"telegram_login"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	_, err := s.authClient.IsLogged(context.Background(), req.TelegramLogin)
	if err == nil {
		http.Error(w, fmt.Sprintf("registration failed: %v", err), http.StatusInternalServerError)
		return
	}

	resp, err := s.authClient.Register(context.Background(), req.TelegramLogin, req.Login, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("registraton failed: %v", err), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"user_id": resp,
	})
}

func (s *GatewayServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login         string `json:"login"`
		Password      string `json:"password"`
		TelegramLogin string `json:"telegram_login"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	_, err := s.authClient.IsLogged(context.Background(), req.TelegramLogin)
	if err == nil {
		http.Error(w, fmt.Sprintf("login failed: %v", err), http.StatusInternalServerError)
		return
	}

	resp, err := s.authClient.Login(context.Background(), req.TelegramLogin, req.Login, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("login failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": resp,
	})
}

func (s *GatewayServer) handleLogout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TelegramLogin string `json:"telegram_login"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	_, err := s.authClient.IsLogged(context.Background(), req.TelegramLogin)
	if err != nil {
		http.Error(w, fmt.Sprintf("logout failed: %v", err), http.StatusInternalServerError)
		return
	}

	err = s.authClient.Logout(context.Background(), req.TelegramLogin)
	if err != nil {
		http.Error(w, fmt.Sprintf("logout failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func (s *GatewayServer) handleGetItem(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TelegramLogin string `json:"telegram_login"`
		Link          string `json:"link"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, err := s.authClient.IsLogged(context.Background(), req.TelegramLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := jwt.GetUserID(token)
	if err != nil {
		log.Printf("JWT validation failed: %v", err)
		http.Error(w, "failed to validate user's token", http.StatusUnauthorized)
		return
	}

	resp, err := s.trackerClient.GetItem(context.Background(), userID, req.Link)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *GatewayServer) handleGetAllItems(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TelegramLogin string `json:"telegram_login"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, err := s.authClient.IsLogged(context.Background(), req.TelegramLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := jwt.GetUserID(token)
	if err != nil {
		http.Error(w, "failed to validate user's token", http.StatusUnauthorized)
		return
	}

	resp, err := s.trackerClient.GetAllItems(context.Background(), userID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
