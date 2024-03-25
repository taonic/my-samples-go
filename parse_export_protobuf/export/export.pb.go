// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.3
// source: exported_workflow.proto

package export

import (
	v1 "go.temporal.io/api/history/v1"
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

type ExportedWorkflows struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Workflows []*Workflow `protobuf:"bytes,1,rep,name=workflows,proto3" json:"workflows,omitempty"`
}

func (x *ExportedWorkflows) Reset() {
	*x = ExportedWorkflows{}
	if protoimpl.UnsafeEnabled {
		mi := &file_exported_workflow_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExportedWorkflows) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExportedWorkflows) ProtoMessage() {}

func (x *ExportedWorkflows) ProtoReflect() protoreflect.Message {
	mi := &file_exported_workflow_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExportedWorkflows.ProtoReflect.Descriptor instead.
func (*ExportedWorkflows) Descriptor() ([]byte, []int) {
	return file_exported_workflow_proto_rawDescGZIP(), []int{0}
}

func (x *ExportedWorkflows) GetWorkflows() []*Workflow {
	if x != nil {
		return x.Workflows
	}
	return nil
}

type Workflow struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	History *v1.History `protobuf:"bytes,1,opt,name=history,proto3" json:"history,omitempty"`
}

func (x *Workflow) Reset() {
	*x = Workflow{}
	if protoimpl.UnsafeEnabled {
		mi := &file_exported_workflow_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Workflow) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Workflow) ProtoMessage() {}

func (x *Workflow) ProtoReflect() protoreflect.Message {
	mi := &file_exported_workflow_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Workflow.ProtoReflect.Descriptor instead.
func (*Workflow) Descriptor() ([]byte, []int) {
	return file_exported_workflow_proto_rawDescGZIP(), []int{1}
}

func (x *Workflow) GetHistory() *v1.History {
	if x != nil {
		return x.History
	}
	return nil
}

var File_exported_workflow_proto protoreflect.FileDescriptor

var file_exported_workflow_proto_rawDesc = []byte{
	0x0a, 0x17, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x5f, 0x77, 0x6f, 0x72, 0x6b, 0x66,
	0x6c, 0x6f, 0x77, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x65, 0x78, 0x70, 0x6f, 0x72,
	0x74, 0x2e, 0x76, 0x31, 0x1a, 0x25, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x46, 0x0a, 0x11, 0x45,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73,
	0x12, 0x31, 0x0a, 0x09, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x76, 0x31, 0x2e,
	0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x52, 0x09, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c,
	0x6f, 0x77, 0x73, 0x22, 0x46, 0x0a, 0x08, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x12,
	0x3a, 0x0a, 0x07, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x20, 0x2e, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x69, 0x73, 0x74, 0x6f,
	0x72, 0x79, 0x52, 0x07, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x42, 0x20, 0x5a, 0x1e, 0x67,
	0x6f, 0x2e, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2e, 0x69, 0x6f, 0x2f, 0x74, 0x65,
	0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x5f, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_exported_workflow_proto_rawDescOnce sync.Once
	file_exported_workflow_proto_rawDescData = file_exported_workflow_proto_rawDesc
)

func file_exported_workflow_proto_rawDescGZIP() []byte {
	file_exported_workflow_proto_rawDescOnce.Do(func() {
		file_exported_workflow_proto_rawDescData = protoimpl.X.CompressGZIP(file_exported_workflow_proto_rawDescData)
	})
	return file_exported_workflow_proto_rawDescData
}

var file_exported_workflow_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_exported_workflow_proto_goTypes = []interface{}{
	(*ExportedWorkflows)(nil), // 0: export.v1.ExportedWorkflows
	(*Workflow)(nil),          // 1: export.v1.Workflow
	(*v1.History)(nil),        // 2: temporal.api.history.v1.History
}
var file_exported_workflow_proto_depIdxs = []int32{
	1, // 0: export.v1.ExportedWorkflows.workflows:type_name -> export.v1.Workflow
	2, // 1: export.v1.Workflow.history:type_name -> temporal.api.history.v1.History
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_exported_workflow_proto_init() }
func file_exported_workflow_proto_init() {
	if File_exported_workflow_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_exported_workflow_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExportedWorkflows); i {
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
		file_exported_workflow_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Workflow); i {
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
			RawDescriptor: file_exported_workflow_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_exported_workflow_proto_goTypes,
		DependencyIndexes: file_exported_workflow_proto_depIdxs,
		MessageInfos:      file_exported_workflow_proto_msgTypes,
	}.Build()
	File_exported_workflow_proto = out.File
	file_exported_workflow_proto_rawDesc = nil
	file_exported_workflow_proto_goTypes = nil
	file_exported_workflow_proto_depIdxs = nil
}