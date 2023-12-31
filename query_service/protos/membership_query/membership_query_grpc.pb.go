// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: membership_query.proto

package membershipQueryService

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

// MembershipQueryServiceClient is the client API for MembershipQueryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MembershipQueryServiceClient interface {
	CreateMembership(ctx context.Context, in *CreateMembershipReq, opts ...grpc.CallOption) (*CreateMembershipRes, error)
	UpdateMembership(ctx context.Context, in *UpdateMembershipReq, opts ...grpc.CallOption) (*UpdateMembershipRes, error)
	GetMembershipById(ctx context.Context, in *GetMembershipByIdReq, opts ...grpc.CallOption) (*GetMembershipByIdRes, error)
	DeleteMembershipByID(ctx context.Context, in *DeleteMembershipByIdReq, opts ...grpc.CallOption) (*DeleteMembershipByIdRes, error)
	GetGroupMembership(ctx context.Context, in *GetGroupMembershipReq, opts ...grpc.CallOption) (*GetGroupMembershipRes, error)
	GetUserMembership(ctx context.Context, in *GetUserMembershipReq, opts ...grpc.CallOption) (*GetUserMembershipRes, error)
}

type membershipQueryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMembershipQueryServiceClient(cc grpc.ClientConnInterface) MembershipQueryServiceClient {
	return &membershipQueryServiceClient{cc}
}

func (c *membershipQueryServiceClient) CreateMembership(ctx context.Context, in *CreateMembershipReq, opts ...grpc.CallOption) (*CreateMembershipRes, error) {
	out := new(CreateMembershipRes)
	err := c.cc.Invoke(ctx, "/membershipQueryService.membershipQueryService/CreateMembership", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *membershipQueryServiceClient) UpdateMembership(ctx context.Context, in *UpdateMembershipReq, opts ...grpc.CallOption) (*UpdateMembershipRes, error) {
	out := new(UpdateMembershipRes)
	err := c.cc.Invoke(ctx, "/membershipQueryService.membershipQueryService/UpdateMembership", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *membershipQueryServiceClient) GetMembershipById(ctx context.Context, in *GetMembershipByIdReq, opts ...grpc.CallOption) (*GetMembershipByIdRes, error) {
	out := new(GetMembershipByIdRes)
	err := c.cc.Invoke(ctx, "/membershipQueryService.membershipQueryService/GetMembershipById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *membershipQueryServiceClient) DeleteMembershipByID(ctx context.Context, in *DeleteMembershipByIdReq, opts ...grpc.CallOption) (*DeleteMembershipByIdRes, error) {
	out := new(DeleteMembershipByIdRes)
	err := c.cc.Invoke(ctx, "/membershipQueryService.membershipQueryService/DeleteMembershipByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *membershipQueryServiceClient) GetGroupMembership(ctx context.Context, in *GetGroupMembershipReq, opts ...grpc.CallOption) (*GetGroupMembershipRes, error) {
	out := new(GetGroupMembershipRes)
	err := c.cc.Invoke(ctx, "/membershipQueryService.membershipQueryService/GetGroupMembership", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *membershipQueryServiceClient) GetUserMembership(ctx context.Context, in *GetUserMembershipReq, opts ...grpc.CallOption) (*GetUserMembershipRes, error) {
	out := new(GetUserMembershipRes)
	err := c.cc.Invoke(ctx, "/membershipQueryService.membershipQueryService/GetUserMembership", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MembershipQueryServiceServer is the server API for MembershipQueryService service.
// All implementations should embed UnimplementedMembershipQueryServiceServer
// for forward compatibility
type MembershipQueryServiceServer interface {
	CreateMembership(context.Context, *CreateMembershipReq) (*CreateMembershipRes, error)
	UpdateMembership(context.Context, *UpdateMembershipReq) (*UpdateMembershipRes, error)
	GetMembershipById(context.Context, *GetMembershipByIdReq) (*GetMembershipByIdRes, error)
	DeleteMembershipByID(context.Context, *DeleteMembershipByIdReq) (*DeleteMembershipByIdRes, error)
	GetGroupMembership(context.Context, *GetGroupMembershipReq) (*GetGroupMembershipRes, error)
	GetUserMembership(context.Context, *GetUserMembershipReq) (*GetUserMembershipRes, error)
}

// UnimplementedMembershipQueryServiceServer should be embedded to have forward compatible implementations.
type UnimplementedMembershipQueryServiceServer struct {
}

func (UnimplementedMembershipQueryServiceServer) CreateMembership(context.Context, *CreateMembershipReq) (*CreateMembershipRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMembership not implemented")
}
func (UnimplementedMembershipQueryServiceServer) UpdateMembership(context.Context, *UpdateMembershipReq) (*UpdateMembershipRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMembership not implemented")
}
func (UnimplementedMembershipQueryServiceServer) GetMembershipById(context.Context, *GetMembershipByIdReq) (*GetMembershipByIdRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMembershipById not implemented")
}
func (UnimplementedMembershipQueryServiceServer) DeleteMembershipByID(context.Context, *DeleteMembershipByIdReq) (*DeleteMembershipByIdRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMembershipByID not implemented")
}
func (UnimplementedMembershipQueryServiceServer) GetGroupMembership(context.Context, *GetGroupMembershipReq) (*GetGroupMembershipRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupMembership not implemented")
}
func (UnimplementedMembershipQueryServiceServer) GetUserMembership(context.Context, *GetUserMembershipReq) (*GetUserMembershipRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserMembership not implemented")
}

// UnsafeMembershipQueryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MembershipQueryServiceServer will
// result in compilation errors.
type UnsafeMembershipQueryServiceServer interface {
	mustEmbedUnimplementedMembershipQueryServiceServer()
}

func RegisterMembershipQueryServiceServer(s grpc.ServiceRegistrar, srv MembershipQueryServiceServer) {
	s.RegisterService(&MembershipQueryService_ServiceDesc, srv)
}

func _MembershipQueryService_CreateMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMembershipReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MembershipQueryServiceServer).CreateMembership(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/membershipQueryService.membershipQueryService/CreateMembership",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MembershipQueryServiceServer).CreateMembership(ctx, req.(*CreateMembershipReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MembershipQueryService_UpdateMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMembershipReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MembershipQueryServiceServer).UpdateMembership(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/membershipQueryService.membershipQueryService/UpdateMembership",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MembershipQueryServiceServer).UpdateMembership(ctx, req.(*UpdateMembershipReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MembershipQueryService_GetMembershipById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMembershipByIdReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MembershipQueryServiceServer).GetMembershipById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/membershipQueryService.membershipQueryService/GetMembershipById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MembershipQueryServiceServer).GetMembershipById(ctx, req.(*GetMembershipByIdReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MembershipQueryService_DeleteMembershipByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMembershipByIdReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MembershipQueryServiceServer).DeleteMembershipByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/membershipQueryService.membershipQueryService/DeleteMembershipByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MembershipQueryServiceServer).DeleteMembershipByID(ctx, req.(*DeleteMembershipByIdReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MembershipQueryService_GetGroupMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupMembershipReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MembershipQueryServiceServer).GetGroupMembership(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/membershipQueryService.membershipQueryService/GetGroupMembership",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MembershipQueryServiceServer).GetGroupMembership(ctx, req.(*GetGroupMembershipReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MembershipQueryService_GetUserMembership_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserMembershipReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MembershipQueryServiceServer).GetUserMembership(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/membershipQueryService.membershipQueryService/GetUserMembership",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MembershipQueryServiceServer).GetUserMembership(ctx, req.(*GetUserMembershipReq))
	}
	return interceptor(ctx, in, info, handler)
}

// MembershipQueryService_ServiceDesc is the grpc.ServiceDesc for MembershipQueryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MembershipQueryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "membershipQueryService.membershipQueryService",
	HandlerType: (*MembershipQueryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMembership",
			Handler:    _MembershipQueryService_CreateMembership_Handler,
		},
		{
			MethodName: "UpdateMembership",
			Handler:    _MembershipQueryService_UpdateMembership_Handler,
		},
		{
			MethodName: "GetMembershipById",
			Handler:    _MembershipQueryService_GetMembershipById_Handler,
		},
		{
			MethodName: "DeleteMembershipByID",
			Handler:    _MembershipQueryService_DeleteMembershipByID_Handler,
		},
		{
			MethodName: "GetGroupMembership",
			Handler:    _MembershipQueryService_GetGroupMembership_Handler,
		},
		{
			MethodName: "GetUserMembership",
			Handler:    _MembershipQueryService_GetUserMembership_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "membership_query.proto",
}
