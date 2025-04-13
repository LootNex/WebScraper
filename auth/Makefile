generate-grpc-api:
	mkdir -p pkg/pb/auth/
	
	protoc \
	--proto_path api/auth \
	--go_out . \
	--go-grpc_out . \
	api/auth/auth.proto
