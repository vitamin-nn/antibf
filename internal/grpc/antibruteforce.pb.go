// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.3
// source: antibruteforce.proto

package grpc

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type CheckAuthRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login    string `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Ip       string `protobuf:"bytes,3,opt,name=ip,proto3" json:"ip,omitempty"`
}

func (x *CheckAuthRequest) Reset() {
	*x = CheckAuthRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_antibruteforce_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckAuthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckAuthRequest) ProtoMessage() {}

func (x *CheckAuthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_antibruteforce_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckAuthRequest.ProtoReflect.Descriptor instead.
func (*CheckAuthRequest) Descriptor() ([]byte, []int) {
	return file_antibruteforce_proto_rawDescGZIP(), []int{0}
}

func (x *CheckAuthRequest) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *CheckAuthRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *CheckAuthRequest) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

type CheckAuthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Result:
	//	*CheckAuthResponse_Ok
	//	*CheckAuthResponse_Error
	Result isCheckAuthResponse_Result `protobuf_oneof:"result"`
}

func (x *CheckAuthResponse) Reset() {
	*x = CheckAuthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_antibruteforce_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckAuthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckAuthResponse) ProtoMessage() {}

func (x *CheckAuthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_antibruteforce_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckAuthResponse.ProtoReflect.Descriptor instead.
func (*CheckAuthResponse) Descriptor() ([]byte, []int) {
	return file_antibruteforce_proto_rawDescGZIP(), []int{1}
}

func (m *CheckAuthResponse) GetResult() isCheckAuthResponse_Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (x *CheckAuthResponse) GetOk() bool {
	if x, ok := x.GetResult().(*CheckAuthResponse_Ok); ok {
		return x.Ok
	}
	return false
}

func (x *CheckAuthResponse) GetError() string {
	if x, ok := x.GetResult().(*CheckAuthResponse_Error); ok {
		return x.Error
	}
	return ""
}

type isCheckAuthResponse_Result interface {
	isCheckAuthResponse_Result()
}

type CheckAuthResponse_Ok struct {
	Ok bool `protobuf:"varint,1,opt,name=ok,proto3,oneof"`
}

type CheckAuthResponse_Error struct {
	Error string `protobuf:"bytes,2,opt,name=error,proto3,oneof"`
}

func (*CheckAuthResponse_Ok) isCheckAuthResponse_Result() {}

func (*CheckAuthResponse_Error) isCheckAuthResponse_Result() {}

var File_antibruteforce_proto protoreflect.FileDescriptor

var file_antibruteforce_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x65, 0x66, 0x6f, 0x72, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x54, 0x0a, 0x10, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x41,
	0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x6f,
	0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e,
	0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x22, 0x47, 0x0a, 0x11,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x10, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52,
	0x02, 0x6f, 0x6b, 0x12, 0x16, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x08, 0x0a, 0x06, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x32, 0x50, 0x0a, 0x15, 0x41, 0x6e, 0x74, 0x69, 0x42, 0x72, 0x75,
	0x74, 0x65, 0x66, 0x6f, 0x72, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x37,
	0x0a, 0x0c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x11,
	0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x12, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x3b, 0x67, 0x72, 0x70,
	0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_antibruteforce_proto_rawDescOnce sync.Once
	file_antibruteforce_proto_rawDescData = file_antibruteforce_proto_rawDesc
)

func file_antibruteforce_proto_rawDescGZIP() []byte {
	file_antibruteforce_proto_rawDescOnce.Do(func() {
		file_antibruteforce_proto_rawDescData = protoimpl.X.CompressGZIP(file_antibruteforce_proto_rawDescData)
	})
	return file_antibruteforce_proto_rawDescData
}

var file_antibruteforce_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_antibruteforce_proto_goTypes = []interface{}{
	(*CheckAuthRequest)(nil),  // 0: CheckAuthRequest
	(*CheckAuthResponse)(nil), // 1: CheckAuthResponse
}
var file_antibruteforce_proto_depIdxs = []int32{
	0, // 0: AntiBruteforceService.CheckRequest:input_type -> CheckAuthRequest
	1, // 1: AntiBruteforceService.CheckRequest:output_type -> CheckAuthResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_antibruteforce_proto_init() }
func file_antibruteforce_proto_init() {
	if File_antibruteforce_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_antibruteforce_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckAuthRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_antibruteforce_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckAuthResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_antibruteforce_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*CheckAuthResponse_Ok)(nil),
		(*CheckAuthResponse_Error)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_antibruteforce_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_antibruteforce_proto_goTypes,
		DependencyIndexes: file_antibruteforce_proto_depIdxs,
		MessageInfos:      file_antibruteforce_proto_msgTypes,
	}.Build()
	File_antibruteforce_proto = out.File
	file_antibruteforce_proto_rawDesc = nil
	file_antibruteforce_proto_goTypes = nil
	file_antibruteforce_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AntiBruteforceServiceClient is the client API for AntiBruteforceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AntiBruteforceServiceClient interface {
	CheckRequest(ctx context.Context, in *CheckAuthRequest, opts ...grpc.CallOption) (*CheckAuthResponse, error)
}

type antiBruteforceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAntiBruteforceServiceClient(cc grpc.ClientConnInterface) AntiBruteforceServiceClient {
	return &antiBruteforceServiceClient{cc}
}

func (c *antiBruteforceServiceClient) CheckRequest(ctx context.Context, in *CheckAuthRequest, opts ...grpc.CallOption) (*CheckAuthResponse, error) {
	out := new(CheckAuthResponse)
	err := c.cc.Invoke(ctx, "/AntiBruteforceService/CheckRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AntiBruteforceServiceServer is the server API for AntiBruteforceService service.
type AntiBruteforceServiceServer interface {
	CheckRequest(context.Context, *CheckAuthRequest) (*CheckAuthResponse, error)
}

// UnimplementedAntiBruteforceServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAntiBruteforceServiceServer struct {
}

func (*UnimplementedAntiBruteforceServiceServer) CheckRequest(context.Context, *CheckAuthRequest) (*CheckAuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckRequest not implemented")
}

func RegisterAntiBruteforceServiceServer(s *grpc.Server, srv AntiBruteforceServiceServer) {
	s.RegisterService(&_AntiBruteforceService_serviceDesc, srv)
}

func _AntiBruteforceService_CheckRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckAuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AntiBruteforceServiceServer).CheckRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AntiBruteforceService/CheckRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AntiBruteforceServiceServer).CheckRequest(ctx, req.(*CheckAuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AntiBruteforceService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "AntiBruteforceService",
	HandlerType: (*AntiBruteforceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckRequest",
			Handler:    _AntiBruteforceService_CheckRequest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "antibruteforce.proto",
}