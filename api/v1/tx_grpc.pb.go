// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: fatal_fruit/auction/v1/tx.proto

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
	Msg_NewAuction_FullMethodName   = "/fatal_fruit.auction.v1.Msg/NewAuction"
	Msg_StartAuction_FullMethodName = "/fatal_fruit.auction.v1.Msg/StartAuction"
	Msg_NewBid_FullMethodName       = "/fatal_fruit.auction.v1.Msg/NewBid"
	Msg_Exec_FullMethodName         = "/fatal_fruit.auction.v1.Msg/Exec"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	// NewAuction creates a new auction.
	NewAuction(ctx context.Context, in *MsgNewAuction, opts ...grpc.CallOption) (*MsgNewAuctionResponse, error)
	// StartAuction initializes the auction
	StartAuction(ctx context.Context, in *MsgStartAuction, opts ...grpc.CallOption) (*MsgStartAuctionResponse, error)
	// NewBid places a new bid on an auction.
	NewBid(ctx context.Context, in *MsgNewBid, opts ...grpc.CallOption) (*MsgNewBidResponse, error)
	// Exec executes an auction, distributing funds and finalizing the auction.
	Exec(ctx context.Context, in *MsgExecAuction, opts ...grpc.CallOption) (*MsgExecAuctionResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) NewAuction(ctx context.Context, in *MsgNewAuction, opts ...grpc.CallOption) (*MsgNewAuctionResponse, error) {
	out := new(MsgNewAuctionResponse)
	err := c.cc.Invoke(ctx, Msg_NewAuction_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) StartAuction(ctx context.Context, in *MsgStartAuction, opts ...grpc.CallOption) (*MsgStartAuctionResponse, error) {
	out := new(MsgStartAuctionResponse)
	err := c.cc.Invoke(ctx, Msg_StartAuction_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) NewBid(ctx context.Context, in *MsgNewBid, opts ...grpc.CallOption) (*MsgNewBidResponse, error) {
	out := new(MsgNewBidResponse)
	err := c.cc.Invoke(ctx, Msg_NewBid_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Exec(ctx context.Context, in *MsgExecAuction, opts ...grpc.CallOption) (*MsgExecAuctionResponse, error) {
	out := new(MsgExecAuctionResponse)
	err := c.cc.Invoke(ctx, Msg_Exec_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	// NewAuction creates a new auction.
	NewAuction(context.Context, *MsgNewAuction) (*MsgNewAuctionResponse, error)
	// StartAuction initializes the auction
	StartAuction(context.Context, *MsgStartAuction) (*MsgStartAuctionResponse, error)
	// NewBid places a new bid on an auction.
	NewBid(context.Context, *MsgNewBid) (*MsgNewBidResponse, error)
	// Exec executes an auction, distributing funds and finalizing the auction.
	Exec(context.Context, *MsgExecAuction) (*MsgExecAuctionResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) NewAuction(context.Context, *MsgNewAuction) (*MsgNewAuctionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewAuction not implemented")
}
func (UnimplementedMsgServer) StartAuction(context.Context, *MsgStartAuction) (*MsgStartAuctionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartAuction not implemented")
}
func (UnimplementedMsgServer) NewBid(context.Context, *MsgNewBid) (*MsgNewBidResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewBid not implemented")
}
func (UnimplementedMsgServer) Exec(context.Context, *MsgExecAuction) (*MsgExecAuctionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exec not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_NewAuction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgNewAuction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).NewAuction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_NewAuction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).NewAuction(ctx, req.(*MsgNewAuction))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_StartAuction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgStartAuction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).StartAuction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_StartAuction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).StartAuction(ctx, req.(*MsgStartAuction))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_NewBid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgNewBid)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).NewBid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_NewBid_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).NewBid(ctx, req.(*MsgNewBid))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Exec_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgExecAuction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Exec(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_Exec_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Exec(ctx, req.(*MsgExecAuction))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fatal_fruit.auction.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewAuction",
			Handler:    _Msg_NewAuction_Handler,
		},
		{
			MethodName: "StartAuction",
			Handler:    _Msg_StartAuction_Handler,
		},
		{
			MethodName: "NewBid",
			Handler:    _Msg_NewBid_Handler,
		},
		{
			MethodName: "Exec",
			Handler:    _Msg_Exec_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fatal_fruit/auction/v1/tx.proto",
}
