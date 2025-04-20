package grpc

import(
	authpb "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/pkg/pb/auth"
)

type Client struct{
	api authpb.Auth_V1Client
}