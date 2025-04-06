//
// Copyright (c) 2025 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: admin/v1/hub_type.proto

package adminv1

import (
	v1 "github.com/innabox/fulfillment-cli/internal/api/shared/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Represents the overall state of a hub.
type HubState int32

const (
	// Unspecified indicates that the state isn't set.
	HubState_HUB_STATE_UNSPECIFIED HubState = 0
	// Unspecified indicates that the hub is ready.
	HubState_HUB_STATE_READY HubState = 1
	// Indicates that hub has been disabled by the administrator.
	HubState_HUB_STATE_DISABLED HubState = 2
)

// Enum value maps for HubState.
var (
	HubState_name = map[int32]string{
		0: "HUB_STATE_UNSPECIFIED",
		1: "HUB_STATE_READY",
		2: "HUB_STATE_DISABLED",
	}
	HubState_value = map[string]int32{
		"HUB_STATE_UNSPECIFIED": 0,
		"HUB_STATE_READY":       1,
		"HUB_STATE_DISABLED":    2,
	}
)

func (x HubState) Enum() *HubState {
	p := new(HubState)
	*p = x
	return p
}

func (x HubState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (HubState) Descriptor() protoreflect.EnumDescriptor {
	return file_admin_v1_hub_type_proto_enumTypes[0].Descriptor()
}

func (HubState) Type() protoreflect.EnumType {
	return &file_admin_v1_hub_type_proto_enumTypes[0]
}

func (x HubState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use HubState.Descriptor instead.
func (HubState) EnumDescriptor() ([]byte, []int) {
	return file_admin_v1_hub_type_proto_rawDescGZIP(), []int{0}
}

// Types of conditions used to describe a hub.
type HubConditionType int32

const (
	// Unspecified indicates that the condition unknown.
	//
	// This will never be appear in the `spec.conditions` field of a management cluster.
	HubConditionType_HUB_CONDITION_TYPE_UNSPECIFIED HubConditionType = 0
	// Accepted indicates that the management cluster is ready.
	HubConditionType_HUB_CONDITION_TYPE_READY HubConditionType = 1
	// Rejected indicates that the management cluster has been disabled by the administrator.
	HubConditionType_HUB_CONDITION_TYPE_DISABLED HubConditionType = 2
)

// Enum value maps for HubConditionType.
var (
	HubConditionType_name = map[int32]string{
		0: "HUB_CONDITION_TYPE_UNSPECIFIED",
		1: "HUB_CONDITION_TYPE_READY",
		2: "HUB_CONDITION_TYPE_DISABLED",
	}
	HubConditionType_value = map[string]int32{
		"HUB_CONDITION_TYPE_UNSPECIFIED": 0,
		"HUB_CONDITION_TYPE_READY":       1,
		"HUB_CONDITION_TYPE_DISABLED":    2,
	}
)

func (x HubConditionType) Enum() *HubConditionType {
	p := new(HubConditionType)
	*p = x
	return p
}

func (x HubConditionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (HubConditionType) Descriptor() protoreflect.EnumDescriptor {
	return file_admin_v1_hub_type_proto_enumTypes[1].Descriptor()
}

func (HubConditionType) Type() protoreflect.EnumType {
	return &file_admin_v1_hub_type_proto_enumTypes[1]
}

func (x HubConditionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use HubConditionType.Descriptor instead.
func (HubConditionType) EnumDescriptor() ([]byte, []int) {
	return file_admin_v1_hub_type_proto_rawDescGZIP(), []int{1}
}

// Contains the details of a hub.
type Hub struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique identifier of the hub.
	//
	// This will be automatically generated by the server when the hub is created.
	Id            string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Spec          *HubSpec   `protobuf:"bytes,2,opt,name=spec,proto3" json:"spec,omitempty"`
	Status        *HubStatus `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Hub) Reset() {
	*x = Hub{}
	mi := &file_admin_v1_hub_type_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Hub) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Hub) ProtoMessage() {}

func (x *Hub) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_hub_type_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Hub.ProtoReflect.Descriptor instead.
func (*Hub) Descriptor() ([]byte, []int) {
	return file_admin_v1_hub_type_proto_rawDescGZIP(), []int{0}
}

func (x *Hub) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Hub) GetSpec() *HubSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

func (x *Hub) GetStatus() *HubStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

// Contains the details that the user provides to create a hub.
type HubSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The Kubeconfig containing the address and credentials that the fulfillment service will use to connect to the hub.
	Kubeconfig []byte `protobuf:"bytes,1,opt,name=kubeconfig,proto3" json:"kubeconfig,omitempty"`
	// Namespace where the cluster orders will be created.
	Namespace     string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HubSpec) Reset() {
	*x = HubSpec{}
	mi := &file_admin_v1_hub_type_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HubSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HubSpec) ProtoMessage() {}

func (x *HubSpec) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_hub_type_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HubSpec.ProtoReflect.Descriptor instead.
func (*HubSpec) Descriptor() ([]byte, []int) {
	return file_admin_v1_hub_type_proto_rawDescGZIP(), []int{1}
}

func (x *HubSpec) GetKubeconfig() []byte {
	if x != nil {
		return x.Kubeconfig
	}
	return nil
}

func (x *HubSpec) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

// Contains the current status of a hub.
type HubStatus struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Indicates the overall state of a hub.
	//
	// For more details check the conditions.
	State HubState `protobuf:"varint,1,opt,name=state,proto3,enum=admin.v1.HubState" json:"state,omitempty"`
	// Contains a list of conditions that describe in detail the hub.
	//
	// For example, hub that has been disabled by the administrator would be represented like this (when converted to
	// JSON):
	//
	//	{
	//	  "id": "123",
	//	  "spec": {
	//	    "kubeconfig": "..."
	//	  },
	//	  "state": "HUB_STATE_DISABLED",
	//	  "status": {
	//	    "conditions": [
	//	      {
	//	        "type: "HUB_CONDITION_TYPE_READY",
	//	        "status": "CONDITION_STATUS_FALSE",
	//	        "last_transition_time": "2025-03-12 20:15:59+00:00"
	//	      },
	//	      {
	//	        "type": "HUB_CONDITION_TYPE_DISABLED",
	//	        "status": "CONDITION_STATUS_FALSE",
	//	        "last_transition_time": "2025-03-12 20:17:16+00:00",
	//	        "reason": "Disabled",
	//	        "message": "The management cluster has been disabled by the administrator"
	//	      },
	//	    ]
	//	  }
	//	}
	//
	// In this example the `DISABLED` condition is true. That tells us that the hub has disabled by the administrator.
	Conditions    []*HubCondition `protobuf:"bytes,2,rep,name=conditions,proto3" json:"conditions,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HubStatus) Reset() {
	*x = HubStatus{}
	mi := &file_admin_v1_hub_type_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HubStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HubStatus) ProtoMessage() {}

func (x *HubStatus) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_hub_type_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HubStatus.ProtoReflect.Descriptor instead.
func (*HubStatus) Descriptor() ([]byte, []int) {
	return file_admin_v1_hub_type_proto_rawDescGZIP(), []int{2}
}

func (x *HubStatus) GetState() HubState {
	if x != nil {
		return x.State
	}
	return HubState_HUB_STATE_UNSPECIFIED
}

func (x *HubStatus) GetConditions() []*HubCondition {
	if x != nil {
		return x.Conditions
	}
	return nil
}

// Contains the details of a condition that describes the status of a hub.
type HubCondition struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Indicates the type of condition.
	Type HubConditionType `protobuf:"varint,1,opt,name=type,proto3,enum=admin.v1.HubConditionType" json:"type,omitempty"`
	// Indicates status of the condition.
	Status v1.ConditionStatus `protobuf:"varint,2,opt,name=status,proto3,enum=shared.v1.ConditionStatus" json:"status,omitempty"`
	// This time is the last time that the condition was updated.
	LastTransitionTime *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=last_transition_time,json=lastTransitionTime,proto3" json:"last_transition_time,omitempty"`
	// Contains a the reason of the condition in a format suitable for use by programs.
	//
	// The possible are documented in the `HubConditionType` object.
	Reason *string `protobuf:"bytes,4,opt,name=reason,proto3,oneof" json:"reason,omitempty"`
	// Contains a text giving more details of the condition. This will usually be progress reports, or error messages, and
	// are intended for use by humans, to debug problems.
	Message       *string `protobuf:"bytes,5,opt,name=message,proto3,oneof" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HubCondition) Reset() {
	*x = HubCondition{}
	mi := &file_admin_v1_hub_type_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HubCondition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HubCondition) ProtoMessage() {}

func (x *HubCondition) ProtoReflect() protoreflect.Message {
	mi := &file_admin_v1_hub_type_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HubCondition.ProtoReflect.Descriptor instead.
func (*HubCondition) Descriptor() ([]byte, []int) {
	return file_admin_v1_hub_type_proto_rawDescGZIP(), []int{3}
}

func (x *HubCondition) GetType() HubConditionType {
	if x != nil {
		return x.Type
	}
	return HubConditionType_HUB_CONDITION_TYPE_UNSPECIFIED
}

func (x *HubCondition) GetStatus() v1.ConditionStatus {
	if x != nil {
		return x.Status
	}
	return v1.ConditionStatus(0)
}

func (x *HubCondition) GetLastTransitionTime() *timestamppb.Timestamp {
	if x != nil {
		return x.LastTransitionTime
	}
	return nil
}

func (x *HubCondition) GetReason() string {
	if x != nil && x.Reason != nil {
		return *x.Reason
	}
	return ""
}

func (x *HubCondition) GetMessage() string {
	if x != nil && x.Message != nil {
		return *x.Message
	}
	return ""
}

var File_admin_v1_hub_type_proto protoreflect.FileDescriptor

var file_admin_v1_hub_type_proto_rawDesc = string([]byte{
	0x0a, 0x17, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x75, 0x62, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x25, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x76, 0x31, 0x2f,
	0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x69, 0x0a, 0x03, 0x48,
	0x75, 0x62, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x25, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x11, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x53,
	0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x12, 0x2b, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x47, 0x0a, 0x07, 0x48, 0x75, 0x62, 0x53, 0x70, 0x65,
	0x63, 0x12, 0x1e, 0x0a, 0x0a, 0x6b, 0x75, 0x62, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x6b, 0x75, 0x62, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x22,
	0x6d, 0x0a, 0x09, 0x48, 0x75, 0x62, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x28, 0x0a, 0x05,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x36, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x93,
	0x02, 0x0a, 0x0c, 0x48, 0x75, 0x62, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x2e, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x43, 0x6f, 0x6e, 0x64,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12,
	0x32, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1a, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x64,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x4c, 0x0a, 0x14, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x12, 0x6c,
	0x61, 0x73, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x1b, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x48, 0x00, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x1d,
	0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x01, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x88, 0x01, 0x01, 0x42, 0x09, 0x0a,
	0x07, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2a, 0x52, 0x0a, 0x08, 0x48, 0x75, 0x62, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x19, 0x0a, 0x15, 0x48, 0x55, 0x42, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x55, 0x4e,
	0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x48,
	0x55, 0x42, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x52, 0x45, 0x41, 0x44, 0x59, 0x10, 0x01,
	0x12, 0x16, 0x0a, 0x12, 0x48, 0x55, 0x42, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x44, 0x49,
	0x53, 0x41, 0x42, 0x4c, 0x45, 0x44, 0x10, 0x02, 0x2a, 0x75, 0x0a, 0x10, 0x48, 0x75, 0x62, 0x43,
	0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x22, 0x0a, 0x1e,
	0x48, 0x55, 0x42, 0x5f, 0x43, 0x4f, 0x4e, 0x44, 0x49, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x1c, 0x0a, 0x18, 0x48, 0x55, 0x42, 0x5f, 0x43, 0x4f, 0x4e, 0x44, 0x49, 0x54, 0x49, 0x4f,
	0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52, 0x45, 0x41, 0x44, 0x59, 0x10, 0x01, 0x12, 0x1f,
	0x0a, 0x1b, 0x48, 0x55, 0x42, 0x5f, 0x43, 0x4f, 0x4e, 0x44, 0x49, 0x54, 0x49, 0x4f, 0x4e, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x49, 0x53, 0x41, 0x42, 0x4c, 0x45, 0x44, 0x10, 0x02, 0x42,
	0x9f, 0x01, 0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x76, 0x31,
	0x42, 0x0c, 0x48, 0x75, 0x62, 0x54, 0x79, 0x70, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x6e,
	0x61, 0x62, 0x6f, 0x78, 0x2f, 0x66, 0x75, 0x6c, 0x66, 0x69, 0x6c, 0x6c, 0x6d, 0x65, 0x6e, 0x74,
	0x2d, 0x63, 0x6c, 0x69, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x58, 0x58, 0xaa, 0x02, 0x08, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x2e, 0x56, 0x31, 0xca, 0x02, 0x08, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x5c, 0x56, 0x31, 0xe2, 0x02,
	0x14, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x09, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x3a, 0x3a, 0x56,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_admin_v1_hub_type_proto_rawDescOnce sync.Once
	file_admin_v1_hub_type_proto_rawDescData []byte
)

func file_admin_v1_hub_type_proto_rawDescGZIP() []byte {
	file_admin_v1_hub_type_proto_rawDescOnce.Do(func() {
		file_admin_v1_hub_type_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_admin_v1_hub_type_proto_rawDesc), len(file_admin_v1_hub_type_proto_rawDesc)))
	})
	return file_admin_v1_hub_type_proto_rawDescData
}

var file_admin_v1_hub_type_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_admin_v1_hub_type_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_admin_v1_hub_type_proto_goTypes = []any{
	(HubState)(0),                 // 0: admin.v1.HubState
	(HubConditionType)(0),         // 1: admin.v1.HubConditionType
	(*Hub)(nil),                   // 2: admin.v1.Hub
	(*HubSpec)(nil),               // 3: admin.v1.HubSpec
	(*HubStatus)(nil),             // 4: admin.v1.HubStatus
	(*HubCondition)(nil),          // 5: admin.v1.HubCondition
	(v1.ConditionStatus)(0),       // 6: shared.v1.ConditionStatus
	(*timestamppb.Timestamp)(nil), // 7: google.protobuf.Timestamp
}
var file_admin_v1_hub_type_proto_depIdxs = []int32{
	3, // 0: admin.v1.Hub.spec:type_name -> admin.v1.HubSpec
	4, // 1: admin.v1.Hub.status:type_name -> admin.v1.HubStatus
	0, // 2: admin.v1.HubStatus.state:type_name -> admin.v1.HubState
	5, // 3: admin.v1.HubStatus.conditions:type_name -> admin.v1.HubCondition
	1, // 4: admin.v1.HubCondition.type:type_name -> admin.v1.HubConditionType
	6, // 5: admin.v1.HubCondition.status:type_name -> shared.v1.ConditionStatus
	7, // 6: admin.v1.HubCondition.last_transition_time:type_name -> google.protobuf.Timestamp
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_admin_v1_hub_type_proto_init() }
func file_admin_v1_hub_type_proto_init() {
	if File_admin_v1_hub_type_proto != nil {
		return
	}
	file_admin_v1_hub_type_proto_msgTypes[3].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_admin_v1_hub_type_proto_rawDesc), len(file_admin_v1_hub_type_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_admin_v1_hub_type_proto_goTypes,
		DependencyIndexes: file_admin_v1_hub_type_proto_depIdxs,
		EnumInfos:         file_admin_v1_hub_type_proto_enumTypes,
		MessageInfos:      file_admin_v1_hub_type_proto_msgTypes,
	}.Build()
	File_admin_v1_hub_type_proto = out.File
	file_admin_v1_hub_type_proto_goTypes = nil
	file_admin_v1_hub_type_proto_depIdxs = nil
}
