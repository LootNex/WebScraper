module gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/gateway

go 1.23.4

replace gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth => ../auth

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.37.0
	google.golang.org/grpc v1.71.1
	google.golang.org/protobuf v1.36.6
)

require (
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
)
