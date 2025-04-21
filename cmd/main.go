package main

import (
	"log"
	// "fmt"
	"net"

	proto "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price-monitoring/gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price-monitoring"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price-monitoring/internal/db"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price-monitoring/internal/driver"
	// "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price-monitoring/internal/parser"
	"gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/price-monitoring/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	driver, finish, err := driver.InitDriver()
	if err != nil {
		log.Fatalf("cannot init chrome driver %v", err)
	}
	defer finish()

	// link := "https://www.ozon.ru/product/zero-mileage-5w-30-maslo-motornoe-sinteticheskoe-1-l-1627408607/?at=Eqtk44V8ghrNGKRJTOPRG9LS0zoNA7UwDRn5mtNR6r8o"

	// fmt.Println(parser.OzonParser(link, driver))

	postg_conn, err := postgres_db.InitDB()
	if err != nil {
		log.Fatalf("cannot connect to Postgres %v", err)
	}

	GrpcServer := grpc.NewServer()
	reflection.Register(GrpcServer)

	proto.RegisterScraperServer(GrpcServer, &service.Service{Postg: postg_conn, Driver: driver})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil{
		log.Fatalf("failed to listen %v", err)
	}

	err = GrpcServer.Serve(lis)
	if err != nil{
		log.Fatalf("failed to serve: %v", err)
	}
	

}
