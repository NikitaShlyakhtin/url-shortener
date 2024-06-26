// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: internal/proto/url_shortener.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UrlShortenerClient is the client API for UrlShortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UrlShortenerClient interface {
	CreateShortUrl(ctx context.Context, in *OriginalUrl, opts ...grpc.CallOption) (*ShortUrl, error)
	GetOriginalUrl(ctx context.Context, in *ShortUrl, opts ...grpc.CallOption) (*OriginalUrl, error)
}

type urlShortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewUrlShortenerClient(cc grpc.ClientConnInterface) UrlShortenerClient {
	return &urlShortenerClient{cc}
}

func (c *urlShortenerClient) CreateShortUrl(ctx context.Context, in *OriginalUrl, opts ...grpc.CallOption) (*ShortUrl, error) {
	out := new(ShortUrl)
	err := c.cc.Invoke(ctx, "/proto.UrlShortener/CreateShortUrl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *urlShortenerClient) GetOriginalUrl(ctx context.Context, in *ShortUrl, opts ...grpc.CallOption) (*OriginalUrl, error) {
	out := new(OriginalUrl)
	err := c.cc.Invoke(ctx, "/proto.UrlShortener/GetOriginalUrl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UrlShortenerServer is the server API for UrlShortener service.
// All implementations must embed UnimplementedUrlShortenerServer
// for forward compatibility
type UrlShortenerServer interface {
	CreateShortUrl(context.Context, *OriginalUrl) (*ShortUrl, error)
	GetOriginalUrl(context.Context, *ShortUrl) (*OriginalUrl, error)
	mustEmbedUnimplementedUrlShortenerServer()
}

// UnimplementedUrlShortenerServer must be embedded to have forward compatible implementations.
type UnimplementedUrlShortenerServer struct {
}

func (UnimplementedUrlShortenerServer) CreateShortUrl(context.Context, *OriginalUrl) (*ShortUrl, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateShortUrl not implemented")
}
func (UnimplementedUrlShortenerServer) GetOriginalUrl(context.Context, *ShortUrl) (*OriginalUrl, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOriginalUrl not implemented")
}
func (UnimplementedUrlShortenerServer) mustEmbedUnimplementedUrlShortenerServer() {}

// UnsafeUrlShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UrlShortenerServer will
// result in compilation errors.
type UnsafeUrlShortenerServer interface {
	mustEmbedUnimplementedUrlShortenerServer()
}

func RegisterUrlShortenerServer(s grpc.ServiceRegistrar, srv UrlShortenerServer) {
	s.RegisterService(&UrlShortener_ServiceDesc, srv)
}

func _UrlShortener_CreateShortUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OriginalUrl)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UrlShortenerServer).CreateShortUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UrlShortener/CreateShortUrl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UrlShortenerServer).CreateShortUrl(ctx, req.(*OriginalUrl))
	}
	return interceptor(ctx, in, info, handler)
}

func _UrlShortener_GetOriginalUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortUrl)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UrlShortenerServer).GetOriginalUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.UrlShortener/GetOriginalUrl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UrlShortenerServer).GetOriginalUrl(ctx, req.(*ShortUrl))
	}
	return interceptor(ctx, in, info, handler)
}

// UrlShortener_ServiceDesc is the grpc.ServiceDesc for UrlShortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UrlShortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.UrlShortener",
	HandlerType: (*UrlShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateShortUrl",
			Handler:    _UrlShortener_CreateShortUrl_Handler,
		},
		{
			MethodName: "GetOriginalUrl",
			Handler:    _UrlShortener_GetOriginalUrl_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/url_shortener.proto",
}
