// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: price_tracker.proto

package price_monitoring

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Scraper_GetItem_FullMethodName     = "/price_tracker.Scraper/GetItem"
	Scraper_GetAllItems_FullMethodName = "/price_tracker.Scraper/GetAllItems"
)

// ScraperClient is the client API for Scraper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ScraperClient interface {
	GetItem(ctx context.Context, in *GetItemRequest, opts ...grpc.CallOption) (*GetItemResponse, error)
	GetAllItems(ctx context.Context, in *GetAllItemsRequest, opts ...grpc.CallOption) (*GetAllItemsResponse, error)
}

type scraperClient struct {
	cc grpc.ClientConnInterface
}

func NewScraperClient(cc grpc.ClientConnInterface) ScraperClient {
	return &scraperClient{cc}
}

func (c *scraperClient) GetItem(ctx context.Context, in *GetItemRequest, opts ...grpc.CallOption) (*GetItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetItemResponse)
	err := c.cc.Invoke(ctx, Scraper_GetItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scraperClient) GetAllItems(ctx context.Context, in *GetAllItemsRequest, opts ...grpc.CallOption) (*GetAllItemsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllItemsResponse)
	err := c.cc.Invoke(ctx, Scraper_GetAllItems_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ScraperServer is the server API for Scraper service.
// All implementations must embed UnimplementedScraperServer
// for forward compatibility.
type ScraperServer interface {
	GetItem(context.Context, *GetItemRequest) (*GetItemResponse, error)
	GetAllItems(context.Context, *GetAllItemsRequest) (*GetAllItemsResponse, error)
	mustEmbedUnimplementedScraperServer()
}

// UnimplementedScraperServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedScraperServer struct{}

func (UnimplementedScraperServer) GetItem(context.Context, *GetItemRequest) (*GetItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItem not implemented")
}
func (UnimplementedScraperServer) GetAllItems(context.Context, *GetAllItemsRequest) (*GetAllItemsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllItems not implemented")
}
func (UnimplementedScraperServer) mustEmbedUnimplementedScraperServer() {}
func (UnimplementedScraperServer) testEmbeddedByValue()                 {}

// UnsafeScraperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ScraperServer will
// result in compilation errors.
type UnsafeScraperServer interface {
	mustEmbedUnimplementedScraperServer()
}

func RegisterScraperServer(s grpc.ServiceRegistrar, srv ScraperServer) {
	// If the following call pancis, it indicates UnimplementedScraperServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Scraper_ServiceDesc, srv)
}

func _Scraper_GetItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScraperServer).GetItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Scraper_GetItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScraperServer).GetItem(ctx, req.(*GetItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scraper_GetAllItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllItemsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScraperServer).GetAllItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Scraper_GetAllItems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScraperServer).GetAllItems(ctx, req.(*GetAllItemsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Scraper_ServiceDesc is the grpc.ServiceDesc for Scraper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Scraper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "price_tracker.Scraper",
	HandlerType: (*ScraperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetItem",
			Handler:    _Scraper_GetItem_Handler,
		},
		{
			MethodName: "GetAllItems",
			Handler:    _Scraper_GetAllItems_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "price_tracker.proto",
}
