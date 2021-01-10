// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.12.3
// source: service_upload.proto

package services

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

//
type UploadSubmission struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Val1 string `protobuf:"bytes,1,opt,name=val1,proto3" json:"val1,omitempty"`
}

func (x *UploadSubmission) Reset() {
	*x = UploadSubmission{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_upload_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadSubmission) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadSubmission) ProtoMessage() {}

func (x *UploadSubmission) ProtoReflect() protoreflect.Message {
	mi := &file_service_upload_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadSubmission.ProtoReflect.Descriptor instead.
func (*UploadSubmission) Descriptor() ([]byte, []int) {
	return file_service_upload_proto_rawDescGZIP(), []int{0}
}

func (x *UploadSubmission) GetVal1() string {
	if x != nil {
		return x.Val1
	}
	return ""
}

//
type UploadSummary struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Val2 string `protobuf:"bytes,1,opt,name=val2,proto3" json:"val2,omitempty"`
}

func (x *UploadSummary) Reset() {
	*x = UploadSummary{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_upload_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadSummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadSummary) ProtoMessage() {}

func (x *UploadSummary) ProtoReflect() protoreflect.Message {
	mi := &file_service_upload_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadSummary.ProtoReflect.Descriptor instead.
func (*UploadSummary) Descriptor() ([]byte, []int) {
	return file_service_upload_proto_rawDescGZIP(), []int{1}
}

func (x *UploadSummary) GetVal2() string {
	if x != nil {
		return x.Val2
	}
	return ""
}

var File_service_upload_proto protoreflect.FileDescriptor

var file_service_upload_proto_rawDesc = []byte{
	0x0a, 0x14, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x22, 0x26, 0x0a, 0x10, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x76, 0x61, 0x6c, 0x31, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x76, 0x61, 0x6c, 0x31, 0x22, 0x23, 0x0a, 0x0d, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x76, 0x61, 0x6c,
	0x32, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x76, 0x61, 0x6c, 0x32, 0x32, 0x4c, 0x0a,
	0x06, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x42, 0x0a, 0x09, 0x52, 0x75, 0x6e, 0x55, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x12, 0x1a, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e,
	0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x1a, 0x17, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x22, 0x00, 0x42, 0x0c, 0x5a, 0x0a, 0x2e,
	0x3b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_service_upload_proto_rawDescOnce sync.Once
	file_service_upload_proto_rawDescData = file_service_upload_proto_rawDesc
)

func file_service_upload_proto_rawDescGZIP() []byte {
	file_service_upload_proto_rawDescOnce.Do(func() {
		file_service_upload_proto_rawDescData = protoimpl.X.CompressGZIP(file_service_upload_proto_rawDescData)
	})
	return file_service_upload_proto_rawDescData
}

var file_service_upload_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_service_upload_proto_goTypes = []interface{}{
	(*UploadSubmission)(nil), // 0: services.UploadSubmission
	(*UploadSummary)(nil),    // 1: services.UploadSummary
}
var file_service_upload_proto_depIdxs = []int32{
	0, // 0: services.Upload.RunUpload:input_type -> services.UploadSubmission
	1, // 1: services.Upload.RunUpload:output_type -> services.UploadSummary
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_service_upload_proto_init() }
func file_service_upload_proto_init() {
	if File_service_upload_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_service_upload_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadSubmission); i {
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
		file_service_upload_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadSummary); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_service_upload_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_upload_proto_goTypes,
		DependencyIndexes: file_service_upload_proto_depIdxs,
		MessageInfos:      file_service_upload_proto_msgTypes,
	}.Build()
	File_service_upload_proto = out.File
	file_service_upload_proto_rawDesc = nil
	file_service_upload_proto_goTypes = nil
	file_service_upload_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// UploadClient is the client API for Upload service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UploadClient interface {
	// RunUpload is an exposed function for the Upload service
	RunUpload(ctx context.Context, in *UploadSubmission, opts ...grpc.CallOption) (*UploadSummary, error)
}

type uploadClient struct {
	cc grpc.ClientConnInterface
}

func NewUploadClient(cc grpc.ClientConnInterface) UploadClient {
	return &uploadClient{cc}
}

func (c *uploadClient) RunUpload(ctx context.Context, in *UploadSubmission, opts ...grpc.CallOption) (*UploadSummary, error) {
	out := new(UploadSummary)
	err := c.cc.Invoke(ctx, "/services.Upload/RunUpload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UploadServer is the server API for Upload service.
type UploadServer interface {
	// RunUpload is an exposed function for the Upload service
	RunUpload(context.Context, *UploadSubmission) (*UploadSummary, error)
}

// UnimplementedUploadServer can be embedded to have forward compatible implementations.
type UnimplementedUploadServer struct {
}

func (*UnimplementedUploadServer) RunUpload(context.Context, *UploadSubmission) (*UploadSummary, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunUpload not implemented")
}

func RegisterUploadServer(s *grpc.Server, srv UploadServer) {
	s.RegisterService(&_Upload_serviceDesc, srv)
}

func _Upload_RunUpload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadSubmission)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UploadServer).RunUpload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.Upload/RunUpload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UploadServer).RunUpload(ctx, req.(*UploadSubmission))
	}
	return interceptor(ctx, in, info, handler)
}

var _Upload_serviceDesc = grpc.ServiceDesc{
	ServiceName: "services.Upload",
	HandlerType: (*UploadServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RunUpload",
			Handler:    _Upload_RunUpload_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_upload.proto",
}
