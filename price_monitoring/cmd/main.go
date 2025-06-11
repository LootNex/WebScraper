package main

import (
	"log"
	"net"

	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price_monitoring/internal/db"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price_monitoring/internal/service"
	proto "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price_monitoring/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	postg_conn, err := postgres_db.InitDB()
	if err != nil {
		log.Fatalf("cannot connect to Postgres %v", err)
	}

	GrpcServer := grpc.NewServer()
	reflection.Register(GrpcServer)

	proto.RegisterScraperServer(GrpcServer, &service.Service{Postg: postg_conn})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	log.Println("WebScraper is running")

	err = GrpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
