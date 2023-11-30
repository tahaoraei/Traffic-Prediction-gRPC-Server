// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: contract/protobuf/time.proto

package time

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

// GetETAClient is the client API for GetETA service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GetETAClient interface {
	GetETA(ctx context.Context, in *TravelRequest, opts ...grpc.CallOption) (*TravelResponse, error)
}

type getETAClient struct {
	cc grpc.ClientConnInterface
}

func NewGetETAClient(cc grpc.ClientConnInterface) GetETAClient {
	return &getETAClient{cc}
}

func (c *getETAClient) GetETA(ctx context.Context, in *TravelRequest, opts ...grpc.CallOption) (*TravelResponse, error) {
	out := new(TravelResponse)
	err := c.cc.Invoke(ctx, "/protobuf.GetETA/GetETA", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetETAServer is the server API for GetETA service.
// All implementations must embed UnimplementedGetETAServer
// for forward compatibility
type GetETAServer interface {
	GetETA(context.Context, *TravelRequest) (*TravelResponse, error)
	mustEmbedUnimplementedGetETAServer()
}

// UnimplementedGetETAServer must be embedded to have forward compatible implementations.
type UnimplementedGetETAServer struct {
}

func (UnimplementedGetETAServer) GetETA(context.Context, *TravelRequest) (*TravelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetETA not implemented")
}
func (UnimplementedGetETAServer) mustEmbedUnimplementedGetETAServer() {}

// UnsafeGetETAServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GetETAServer will
// result in compilation errors.
type UnsafeGetETAServer interface {
	mustEmbedUnimplementedGetETAServer()
}

func RegisterGetETAServer(s grpc.ServiceRegistrar, srv GetETAServer) {
	s.RegisterService(&GetETA_ServiceDesc, srv)
}

func _GetETA_GetETA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TravelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GetETAServer).GetETA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.GetETA/GetETA",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GetETAServer).GetETA(ctx, req.(*TravelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GetETA_ServiceDesc is the grpc.ServiceDesc for GetETA service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GetETA_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.GetETA",
	HandlerType: (*GetETAServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetETA",
			Handler:    _GetETA_GetETA_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "contract/protobuf/time.proto",
}