package grpc

import(
	authpb "gitlab.crja72.ru/golang/2025/spring/course/projects/go2/price-tracker/auth/pkg/pb"
)

type Client struct{
	api authpb.AuthClient
}