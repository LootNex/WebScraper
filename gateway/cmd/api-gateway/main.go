package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/gateway/proto"
	"google.golang.org/grpc"
)

type GatewayServer struct {
	authClient  proto.AuthServiceClient
	priceClient proto.PriceMonitoringServiceClient
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	authAddr := os.Getenv("AUTH_SERVICE_ADDR")
	priceAddr := os.Getenv("PRICE_SERVICE_ADDR")

	authConn, err := grpc.Dial(authAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect to auth service:", err)
	}
	defer authConn.Close()

	priceConn, err := grpc.Dial(priceAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect to price service:", err)
	}
	defer priceConn.Close()

	server := &GatewayServer{
		authClient:  proto.NewAuthServiceClient(authConn),
		priceClient: proto.NewPriceMonitoringServiceClient(priceConn),
	}

	r := mux.NewRouter()
	r.HandleFunc("/login", server.handleLogin).Methods("POST")
	r.HandleFunc("/additem", server.handleAddItem).Methods("POST")
	r.HandleFunc("/checkitem", server.handleCheckItem).Methods("GET")
	r.HandleFunc("/getallitems", server.handleGetAllItems).Methods("GET")

	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func (s *GatewayServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	resp, err := s.authClient.AuthenticateUser(context.Background(), &proto.AuthRequest{
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
	})
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *GatewayServer) handleAddItem(w http.ResponseWriter, r *http.Request) {
	userID, err := s.validateToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var reqBody struct {
		Link string `json:"link"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := s.priceClient.AddItem(context.Background(), &proto.AddItemRequest{
		UserId: userID,
		Link:   reqBody.Link,
	})
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *GatewayServer) handleCheckItem(w http.ResponseWriter, r *http.Request) {
	userID, err := s.validateToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	link := r.URL.Query().Get("link")
	if link == "" {
		http.Error(w, "Link is required", http.StatusBadRequest)
		return
	}

	resp, err := s.priceClient.CheckItem(context.Background(), &proto.CheckItemRequest{
		UserId: userID,
		Link:   link,
	})
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
	userID, err := s.validateToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp, err := s.priceClient.GetAllItems(context.Background(), &proto.GetAllItemsRequest{
		UserId: userID,
	})
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *GatewayServer) validateToken(r *http.Request) (int64, error) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if token == "" {
		return 0, fmt.Errorf("missing token")
	}

	resp, err := s.authClient.ValidateToken(context.Background(), &proto.ValidateTokenRequest{Token: token})
	if err != nil || !resp.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	return resp.UserId, nil
}