// Code generated by protoc-gen-go-pulsar. DO NOT EDIT.
package auctionv1

import (
	_ "cosmossdk.io/api/cosmos/base/v1beta1"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.0
// 	protoc        (unknown)
// source: fatal_fruit/auction/v1/event.proto

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_fatal_fruit_auction_v1_event_proto protoreflect.FileDescriptor

var file_fatal_fruit_auction_v1_event_proto_rawDesc = []byte{
	0x0a, 0x22, 0x66, 0x61, 0x74, 0x61, 0x6c, 0x5f, 0x66, 0x72, 0x75, 0x69, 0x74, 0x2f, 0x61, 0x75,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x66, 0x61, 0x74, 0x61, 0x6c, 0x5f, 0x66, 0x72, 0x75, 0x69,
	0x74, 0x2e, 0x61, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x19, 0x63, 0x6f,
	0x73, 0x6d, 0x6f, 0x73, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x73, 0x6d, 0x6f,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2f,
	0x6d, 0x73, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x14, 0x67, 0x6f, 0x67, 0x6f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x67, 0x6f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2f, 0x62,
	0x61, 0x73, 0x65, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x63, 0x6f, 0x69, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x61, 0x6d, 0x69, 0x6e, 0x6f, 0x2f, 0x61, 0x6d,
	0x69, 0x6e, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x42, 0xe3, 0x01, 0x0a, 0x1a, 0x63, 0x6f,
	0x6d, 0x2e, 0x66, 0x61, 0x74, 0x61, 0x6c, 0x5f, 0x66, 0x72, 0x75, 0x69, 0x74, 0x2e, 0x61, 0x75,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x43, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x66, 0x61, 0x74, 0x61, 0x6c, 0x2d, 0x66, 0x72, 0x75, 0x69, 0x74, 0x2f, 0x61,
	0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x61, 0x74, 0x61, 0x6c,
	0x5f, 0x66, 0x72, 0x75, 0x69, 0x74, 0x2f, 0x61, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76,
	0x31, 0x3b, 0x61, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x46, 0x41,
	0x58, 0xaa, 0x02, 0x15, 0x46, 0x61, 0x74, 0x61, 0x6c, 0x46, 0x72, 0x75, 0x69, 0x74, 0x2e, 0x41,
	0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x15, 0x46, 0x61, 0x74, 0x61,
	0x6c, 0x46, 0x72, 0x75, 0x69, 0x74, 0x5c, 0x41, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x56,
	0x31, 0xe2, 0x02, 0x21, 0x46, 0x61, 0x74, 0x61, 0x6c, 0x46, 0x72, 0x75, 0x69, 0x74, 0x5c, 0x41,
	0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x17, 0x46, 0x61, 0x74, 0x61, 0x6c, 0x46, 0x72, 0x75,
	0x69, 0x74, 0x3a, 0x3a, 0x41, 0x75, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x3a, 0x56, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_fatal_fruit_auction_v1_event_proto_goTypes = []interface{}{}
var file_fatal_fruit_auction_v1_event_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_fatal_fruit_auction_v1_event_proto_init() }
func file_fatal_fruit_auction_v1_event_proto_init() {
	if File_fatal_fruit_auction_v1_event_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_fatal_fruit_auction_v1_event_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_fatal_fruit_auction_v1_event_proto_goTypes,
		DependencyIndexes: file_fatal_fruit_auction_v1_event_proto_depIdxs,
	}.Build()
	File_fatal_fruit_auction_v1_event_proto = out.File
	file_fatal_fruit_auction_v1_event_proto_rawDesc = nil
	file_fatal_fruit_auction_v1_event_proto_goTypes = nil
	file_fatal_fruit_auction_v1_event_proto_depIdxs = nil
}
