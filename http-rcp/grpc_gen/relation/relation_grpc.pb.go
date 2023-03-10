// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: idl/relation.proto

package relation

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

// RelationManagementClient is the client API for RelationManagement service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RelationManagementClient interface {
	CreateRelation(ctx context.Context, in *CreateRelationRequest, opts ...grpc.CallOption) (*CreateRelationResponse, error)
	UpdateRelation(ctx context.Context, in *UpdateRelationRequest, opts ...grpc.CallOption) (*UpdateRelationResponse, error)
	QueryRelation(ctx context.Context, in *QueryRelationRequest, opts ...grpc.CallOption) (*QueryRelationResponse, error)
	QueryRelations(ctx context.Context, in *QueryRelationsRequest, opts ...grpc.CallOption) (*QueryRelationsResponse, error)
}

type relationManagementClient struct {
	cc grpc.ClientConnInterface
}

func NewRelationManagementClient(cc grpc.ClientConnInterface) RelationManagementClient {
	return &relationManagementClient{cc}
}

func (c *relationManagementClient) CreateRelation(ctx context.Context, in *CreateRelationRequest, opts ...grpc.CallOption) (*CreateRelationResponse, error) {
	out := new(CreateRelationResponse)
	err := c.cc.Invoke(ctx, "/relation.RelationManagement/CreateRelation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationManagementClient) UpdateRelation(ctx context.Context, in *UpdateRelationRequest, opts ...grpc.CallOption) (*UpdateRelationResponse, error) {
	out := new(UpdateRelationResponse)
	err := c.cc.Invoke(ctx, "/relation.RelationManagement/UpdateRelation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationManagementClient) QueryRelation(ctx context.Context, in *QueryRelationRequest, opts ...grpc.CallOption) (*QueryRelationResponse, error) {
	out := new(QueryRelationResponse)
	err := c.cc.Invoke(ctx, "/relation.RelationManagement/QueryRelation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationManagementClient) QueryRelations(ctx context.Context, in *QueryRelationsRequest, opts ...grpc.CallOption) (*QueryRelationsResponse, error) {
	out := new(QueryRelationsResponse)
	err := c.cc.Invoke(ctx, "/relation.RelationManagement/QueryRelations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RelationManagementServer is the server API for RelationManagement service.
// All implementations must embed UnimplementedRelationManagementServer
// for forward compatibility
type RelationManagementServer interface {
	CreateRelation(context.Context, *CreateRelationRequest) (*CreateRelationResponse, error)
	UpdateRelation(context.Context, *UpdateRelationRequest) (*UpdateRelationResponse, error)
	QueryRelation(context.Context, *QueryRelationRequest) (*QueryRelationResponse, error)
	QueryRelations(context.Context, *QueryRelationsRequest) (*QueryRelationsResponse, error)
	mustEmbedUnimplementedRelationManagementServer()
}

// UnimplementedRelationManagementServer must be embedded to have forward compatible implementations.
type UnimplementedRelationManagementServer struct {
}

func (UnimplementedRelationManagementServer) CreateRelation(context.Context, *CreateRelationRequest) (*CreateRelationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRelation not implemented")
}
func (UnimplementedRelationManagementServer) UpdateRelation(context.Context, *UpdateRelationRequest) (*UpdateRelationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRelation not implemented")
}
func (UnimplementedRelationManagementServer) QueryRelation(context.Context, *QueryRelationRequest) (*QueryRelationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryRelation not implemented")
}
func (UnimplementedRelationManagementServer) QueryRelations(context.Context, *QueryRelationsRequest) (*QueryRelationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryRelations not implemented")
}
func (UnimplementedRelationManagementServer) mustEmbedUnimplementedRelationManagementServer() {}

// UnsafeRelationManagementServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RelationManagementServer will
// result in compilation errors.
type UnsafeRelationManagementServer interface {
	mustEmbedUnimplementedRelationManagementServer()
}

func RegisterRelationManagementServer(s grpc.ServiceRegistrar, srv RelationManagementServer) {
	s.RegisterService(&RelationManagement_ServiceDesc, srv)
}

func _RelationManagement_CreateRelation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRelationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationManagementServer).CreateRelation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relation.RelationManagement/CreateRelation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationManagementServer).CreateRelation(ctx, req.(*CreateRelationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelationManagement_UpdateRelation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRelationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationManagementServer).UpdateRelation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relation.RelationManagement/UpdateRelation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationManagementServer).UpdateRelation(ctx, req.(*UpdateRelationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelationManagement_QueryRelation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRelationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationManagementServer).QueryRelation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relation.RelationManagement/QueryRelation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationManagementServer).QueryRelation(ctx, req.(*QueryRelationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelationManagement_QueryRelations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRelationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationManagementServer).QueryRelations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relation.RelationManagement/QueryRelations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationManagementServer).QueryRelations(ctx, req.(*QueryRelationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RelationManagement_ServiceDesc is the grpc.ServiceDesc for RelationManagement service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RelationManagement_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "relation.RelationManagement",
	HandlerType: (*RelationManagementServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRelation",
			Handler:    _RelationManagement_CreateRelation_Handler,
		},
		{
			MethodName: "UpdateRelation",
			Handler:    _RelationManagement_UpdateRelation_Handler,
		},
		{
			MethodName: "QueryRelation",
			Handler:    _RelationManagement_QueryRelation_Handler,
		},
		{
			MethodName: "QueryRelations",
			Handler:    _RelationManagement_QueryRelations_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "idl/relation.proto",
}
