// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: internal/protobuf/resticky.proto

package resticky

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
	RestickyService_LockAll_FullMethodName   = "/resticky.RestickyService/LockAll"
	RestickyService_UnlockAll_FullMethodName = "/resticky.RestickyService/UnlockAll"
)

// RestickyServiceClient is the client API for RestickyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RestickyServiceClient interface {
	LockAll(ctx context.Context, in *RestickyRequest, opts ...grpc.CallOption) (*RestickyResponse, error)
	UnlockAll(ctx context.Context, in *RestickyRequest, opts ...grpc.CallOption) (*RestickyResponse, error)
}

type restickyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRestickyServiceClient(cc grpc.ClientConnInterface) RestickyServiceClient {
	return &restickyServiceClient{cc}
}

func (c *restickyServiceClient) LockAll(ctx context.Context, in *RestickyRequest, opts ...grpc.CallOption) (*RestickyResponse, error) {
	out := new(RestickyResponse)
	err := c.cc.Invoke(ctx, RestickyService_LockAll_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *restickyServiceClient) UnlockAll(ctx context.Context, in *RestickyRequest, opts ...grpc.CallOption) (*RestickyResponse, error) {
	out := new(RestickyResponse)
	err := c.cc.Invoke(ctx, RestickyService_UnlockAll_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RestickyServiceServer is the server API for RestickyService service.
// All implementations must embed UnimplementedRestickyServiceServer
// for forward compatibility
type RestickyServiceServer interface {
	LockAll(context.Context, *RestickyRequest) (*RestickyResponse, error)
	UnlockAll(context.Context, *RestickyRequest) (*RestickyResponse, error)
	mustEmbedUnimplementedRestickyServiceServer()
}

// UnimplementedRestickyServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRestickyServiceServer struct {
}

func (UnimplementedRestickyServiceServer) LockAll(context.Context, *RestickyRequest) (*RestickyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LockAll not implemented")
}
func (UnimplementedRestickyServiceServer) UnlockAll(context.Context, *RestickyRequest) (*RestickyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnlockAll not implemented")
}
func (UnimplementedRestickyServiceServer) mustEmbedUnimplementedRestickyServiceServer() {}

// UnsafeRestickyServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RestickyServiceServer will
// result in compilation errors.
type UnsafeRestickyServiceServer interface {
	mustEmbedUnimplementedRestickyServiceServer()
}

func RegisterRestickyServiceServer(s grpc.ServiceRegistrar, srv RestickyServiceServer) {
	s.RegisterService(&RestickyService_ServiceDesc, srv)
}

func _RestickyService_LockAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RestickyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RestickyServiceServer).LockAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RestickyService_LockAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RestickyServiceServer).LockAll(ctx, req.(*RestickyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RestickyService_UnlockAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RestickyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RestickyServiceServer).UnlockAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RestickyService_UnlockAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RestickyServiceServer).UnlockAll(ctx, req.(*RestickyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RestickyService_ServiceDesc is the grpc.ServiceDesc for RestickyService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RestickyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "resticky.RestickyService",
	HandlerType: (*RestickyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LockAll",
			Handler:    _RestickyService_LockAll_Handler,
		},
		{
			MethodName: "UnlockAll",
			Handler:    _RestickyService_UnlockAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/protobuf/resticky.proto",
}
