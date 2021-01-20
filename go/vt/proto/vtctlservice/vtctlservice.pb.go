// Code generated by protoc-gen-go. DO NOT EDIT.
// source: vtctlservice.proto

package vtctlservice

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	vtctldata "vitess.io/vitess/go/vt/proto/vtctldata"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("vtctlservice.proto", fileDescriptor_27055cdbb1148d2b) }

var fileDescriptor_27055cdbb1148d2b = []byte{
	// 375 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x94, 0xdf, 0x4b, 0x02, 0x41,
	0x10, 0xc7, 0xeb, 0x21, 0xa1, 0xcd, 0x30, 0xa6, 0x87, 0xc0, 0xd2, 0xca, 0x28, 0xb0, 0xc0, 0x0b,
	0xfb, 0x0b, 0x4c, 0xca, 0x44, 0x90, 0x4a, 0xf1, 0x41, 0xe8, 0x61, 0xbd, 0x9b, 0xf2, 0xe0, 0x7e,
	0xe8, 0xcd, 0x7a, 0xe4, 0x3f, 0xde, 0x73, 0x74, 0xb6, 0xeb, 0x7a, 0xde, 0x5a, 0x6f, 0xde, 0x7c,
	0xbe, 0xf3, 0xd9, 0x71, 0x19, 0x96, 0x41, 0x2c, 0x6c, 0xe1, 0x11, 0x46, 0xb1, 0x6b, 0x63, 0x6d,
	0x12, 0x85, 0x22, 0x84, 0xbc, 0x5e, 0x2b, 0x16, 0x92, 0x2f, 0x87, 0x0b, 0xbe, 0xc0, 0xf5, 0x29,
	0xdb, 0x19, 0xfc, 0x94, 0x60, 0xcc, 0x0e, 0x1f, 0x3e, 0xd1, 0x9e, 0x09, 0x4c, 0xbe, 0x9b, 0xa1,
	0xef, 0xf3, 0xc0, 0x81, 0xcb, 0xda, 0xb2, 0x23, 0x83, 0xbf, 0xe2, 0x74, 0x86, 0x24, 0x8a, 0x57,
	0x7f, 0xc5, 0x68, 0x12, 0x06, 0x84, 0x95, 0xad, 0xdb, 0xed, 0xfa, 0x57, 0x8e, 0xe5, 0x12, 0xe8,
	0x40, 0xc4, 0x8e, 0x1e, 0xdd, 0xc0, 0x69, 0x78, 0x5e, 0x6f, 0xcc, 0x23, 0x87, 0xda, 0x41, 0x07,
	0xe7, 0x34, 0xe1, 0x36, 0x42, 0x55, 0x33, 0x1a, 0x32, 0xf2, 0xf0, 0xeb, 0xff, 0x44, 0xe5, 0x00,
	0xf0, 0xc6, 0x0e, 0x5a, 0x28, 0x9a, 0xe8, 0x79, 0xed, 0xe0, 0x3d, 0xec, 0x72, 0x1f, 0x09, 0x2a,
	0x9a, 0x21, 0x0d, 0xe5, 0x29, 0x17, 0x1b, 0x33, 0x4a, 0xdf, 0x65, 0x7b, 0x1a, 0x85, 0x52, 0x76,
	0x97, 0x94, 0x96, 0x4d, 0x58, 0xf9, 0x86, 0xac, 0xf0, 0x0b, 0xa8, 0xe1, 0xb9, 0x9c, 0x90, 0xe0,
	0x7c, 0xbd, 0x49, 0x32, 0xe9, 0xad, 0x6c, 0x8a, 0xa4, 0x66, 0x55, 0x57, 0x9e, 0x9a, 0x35, 0x7d,
	0xcd, 0x65, 0x13, 0x56, 0xbe, 0x17, 0x96, 0xd7, 0x00, 0x81, 0xa1, 0x43, 0x4d, 0x79, 0x6a, 0xe4,
	0x4a, 0xd9, 0x67, 0xfb, 0x2d, 0x14, 0xbd, 0x28, 0x1e, 0xf4, 0xec, 0x31, 0xfa, 0x1c, 0x52, 0x3d,
	0x4b, 0x22, 0xa5, 0x67, 0xe6, 0x80, 0xb2, 0x3e, 0xb1, 0xdd, 0x16, 0x8a, 0x3e, 0x1f, 0x79, 0x28,
	0xe0, 0x78, 0xb5, 0x61, 0x51, 0x95, 0xb6, 0x93, 0x6c, 0xa8, 0x4c, 0x1d, 0xc6, 0x54, 0x99, 0x20,
	0x33, 0xad, 0xfe, 0x6e, 0xc9, 0x40, 0xf5, 0xd5, 0x6c, 0x07, 0xae, 0x48, 0x96, 0xf7, 0x39, 0x72,
	0x7d, 0x1e, 0xcd, 0x57, 0x56, 0x33, 0x0d, 0xb3, 0x56, 0x73, 0x3d, 0x23, 0xf5, 0xf7, 0x37, 0xc3,
	0x6a, 0xec, 0x0a, 0x24, 0xaa, 0xb9, 0xa1, 0xb5, 0xf8, 0x65, 0x7d, 0x84, 0x56, 0x2c, 0xac, 0xe4,
	0x2d, 0xb0, 0xf4, 0x97, 0x62, 0x94, 0x4b, 0x6a, 0x77, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x12,
	0x19, 0xd7, 0x32, 0x54, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// VtctlClient is the client API for Vtctl service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type VtctlClient interface {
	ExecuteVtctlCommand(ctx context.Context, in *vtctldata.ExecuteVtctlCommandRequest, opts ...grpc.CallOption) (Vtctl_ExecuteVtctlCommandClient, error)
}

type vtctlClient struct {
	cc *grpc.ClientConn
}

func NewVtctlClient(cc *grpc.ClientConn) VtctlClient {
	return &vtctlClient{cc}
}

func (c *vtctlClient) ExecuteVtctlCommand(ctx context.Context, in *vtctldata.ExecuteVtctlCommandRequest, opts ...grpc.CallOption) (Vtctl_ExecuteVtctlCommandClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Vtctl_serviceDesc.Streams[0], "/vtctlservice.Vtctl/ExecuteVtctlCommand", opts...)
	if err != nil {
		return nil, err
	}
	x := &vtctlExecuteVtctlCommandClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Vtctl_ExecuteVtctlCommandClient interface {
	Recv() (*vtctldata.ExecuteVtctlCommandResponse, error)
	grpc.ClientStream
}

type vtctlExecuteVtctlCommandClient struct {
	grpc.ClientStream
}

func (x *vtctlExecuteVtctlCommandClient) Recv() (*vtctldata.ExecuteVtctlCommandResponse, error) {
	m := new(vtctldata.ExecuteVtctlCommandResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// VtctlServer is the server API for Vtctl service.
type VtctlServer interface {
	ExecuteVtctlCommand(*vtctldata.ExecuteVtctlCommandRequest, Vtctl_ExecuteVtctlCommandServer) error
}

// UnimplementedVtctlServer can be embedded to have forward compatible implementations.
type UnimplementedVtctlServer struct {
}

func (*UnimplementedVtctlServer) ExecuteVtctlCommand(req *vtctldata.ExecuteVtctlCommandRequest, srv Vtctl_ExecuteVtctlCommandServer) error {
	return status.Errorf(codes.Unimplemented, "method ExecuteVtctlCommand not implemented")
}

func RegisterVtctlServer(s *grpc.Server, srv VtctlServer) {
	s.RegisterService(&_Vtctl_serviceDesc, srv)
}

func _Vtctl_ExecuteVtctlCommand_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(vtctldata.ExecuteVtctlCommandRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(VtctlServer).ExecuteVtctlCommand(m, &vtctlExecuteVtctlCommandServer{stream})
}

type Vtctl_ExecuteVtctlCommandServer interface {
	Send(*vtctldata.ExecuteVtctlCommandResponse) error
	grpc.ServerStream
}

type vtctlExecuteVtctlCommandServer struct {
	grpc.ServerStream
}

func (x *vtctlExecuteVtctlCommandServer) Send(m *vtctldata.ExecuteVtctlCommandResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _Vtctl_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vtctlservice.Vtctl",
	HandlerType: (*VtctlServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ExecuteVtctlCommand",
			Handler:       _Vtctl_ExecuteVtctlCommand_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "vtctlservice.proto",
}

// VtctldClient is the client API for Vtctld service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type VtctldClient interface {
	// FindAllShardsInKeyspace returns a map of shard names to shard references
	// for a given keyspace.
	FindAllShardsInKeyspace(ctx context.Context, in *vtctldata.FindAllShardsInKeyspaceRequest, opts ...grpc.CallOption) (*vtctldata.FindAllShardsInKeyspaceResponse, error)
	// GetCellInfoNames returns all the cells for which we have a CellInfo object,
	// meaning we have a topology service registered.
	GetCellInfoNames(ctx context.Context, in *vtctldata.GetCellInfoNamesRequest, opts ...grpc.CallOption) (*vtctldata.GetCellInfoNamesResponse, error)
	// GetCellInfo returns the information for a cell.
	GetCellInfo(ctx context.Context, in *vtctldata.GetCellInfoRequest, opts ...grpc.CallOption) (*vtctldata.GetCellInfoResponse, error)
	// GetCellsAliases returns a mapping of cell alias to cells identified by that
	// alias.
	GetCellsAliases(ctx context.Context, in *vtctldata.GetCellsAliasesRequest, opts ...grpc.CallOption) (*vtctldata.GetCellsAliasesResponse, error)
	// GetKeyspace reads the given keyspace from the topo and returns it.
	GetKeyspace(ctx context.Context, in *vtctldata.GetKeyspaceRequest, opts ...grpc.CallOption) (*vtctldata.GetKeyspaceResponse, error)
	// GetKeyspaces returns the keyspace struct of all keyspaces in the topo.
	GetKeyspaces(ctx context.Context, in *vtctldata.GetKeyspacesRequest, opts ...grpc.CallOption) (*vtctldata.GetKeyspacesResponse, error)
	// GetSrvVSchema returns a the SrvVSchema for a cell.
	GetSrvVSchema(ctx context.Context, in *vtctldata.GetSrvVSchemaRequest, opts ...grpc.CallOption) (*vtctldata.GetSrvVSchemaResponse, error)
	// GetTablet returns information about a tablet.
	GetTablet(ctx context.Context, in *vtctldata.GetTabletRequest, opts ...grpc.CallOption) (*vtctldata.GetTabletResponse, error)
	// GetTablets returns tablets, optionally filtered by keyspace and shard.
	GetTablets(ctx context.Context, in *vtctldata.GetTabletsRequest, opts ...grpc.CallOption) (*vtctldata.GetTabletsResponse, error)
	// InitShardPrimary sets the initial primary for a shard. Will make all other
	// tablets in the shard replicas of the provided primary.
	//
	// WARNING: This could cause data loss on an already replicating shard.
	// PlannedReparentShard or EmergencyReparentShard should be used in those
	// cases instead.
	InitShardPrimary(ctx context.Context, in *vtctldata.InitShardPrimaryRequest, opts ...grpc.CallOption) (*vtctldata.InitShardPrimaryResponse, error)
}

type vtctldClient struct {
	cc *grpc.ClientConn
}

func NewVtctldClient(cc *grpc.ClientConn) VtctldClient {
	return &vtctldClient{cc}
}

func (c *vtctldClient) FindAllShardsInKeyspace(ctx context.Context, in *vtctldata.FindAllShardsInKeyspaceRequest, opts ...grpc.CallOption) (*vtctldata.FindAllShardsInKeyspaceResponse, error) {
	out := new(vtctldata.FindAllShardsInKeyspaceResponse)
	err := c.cc.Invoke(ctx, "/vtctlservice.Vtctld/FindAllShardsInKeyspace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vtctldClient) GetCellInfoNames(ctx context.Context, in *vtctldata.GetCellInfoNamesRequest, opts ...grpc.CallOption) (*vtctldata.GetCellInfoNamesResponse, error) {
	out := new(vtctldata.GetCellInfoNamesResponse)
	err := c.cc.Invoke(ctx, "/vtctlservice.Vtctld/GetCellInfoNames", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vtctldClient) GetCellInfo(ctx context.Context, in *vtctldata.GetCellInfoRequest, opts ...grpc.CallOption) (*vtctldata.GetCellInfoResponse, error) {
	out := new(vtctldata.GetCellInfoResponse)
	err := c.cc.Invoke(ctx, "/vtctlservice.Vtctld/GetCellInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vtctldClient) GetCellsAliases(ctx context.Context, in *vtctldata.GetCellsAliasesRequest, opts ...grpc.CallOption) (*vtctldata.GetCellsAliasesResponse, error) {
	out := new(vtctldata.GetCellsAliasesResponse)
	err := c.cc.Invoke(ctx, "/vtctlservice.Vtctld/GetCellsAliases", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vtctldClient) GetKeyspace(ctx context.Context, in *vtctldata.GetKeyspaceRequest, opts ...grpc.CallOption) (*vtctldata.GetKeyspaceResponse, error) {
	out := new(vtctldata.GetKeyspaceResponse)
	err := c.cc.Invoke(ctx, "/vtctlservice.Vtctld/GetKeyspace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vtctldClient) GetKeyspaces(ctx context.Context, in *vtctldata.GetKeyspacesRequest, opts ...grpc.CallOption) (*vtctldata.GetKeyspacesResponse, error) {
	out := new(vtctldata.GetKeyspacesResponse)
	err := c.cc.Invoke(ctx, "/vtctlservice.Vtctld/GetKeyspaces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vtctldClient) GetSrvVSchema(ctx context.Context, in *vtctldata.GetSrvVSchemaRequest, opts ...grpc.CallOption) (*vtctldata.GetSrvVSchemaResponse, error) {
	out := new(vtctldata.GetSrvVSchemaResponse)
	err := c.cc.Invoke(ctx, "/vtctlservice.Vtctld/GetSrvVSchema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vtctldClient) GetTablet(ctx context.Context, in *vtctldata.GetTabletRequest, opts ...grpc.CallOption) (*vtctldata.GetTabletResponse, error) {
	out := new(vtctldata.GetTabletResponse)
	err := c.cc.Invoke(ctx, "/vtctlservice.Vtctld/GetTablet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vtctldClient) GetTablets(ctx context.Context, in *vtctldata.GetTabletsRequest, opts ...grpc.CallOption) (*vtctldata.GetTabletsResponse, error) {
	out := new(vtctldata.GetTabletsResponse)
	err := c.cc.Invoke(ctx, "/vtctlservice.Vtctld/GetTablets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vtctldClient) InitShardPrimary(ctx context.Context, in *vtctldata.InitShardPrimaryRequest, opts ...grpc.CallOption) (*vtctldata.InitShardPrimaryResponse, error) {
	out := new(vtctldata.InitShardPrimaryResponse)
	err := c.cc.Invoke(ctx, "/vtctlservice.Vtctld/InitShardPrimary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VtctldServer is the server API for Vtctld service.
type VtctldServer interface {
	// FindAllShardsInKeyspace returns a map of shard names to shard references
	// for a given keyspace.
	FindAllShardsInKeyspace(context.Context, *vtctldata.FindAllShardsInKeyspaceRequest) (*vtctldata.FindAllShardsInKeyspaceResponse, error)
	// GetCellInfoNames returns all the cells for which we have a CellInfo object,
	// meaning we have a topology service registered.
	GetCellInfoNames(context.Context, *vtctldata.GetCellInfoNamesRequest) (*vtctldata.GetCellInfoNamesResponse, error)
	// GetCellInfo returns the information for a cell.
	GetCellInfo(context.Context, *vtctldata.GetCellInfoRequest) (*vtctldata.GetCellInfoResponse, error)
	// GetCellsAliases returns a mapping of cell alias to cells identified by that
	// alias.
	GetCellsAliases(context.Context, *vtctldata.GetCellsAliasesRequest) (*vtctldata.GetCellsAliasesResponse, error)
	// GetKeyspace reads the given keyspace from the topo and returns it.
	GetKeyspace(context.Context, *vtctldata.GetKeyspaceRequest) (*vtctldata.GetKeyspaceResponse, error)
	// GetKeyspaces returns the keyspace struct of all keyspaces in the topo.
	GetKeyspaces(context.Context, *vtctldata.GetKeyspacesRequest) (*vtctldata.GetKeyspacesResponse, error)
	// GetSrvVSchema returns a the SrvVSchema for a cell.
	GetSrvVSchema(context.Context, *vtctldata.GetSrvVSchemaRequest) (*vtctldata.GetSrvVSchemaResponse, error)
	// GetTablet returns information about a tablet.
	GetTablet(context.Context, *vtctldata.GetTabletRequest) (*vtctldata.GetTabletResponse, error)
	// GetTablets returns tablets, optionally filtered by keyspace and shard.
	GetTablets(context.Context, *vtctldata.GetTabletsRequest) (*vtctldata.GetTabletsResponse, error)
	// InitShardPrimary sets the initial primary for a shard. Will make all other
	// tablets in the shard replicas of the provided primary.
	//
	// WARNING: This could cause data loss on an already replicating shard.
	// PlannedReparentShard or EmergencyReparentShard should be used in those
	// cases instead.
	InitShardPrimary(context.Context, *vtctldata.InitShardPrimaryRequest) (*vtctldata.InitShardPrimaryResponse, error)
}

// UnimplementedVtctldServer can be embedded to have forward compatible implementations.
type UnimplementedVtctldServer struct {
}

func (*UnimplementedVtctldServer) FindAllShardsInKeyspace(ctx context.Context, req *vtctldata.FindAllShardsInKeyspaceRequest) (*vtctldata.FindAllShardsInKeyspaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAllShardsInKeyspace not implemented")
}
func (*UnimplementedVtctldServer) GetCellInfoNames(ctx context.Context, req *vtctldata.GetCellInfoNamesRequest) (*vtctldata.GetCellInfoNamesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCellInfoNames not implemented")
}
func (*UnimplementedVtctldServer) GetCellInfo(ctx context.Context, req *vtctldata.GetCellInfoRequest) (*vtctldata.GetCellInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCellInfo not implemented")
}
func (*UnimplementedVtctldServer) GetCellsAliases(ctx context.Context, req *vtctldata.GetCellsAliasesRequest) (*vtctldata.GetCellsAliasesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCellsAliases not implemented")
}
func (*UnimplementedVtctldServer) GetKeyspace(ctx context.Context, req *vtctldata.GetKeyspaceRequest) (*vtctldata.GetKeyspaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKeyspace not implemented")
}
func (*UnimplementedVtctldServer) GetKeyspaces(ctx context.Context, req *vtctldata.GetKeyspacesRequest) (*vtctldata.GetKeyspacesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKeyspaces not implemented")
}
func (*UnimplementedVtctldServer) GetSrvVSchema(ctx context.Context, req *vtctldata.GetSrvVSchemaRequest) (*vtctldata.GetSrvVSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSrvVSchema not implemented")
}
func (*UnimplementedVtctldServer) GetTablet(ctx context.Context, req *vtctldata.GetTabletRequest) (*vtctldata.GetTabletResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTablet not implemented")
}
func (*UnimplementedVtctldServer) GetTablets(ctx context.Context, req *vtctldata.GetTabletsRequest) (*vtctldata.GetTabletsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTablets not implemented")
}
func (*UnimplementedVtctldServer) InitShardPrimary(ctx context.Context, req *vtctldata.InitShardPrimaryRequest) (*vtctldata.InitShardPrimaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitShardPrimary not implemented")
}

func RegisterVtctldServer(s *grpc.Server, srv VtctldServer) {
	s.RegisterService(&_Vtctld_serviceDesc, srv)
}

func _Vtctld_FindAllShardsInKeyspace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vtctldata.FindAllShardsInKeyspaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VtctldServer).FindAllShardsInKeyspace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vtctlservice.Vtctld/FindAllShardsInKeyspace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VtctldServer).FindAllShardsInKeyspace(ctx, req.(*vtctldata.FindAllShardsInKeyspaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vtctld_GetCellInfoNames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vtctldata.GetCellInfoNamesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VtctldServer).GetCellInfoNames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vtctlservice.Vtctld/GetCellInfoNames",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VtctldServer).GetCellInfoNames(ctx, req.(*vtctldata.GetCellInfoNamesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vtctld_GetCellInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vtctldata.GetCellInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VtctldServer).GetCellInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vtctlservice.Vtctld/GetCellInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VtctldServer).GetCellInfo(ctx, req.(*vtctldata.GetCellInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vtctld_GetCellsAliases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vtctldata.GetCellsAliasesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VtctldServer).GetCellsAliases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vtctlservice.Vtctld/GetCellsAliases",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VtctldServer).GetCellsAliases(ctx, req.(*vtctldata.GetCellsAliasesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vtctld_GetKeyspace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vtctldata.GetKeyspaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VtctldServer).GetKeyspace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vtctlservice.Vtctld/GetKeyspace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VtctldServer).GetKeyspace(ctx, req.(*vtctldata.GetKeyspaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vtctld_GetKeyspaces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vtctldata.GetKeyspacesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VtctldServer).GetKeyspaces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vtctlservice.Vtctld/GetKeyspaces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VtctldServer).GetKeyspaces(ctx, req.(*vtctldata.GetKeyspacesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vtctld_GetSrvVSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vtctldata.GetSrvVSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VtctldServer).GetSrvVSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vtctlservice.Vtctld/GetSrvVSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VtctldServer).GetSrvVSchema(ctx, req.(*vtctldata.GetSrvVSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vtctld_GetTablet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vtctldata.GetTabletRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VtctldServer).GetTablet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vtctlservice.Vtctld/GetTablet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VtctldServer).GetTablet(ctx, req.(*vtctldata.GetTabletRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vtctld_GetTablets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vtctldata.GetTabletsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VtctldServer).GetTablets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vtctlservice.Vtctld/GetTablets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VtctldServer).GetTablets(ctx, req.(*vtctldata.GetTabletsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Vtctld_InitShardPrimary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(vtctldata.InitShardPrimaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VtctldServer).InitShardPrimary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vtctlservice.Vtctld/InitShardPrimary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VtctldServer).InitShardPrimary(ctx, req.(*vtctldata.InitShardPrimaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Vtctld_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vtctlservice.Vtctld",
	HandlerType: (*VtctldServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindAllShardsInKeyspace",
			Handler:    _Vtctld_FindAllShardsInKeyspace_Handler,
		},
		{
			MethodName: "GetCellInfoNames",
			Handler:    _Vtctld_GetCellInfoNames_Handler,
		},
		{
			MethodName: "GetCellInfo",
			Handler:    _Vtctld_GetCellInfo_Handler,
		},
		{
			MethodName: "GetCellsAliases",
			Handler:    _Vtctld_GetCellsAliases_Handler,
		},
		{
			MethodName: "GetKeyspace",
			Handler:    _Vtctld_GetKeyspace_Handler,
		},
		{
			MethodName: "GetKeyspaces",
			Handler:    _Vtctld_GetKeyspaces_Handler,
		},
		{
			MethodName: "GetSrvVSchema",
			Handler:    _Vtctld_GetSrvVSchema_Handler,
		},
		{
			MethodName: "GetTablet",
			Handler:    _Vtctld_GetTablet_Handler,
		},
		{
			MethodName: "GetTablets",
			Handler:    _Vtctld_GetTablets_Handler,
		},
		{
			MethodName: "InitShardPrimary",
			Handler:    _Vtctld_InitShardPrimary_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "vtctlservice.proto",
}
