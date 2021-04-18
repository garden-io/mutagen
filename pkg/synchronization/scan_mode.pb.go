// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: synchronization/scan_mode.proto

package synchronization

import (
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

// ScanMode specifies the mode for synchronization root scanning.
type ScanMode int32

const (
	// ScanMode_ScanModeDefault represents an unspecified scan mode. It should
	// be converted to one of the following values based on the desired default
	// behavior.
	ScanMode_ScanModeDefault ScanMode = 0
	// ScanMode_ScanModeFull specifies that full scans should be performed on
	// each synchronization cycle.
	ScanMode_ScanModeFull ScanMode = 1
	// ScanMode_ScanModeAccelerated specifies that scans should attempt to use
	// watch-based acceleration.
	ScanMode_ScanModeAccelerated ScanMode = 2
)

// Enum value maps for ScanMode.
var (
	ScanMode_name = map[int32]string{
		0: "ScanModeDefault",
		1: "ScanModeFull",
		2: "ScanModeAccelerated",
	}
	ScanMode_value = map[string]int32{
		"ScanModeDefault":     0,
		"ScanModeFull":        1,
		"ScanModeAccelerated": 2,
	}
)

func (x ScanMode) Enum() *ScanMode {
	p := new(ScanMode)
	*p = x
	return p
}

func (x ScanMode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ScanMode) Descriptor() protoreflect.EnumDescriptor {
	return file_synchronization_scan_mode_proto_enumTypes[0].Descriptor()
}

func (ScanMode) Type() protoreflect.EnumType {
	return &file_synchronization_scan_mode_proto_enumTypes[0]
}

func (x ScanMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ScanMode.Descriptor instead.
func (ScanMode) EnumDescriptor() ([]byte, []int) {
	return file_synchronization_scan_mode_proto_rawDescGZIP(), []int{0}
}

var File_synchronization_scan_mode_proto protoreflect.FileDescriptor

var file_synchronization_scan_mode_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x73, 0x79, 0x6e, 0x63, 0x68, 0x72, 0x6f, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2f, 0x73, 0x63, 0x61, 0x6e, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0f, 0x73, 0x79, 0x6e, 0x63, 0x68, 0x72, 0x6f, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2a, 0x4a, 0x0a, 0x08, 0x53, 0x63, 0x61, 0x6e, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x13,
	0x0a, 0x0f, 0x53, 0x63, 0x61, 0x6e, 0x4d, 0x6f, 0x64, 0x65, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c,
	0x74, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x53, 0x63, 0x61, 0x6e, 0x4d, 0x6f, 0x64, 0x65, 0x46,
	0x75, 0x6c, 0x6c, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x53, 0x63, 0x61, 0x6e, 0x4d, 0x6f, 0x64,
	0x65, 0x41, 0x63, 0x63, 0x65, 0x6c, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x10, 0x02, 0x42, 0x33,
	0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x75, 0x74,
	0x61, 0x67, 0x65, 0x6e, 0x2d, 0x69, 0x6f, 0x2f, 0x6d, 0x75, 0x74, 0x61, 0x67, 0x65, 0x6e, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x68, 0x72, 0x6f, 0x6e, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_synchronization_scan_mode_proto_rawDescOnce sync.Once
	file_synchronization_scan_mode_proto_rawDescData = file_synchronization_scan_mode_proto_rawDesc
)

func file_synchronization_scan_mode_proto_rawDescGZIP() []byte {
	file_synchronization_scan_mode_proto_rawDescOnce.Do(func() {
		file_synchronization_scan_mode_proto_rawDescData = protoimpl.X.CompressGZIP(file_synchronization_scan_mode_proto_rawDescData)
	})
	return file_synchronization_scan_mode_proto_rawDescData
}

var file_synchronization_scan_mode_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_synchronization_scan_mode_proto_goTypes = []interface{}{
	(ScanMode)(0), // 0: synchronization.ScanMode
}
var file_synchronization_scan_mode_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_synchronization_scan_mode_proto_init() }
func file_synchronization_scan_mode_proto_init() {
	if File_synchronization_scan_mode_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_synchronization_scan_mode_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_synchronization_scan_mode_proto_goTypes,
		DependencyIndexes: file_synchronization_scan_mode_proto_depIdxs,
		EnumInfos:         file_synchronization_scan_mode_proto_enumTypes,
	}.Build()
	File_synchronization_scan_mode_proto = out.File
	file_synchronization_scan_mode_proto_rawDesc = nil
	file_synchronization_scan_mode_proto_goTypes = nil
	file_synchronization_scan_mode_proto_depIdxs = nil
}
