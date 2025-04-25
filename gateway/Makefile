generate-grpc-api:
	mkdir -p pkg/pb/auth/
	
	protoc \
	--proto_path api/auth \
	--go_out . \
	--go-grpc_out . \
	api/auth/auth.proto

	protoc \
	--proto_path api/tracker \
	--go_out . \
	--go-grpc_out . \
	api/tracker/price_tracker.proto