// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gitlab.sz.sensetime.com/viper/resource-allocator/pb/resource-allocator.proto

/*
Package api is a generated protocol buffer package.

It is generated from these files:
	gitlab.sz.sensetime.com/viper/resource-allocator/pb/resource-allocator.proto

It has these top-level messages:
	NodeStatus
	NodeSpec
	ResourceStatus
	ResourceSpec
*/
package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type NodeStatus struct {
	Quotas map[string]uint64 `protobuf:"bytes,1,rep,name=quotas" json:"quotas,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
}

func (m *NodeStatus) Reset()                    { *m = NodeStatus{} }
func (m *NodeStatus) String() string            { return proto.CompactTextString(m) }
func (*NodeStatus) ProtoMessage()               {}
func (*NodeStatus) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *NodeStatus) GetQuotas() map[string]uint64 {
	if m != nil {
		return m.Quotas
	}
	return nil
}

type NodeSpec struct {
	ClusterId string            `protobuf:"bytes,1,opt,name=cluster_id,json=clusterId" json:"cluster_id,omitempty"`
	Id        string            `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
	Quotas    map[string]uint64 `protobuf:"bytes,3,rep,name=quotas" json:"quotas,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	//use to match Type, all the labels is match
	Labels    map[string]string `protobuf:"bytes,4,rep,name=labels" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Version   int32             `protobuf:"varint,5,opt,name=version" json:"version,omitempty"`
	//当前节点的计算能力状态
	Status    *NodeStatus       `protobuf:"bytes,6,opt,name=status" json:"status,omitempty"`
	Payload   []byte            `protobuf:"bytes,7,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *NodeSpec) Reset()                    { *m = NodeSpec{} }
func (m *NodeSpec) String() string            { return proto.CompactTextString(m) }
func (*NodeSpec) ProtoMessage()               {}
func (*NodeSpec) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *NodeSpec) GetClusterId() string {
	if m != nil {
		return m.ClusterId
	}
	return ""
}

func (m *NodeSpec) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *NodeSpec) GetQuotas() map[string]uint64 {
	if m != nil {
		return m.Quotas
	}
	return nil
}

func (m *NodeSpec) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *NodeSpec) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *NodeSpec) GetStatus() *NodeStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *NodeSpec) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type ResourceStatus struct {
	Version int64  `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *ResourceStatus) Reset()                    { *m = ResourceStatus{} }
func (m *ResourceStatus) String() string            { return proto.CompactTextString(m) }
func (*ResourceStatus) ProtoMessage()               {}
func (*ResourceStatus) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ResourceStatus) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *ResourceStatus) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

type ResourceSpec struct {
	Id             string            `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Limits         map[string]uint64 `protobuf:"bytes,2,rep,name=limits" json:"limits,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	RequiredLabels map[string]string `protobuf:"bytes,3,rep,name=required_labels,json=requiredLabels" json:"required_labels,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	AssignedNode   string            `protobuf:"bytes,4,opt,name=assigned_node,json=assignedNode" json:"assigned_node,omitempty"`
	Version        int32             `protobuf:"varint,5,opt,name=version" json:"version,omitempty"`
	Payload        []byte            `protobuf:"bytes,6,opt,name=payload,proto3" json:"payload,omitempty"`
	Status         *ResourceStatus   `protobuf:"bytes,7,opt,name=status" json:"status,omitempty"`
	CreationTime   int64             `protobuf:"varint,8,opt,name=creation_time,json=creationTime" json:"creation_time,omitempty"`
}

func (m *ResourceSpec) Reset()                    { *m = ResourceSpec{} }
func (m *ResourceSpec) String() string            { return proto.CompactTextString(m) }
func (*ResourceSpec) ProtoMessage()               {}
func (*ResourceSpec) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ResourceSpec) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ResourceSpec) GetLimits() map[string]uint64 {
	if m != nil {
		return m.Limits
	}
	return nil
}

func (m *ResourceSpec) GetRequiredLabels() map[string]string {
	if m != nil {
		return m.RequiredLabels
	}
	return nil
}

func (m *ResourceSpec) GetAssignedNode() string {
	if m != nil {
		return m.AssignedNode
	}
	return ""
}

func (m *ResourceSpec) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *ResourceSpec) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *ResourceSpec) GetStatus() *ResourceStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *ResourceSpec) GetCreationTime() int64 {
	if m != nil {
		return m.CreationTime
	}
	return 0
}

func init() {
	proto.RegisterType((*NodeStatus)(nil), "sensetime.viper.resource_allocator.NodeStatus")
	proto.RegisterType((*NodeSpec)(nil), "sensetime.viper.resource_allocator.NodeSpec")
	proto.RegisterType((*ResourceStatus)(nil), "sensetime.viper.resource_allocator.ResourceStatus")
	proto.RegisterType((*ResourceSpec)(nil), "sensetime.viper.resource_allocator.ResourceSpec")
}

func init() {
	proto.RegisterFile("gitlab.sz.sensetime.com/viper/resource-allocator/pb/resource-allocator.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 516 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0xd1, 0x8a, 0xd3, 0x40,
	0x14, 0x86, 0x49, 0xda, 0x4d, 0xb7, 0xa7, 0xdd, 0x2a, 0xa3, 0x17, 0x83, 0x20, 0x84, 0x2e, 0x48,
	0x6e, 0x4c, 0xa1, 0x7a, 0xb1, 0x2e, 0xde, 0x28, 0xab, 0xa0, 0x14, 0xd1, 0xd9, 0xbd, 0xf2, 0xc2,
	0x32, 0x4d, 0x0e, 0xcb, 0x60, 0x92, 0xc9, 0xce, 0x4c, 0x0a, 0xf5, 0x55, 0x7c, 0x04, 0x1f, 0xc6,
	0x57, 0x92, 0x4c, 0x93, 0x34, 0x85, 0xe2, 0x36, 0x7a, 0xd7, 0x73, 0xc8, 0xff, 0xf5, 0x9c, 0xf3,
	0xff, 0x0c, 0x2c, 0x6e, 0x85, 0x49, 0xf8, 0x2a, 0xd4, 0x3f, 0x42, 0x8d, 0x99, 0x46, 0x23, 0x52,
	0x0c, 0x23, 0x99, 0xce, 0xd6, 0x22, 0x47, 0x35, 0x53, 0xa8, 0x65, 0xa1, 0x22, 0x7c, 0xce, 0x93,
	0x44, 0x46, 0xdc, 0x48, 0x35, 0xcb, 0x57, 0x07, 0xba, 0x61, 0xae, 0xa4, 0x91, 0x64, 0xba, 0x63,
	0x58, 0x7d, 0x58, 0x7f, 0xb9, 0x6c, 0xbe, 0x9c, 0xfe, 0x74, 0x00, 0x3e, 0xc9, 0x18, 0xaf, 0x0d,
	0x37, 0x85, 0x26, 0x0c, 0xbc, 0xbb, 0x42, 0x1a, 0xae, 0xa9, 0xe3, 0xf7, 0x82, 0xd1, 0xfc, 0x32,
	0xbc, 0x9f, 0x11, 0xee, 0xf4, 0xe1, 0x17, 0x2b, 0x7e, 0x97, 0x19, 0xb5, 0x61, 0x15, 0xe9, 0xc9,
	0x2b, 0x18, 0xb5, 0xda, 0xe4, 0x21, 0xf4, 0xbe, 0xe3, 0x86, 0x3a, 0xbe, 0x13, 0x0c, 0x59, 0xf9,
	0x93, 0x3c, 0x86, 0x93, 0x35, 0x4f, 0x0a, 0xa4, 0xae, 0xef, 0x04, 0x7d, 0xb6, 0x2d, 0x2e, 0xdd,
	0x0b, 0x67, 0xfa, 0xbb, 0x07, 0xa7, 0x96, 0x9e, 0x63, 0x44, 0x9e, 0x02, 0x44, 0x49, 0xa1, 0x0d,
	0xaa, 0xa5, 0x88, 0x2b, 0xfd, 0xb0, 0xea, 0x7c, 0x88, 0xc9, 0x04, 0x5c, 0x11, 0x5b, 0xc4, 0x90,
	0xb9, 0x22, 0x26, 0x9f, 0x9b, 0x55, 0x7a, 0x76, 0x95, 0x8b, 0xa3, 0x57, 0xc9, 0x31, 0x3a, 0xb4,
	0x48, 0x49, 0x4c, 0xf8, 0x0a, 0x13, 0x4d, 0xfb, 0xff, 0x40, 0x5c, 0x58, 0x69, 0x45, 0xdc, 0x72,
	0x08, 0x85, 0xc1, 0x1a, 0x95, 0x16, 0x32, 0xa3, 0x27, 0xbe, 0x13, 0x9c, 0xb0, 0xba, 0x24, 0xef,
	0xc1, 0xd3, 0xf6, 0xa4, 0xd4, 0xf3, 0x9d, 0x60, 0x34, 0x0f, 0xbb, 0x19, 0xc1, 0x2a, 0x75, 0xf9,
	0x0f, 0x39, 0xdf, 0x24, 0x92, 0xc7, 0x74, 0xe0, 0x3b, 0xc1, 0x98, 0xd5, 0xe5, 0x7f, 0xd8, 0x52,
	0x4a, 0x5b, 0xdb, 0xdc, 0x27, 0x1d, 0xb6, 0x1d, 0xbd, 0x82, 0x09, 0xab, 0x06, 0xbf, 0x6e, 0x26,
	0xac, 0x6f, 0x50, 0x12, 0x7a, 0xbb, 0x1b, 0xb4, 0x66, 0x77, 0xf7, 0x66, 0x9f, 0xfe, 0xea, 0xc3,
	0xb8, 0xc1, 0x94, 0xd9, 0xd8, 0x9a, 0xef, 0x34, 0xe6, 0xdf, 0x80, 0x97, 0x88, 0x54, 0x18, 0x4d,
	0x5d, 0x6b, 0xd5, 0xeb, 0x63, 0xce, 0xd7, 0x26, 0x86, 0x0b, 0x2b, 0xaf, 0xed, 0xb2, 0x05, 0x49,
	0xe1, 0x81, 0xc2, 0xbb, 0x42, 0x28, 0x8c, 0x97, 0x55, 0x12, 0xb6, 0xd9, 0xba, 0xea, 0x8c, 0x67,
	0x15, 0xa7, 0x9d, 0x8a, 0x89, 0xda, 0x6b, 0x92, 0x73, 0x38, 0xe3, 0x5a, 0x8b, 0xdb, 0x0c, 0xe3,
	0x65, 0x26, 0x63, 0xa4, 0x7d, 0xbb, 0xdf, 0xb8, 0x6e, 0x96, 0x76, 0xff, 0x25, 0x42, 0xad, 0xf3,
	0x79, 0x7b, 0xe7, 0x23, 0x1f, 0x9b, 0x70, 0x0d, 0x6c, 0xb8, 0xe6, 0x9d, 0xc6, 0xdf, 0x0f, 0xd8,
	0x39, 0x9c, 0x45, 0x0a, 0xb9, 0x11, 0x32, 0x5b, 0x96, 0x00, 0x7a, 0x6a, 0x4d, 0x1c, 0xd7, 0xcd,
	0x1b, 0x91, 0xa2, 0x0d, 0xcc, 0xee, 0x9e, 0x9d, 0xb2, 0xf6, 0x06, 0x1e, 0x1d, 0xb8, 0x55, 0x97,
	0xcc, 0xbd, 0xfd, 0x06, 0xcf, 0x22, 0x99, 0x1e, 0xb1, 0xe3, 0xd7, 0x97, 0x9d, 0xdf, 0x5f, 0x9e,
	0x8b, 0x95, 0x67, 0x9f, 0xdb, 0x17, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x0f, 0x4f, 0x01, 0x76,
	0xbe, 0x05, 0x00, 0x00,
}
