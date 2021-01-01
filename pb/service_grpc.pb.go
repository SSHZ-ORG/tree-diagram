// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// TreeDiagramServiceClient is the client API for TreeDiagramService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TreeDiagramServiceClient interface {
	RenderEvent(ctx context.Context, in *RenderEventRequest, opts ...grpc.CallOption) (*RenderEventResponse, error)
	RenderPlace(ctx context.Context, in *RenderPlaceRequest, opts ...grpc.CallOption) (*RenderPlaceResponse, error)
}

type treeDiagramServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTreeDiagramServiceClient(cc grpc.ClientConnInterface) TreeDiagramServiceClient {
	return &treeDiagramServiceClient{cc}
}

func (c *treeDiagramServiceClient) RenderEvent(ctx context.Context, in *RenderEventRequest, opts ...grpc.CallOption) (*RenderEventResponse, error) {
	out := new(RenderEventResponse)
	err := c.cc.Invoke(ctx, "/treediagram.pb.TreeDiagramService/RenderEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *treeDiagramServiceClient) RenderPlace(ctx context.Context, in *RenderPlaceRequest, opts ...grpc.CallOption) (*RenderPlaceResponse, error) {
	out := new(RenderPlaceResponse)
	err := c.cc.Invoke(ctx, "/treediagram.pb.TreeDiagramService/RenderPlace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TreeDiagramServiceServer is the server API for TreeDiagramService service.
// All implementations must embed UnimplementedTreeDiagramServiceServer
// for forward compatibility
type TreeDiagramServiceServer interface {
	RenderEvent(context.Context, *RenderEventRequest) (*RenderEventResponse, error)
	RenderPlace(context.Context, *RenderPlaceRequest) (*RenderPlaceResponse, error)
	mustEmbedUnimplementedTreeDiagramServiceServer()
}

// UnimplementedTreeDiagramServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTreeDiagramServiceServer struct {
}

func (UnimplementedTreeDiagramServiceServer) RenderEvent(context.Context, *RenderEventRequest) (*RenderEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RenderEvent not implemented")
}
func (UnimplementedTreeDiagramServiceServer) RenderPlace(context.Context, *RenderPlaceRequest) (*RenderPlaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RenderPlace not implemented")
}
func (UnimplementedTreeDiagramServiceServer) mustEmbedUnimplementedTreeDiagramServiceServer() {}

// UnsafeTreeDiagramServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TreeDiagramServiceServer will
// result in compilation errors.
type UnsafeTreeDiagramServiceServer interface {
	mustEmbedUnimplementedTreeDiagramServiceServer()
}

func RegisterTreeDiagramServiceServer(s grpc.ServiceRegistrar, srv TreeDiagramServiceServer) {
	s.RegisterService(&TreeDiagramService_ServiceDesc, srv)
}

func _TreeDiagramService_RenderEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RenderEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TreeDiagramServiceServer).RenderEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/treediagram.pb.TreeDiagramService/RenderEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TreeDiagramServiceServer).RenderEvent(ctx, req.(*RenderEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TreeDiagramService_RenderPlace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RenderPlaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TreeDiagramServiceServer).RenderPlace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/treediagram.pb.TreeDiagramService/RenderPlace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TreeDiagramServiceServer).RenderPlace(ctx, req.(*RenderPlaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TreeDiagramService_ServiceDesc is the grpc.ServiceDesc for TreeDiagramService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TreeDiagramService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "treediagram.pb.TreeDiagramService",
	HandlerType: (*TreeDiagramServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RenderEvent",
			Handler:    _TreeDiagramService_RenderEvent_Handler,
		},
		{
			MethodName: "RenderPlace",
			Handler:    _TreeDiagramService_RenderPlace_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
