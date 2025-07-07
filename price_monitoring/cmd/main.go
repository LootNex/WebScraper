package main

import (
	"log"
	"net"

	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price_monitoring/config"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price_monitoring/internal/db"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price_monitoring/internal/handlers"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price_monitoring/internal/service"
	proto "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price_monitoring/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	config, err := config.InitConfig()
	if err != nil {
		log.Fatalf("cannot init config err:%v", err)
	}

	PG_conn, err := postgres_db.InitDB(config.Postgres.DB_HOST, config.Postgres.DB_PORT, config.Postgres.DB_USER, config.Postgres.DB_PASSWORD, config.Postgres.DB_NAME)
	if err != nil {
		log.Fatalf("cannot connect to Postgres %v", err)
	}

	service := service.NewService(PG_conn)

	handler := handlers.NewHandler(service)

	GrpcServer := grpc.NewServer()
	reflection.Register(GrpcServer)

	proto.RegisterScraperServer(GrpcServer, handler)

	lis, err := net.Listen("tcp", ":"+config.Server.Port)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	log.Println("WebScraper is running")

	err = GrpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
