// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: fatal_fruit/auction/v1/query.proto

package auctionv1

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

const (
	Query_Auction_FullMethodName       = "/fatal_fruit.auction.v1.Query/Auction"
	Query_OwnerAuctions_FullMethodName = "/fatal_fruit.auction.v1.Query/OwnerAuctions"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	Auction(ctx context.Context, in *QueryAuctionRequest, opts ...grpc.CallOption) (*QueryAuctionResponse, error)
	OwnerAuctions(ctx context.Context, in *QueryOwnerAuctionsRequest, opts ...grpc.CallOption) (*QueryOwnerAuctionsResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Auction(ctx context.Context, in *QueryAuctionRequest, opts ...grpc.CallOption) (*QueryAuctionResponse, error) {
	out := new(QueryAuctionResponse)
	err := c.cc.Invoke(ctx, Query_Auction_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) OwnerAuctions(ctx context.Context, in *QueryOwnerAuctionsRequest, opts ...grpc.CallOption) (*QueryOwnerAuctionsResponse, error) {
	out := new(QueryOwnerAuctionsResponse)
	err := c.cc.Invoke(ctx, Query_OwnerAuctions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	Auction(context.Context, *QueryAuctionRequest) (*QueryAuctionResponse, error)
	OwnerAuctions(context.Context, *QueryOwnerAuctionsRequest) (*QueryOwnerAuctionsResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Auction(context.Context, *QueryAuctionRequest) (*QueryAuctionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Auction not implemented")
}
func (UnimplementedQueryServer) OwnerAuctions(context.Context, *QueryOwnerAuctionsRequest) (*QueryOwnerAuctionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OwnerAuctions not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Auction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAuctionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Auction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Auction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Auction(ctx, req.(*QueryAuctionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_OwnerAuctions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryOwnerAuctionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).OwnerAuctions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_OwnerAuctions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).OwnerAuctions(ctx, req.(*QueryOwnerAuctionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fatal_fruit.auction.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Auction",
			Handler:    _Query_Auction_Handler,
		},
		{
			MethodName: "OwnerAuctions",
			Handler:    _Query_OwnerAuctions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fatal_fruit/auction/v1/query.proto",
}
