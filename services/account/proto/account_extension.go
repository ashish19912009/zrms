package proto

import (
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Extension field numbers from your proto
const (
	Field_Resource = 50001
	Field_Action   = 50002
)

var (
	// E_Resource is a protobuf extension for MethodOptions.resource
	E_Resource = &protoimpl.ExtensionInfo{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         Field_Resource,
		Name:          "account.resource",
		Tag:           "bytes,50001,opt,name=resource",
		Filename:      "account.proto",
	}

	// E_Action is a protobuf extension for MethodOptions.action
	E_Action = &protoimpl.ExtensionInfo{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         Field_Action,
		Name:          "account.action",
		Tag:           "bytes,50002,opt,name=action",
		Filename:      "account.proto",
	}
)
