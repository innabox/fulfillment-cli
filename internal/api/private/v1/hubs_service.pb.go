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
// source: private/v1/hubs_service.proto

//go:build !protoopaque

package privatev1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type HubsListRequest struct {
	state         protoimpl.MessageState `protogen:"hybrid.v1"`
	Offset        *int32                 `protobuf:"varint,1,opt,name=offset,proto3,oneof" json:"offset,omitempty"`
	Limit         *int32                 `protobuf:"varint,2,opt,name=limit,proto3,oneof" json:"limit,omitempty"`
	Filter        *string                `protobuf:"bytes,3,opt,name=filter,proto3,oneof" json:"filter,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HubsListRequest) Reset() {
	*x = HubsListRequest{}
	mi := &file_private_v1_hubs_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HubsListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HubsListRequest) ProtoMessage() {}

func (x *HubsListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_private_v1_hubs_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *HubsListRequest) GetOffset() int32 {
	if x != nil && x.Offset != nil {
		return *x.Offset
	}
	return 0
}

func (x *HubsListRequest) GetLimit() int32 {
	if x != nil && x.Limit != nil {
		return *x.Limit
	}
	return 0
}

func (x *HubsListRequest) GetFilter() string {
	if x != nil && x.Filter != nil {
		return *x.Filter
	}
	return ""
}

func (x *HubsListRequest) SetOffset(v int32) {
	x.Offset = &v
}

func (x *HubsListRequest) SetLimit(v int32) {
	x.Limit = &v
}

func (x *HubsListRequest) SetFilter(v string) {
	x.Filter = &v
}

func (x *HubsListRequest) HasOffset() bool {
	if x == nil {
		return false
	}
	return x.Offset != nil
}

func (x *HubsListRequest) HasLimit() bool {
	if x == nil {
		return false
	}
	return x.Limit != nil
}

func (x *HubsListRequest) HasFilter() bool {
	if x == nil {
		return false
	}
	return x.Filter != nil
}

func (x *HubsListRequest) ClearOffset() {
	x.Offset = nil
}

func (x *HubsListRequest) ClearLimit() {
	x.Limit = nil
}

func (x *HubsListRequest) ClearFilter() {
	x.Filter = nil
}

type HubsListRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Offset *int32
	Limit  *int32
	Filter *string
}

func (b0 HubsListRequest_builder) Build() *HubsListRequest {
	m0 := &HubsListRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.Offset = b.Offset
	x.Limit = b.Limit
	x.Filter = b.Filter
	return m0
}

type HubsListResponse struct {
	state         protoimpl.MessageState `protogen:"hybrid.v1"`
	Size          *int32                 `protobuf:"varint,1,opt,name=size,proto3,oneof" json:"size,omitempty"`
	Total         *int32                 `protobuf:"varint,2,opt,name=total,proto3,oneof" json:"total,omitempty"`
	Items         []*Hub                 `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HubsListResponse) Reset() {
	*x = HubsListResponse{}
	mi := &file_private_v1_hubs_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HubsListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HubsListResponse) ProtoMessage() {}

func (x *HubsListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_private_v1_hubs_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *HubsListResponse) GetSize() int32 {
	if x != nil && x.Size != nil {
		return *x.Size
	}
	return 0
}

func (x *HubsListResponse) GetTotal() int32 {
	if x != nil && x.Total != nil {
		return *x.Total
	}
	return 0
}

func (x *HubsListResponse) GetItems() []*Hub {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *HubsListResponse) SetSize(v int32) {
	x.Size = &v
}

func (x *HubsListResponse) SetTotal(v int32) {
	x.Total = &v
}

func (x *HubsListResponse) SetItems(v []*Hub) {
	x.Items = v
}

func (x *HubsListResponse) HasSize() bool {
	if x == nil {
		return false
	}
	return x.Size != nil
}

func (x *HubsListResponse) HasTotal() bool {
	if x == nil {
		return false
	}
	return x.Total != nil
}

func (x *HubsListResponse) ClearSize() {
	x.Size = nil
}

func (x *HubsListResponse) ClearTotal() {
	x.Total = nil
}

type HubsListResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Size  *int32
	Total *int32
	Items []*Hub
}

func (b0 HubsListResponse_builder) Build() *HubsListResponse {
	m0 := &HubsListResponse{}
	b, x := &b0, m0
	_, _ = b, x
	x.Size = b.Size
	x.Total = b.Total
	x.Items = b.Items
	return m0
}

type HubsGetRequest struct {
	state         protoimpl.MessageState `protogen:"hybrid.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HubsGetRequest) Reset() {
	*x = HubsGetRequest{}
	mi := &file_private_v1_hubs_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HubsGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HubsGetRequest) ProtoMessage() {}

func (x *HubsGetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_private_v1_hubs_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *HubsGetRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *HubsGetRequest) SetId(v string) {
	x.Id = v
}

type HubsGetRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Id string
}

func (b0 HubsGetRequest_builder) Build() *HubsGetRequest {
	m0 := &HubsGetRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.Id = b.Id
	return m0
}

type HubsGetResponse struct {
	state         protoimpl.MessageState `protogen:"hybrid.v1"`
	Object        *Hub                   `protobuf:"bytes,1,opt,name=object,proto3" json:"object,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HubsGetResponse) Reset() {
	*x = HubsGetResponse{}
	mi := &file_private_v1_hubs_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HubsGetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HubsGetResponse) ProtoMessage() {}

func (x *HubsGetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_private_v1_hubs_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *HubsGetResponse) GetObject() *Hub {
	if x != nil {
		return x.Object
	}
	return nil
}

func (x *HubsGetResponse) SetObject(v *Hub) {
	x.Object = v
}

func (x *HubsGetResponse) HasObject() bool {
	if x == nil {
		return false
	}
	return x.Object != nil
}

func (x *HubsGetResponse) ClearObject() {
	x.Object = nil
}

type HubsGetResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Object *Hub
}

func (b0 HubsGetResponse_builder) Build() *HubsGetResponse {
	m0 := &HubsGetResponse{}
	b, x := &b0, m0
	_, _ = b, x
	x.Object = b.Object
	return m0
}

type HubsCreateRequest struct {
	state         protoimpl.MessageState `protogen:"hybrid.v1"`
	Object        *Hub                   `protobuf:"bytes,1,opt,name=object,proto3" json:"object,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HubsCreateRequest) Reset() {
	*x = HubsCreateRequest{}
	mi := &file_private_v1_hubs_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HubsCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HubsCreateRequest) ProtoMessage() {}

func (x *HubsCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_private_v1_hubs_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *HubsCreateRequest) GetObject() *Hub {
	if x != nil {
		return x.Object
	}
	return nil
}

func (x *HubsCreateRequest) SetObject(v *Hub) {
	x.Object = v
}

func (x *HubsCreateRequest) HasObject() bool {
	if x == nil {
		return false
	}
	return x.Object != nil
}

func (x *HubsCreateRequest) ClearObject() {
	x.Object = nil
}

type HubsCreateRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Object *Hub
}

func (b0 HubsCreateRequest_builder) Build() *HubsCreateRequest {
	m0 := &HubsCreateRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.Object = b.Object
	return m0
}

type HubsCreateResponse struct {
	state         protoimpl.MessageState `protogen:"hybrid.v1"`
	Object        *Hub                   `protobuf:"bytes,1,opt,name=object,proto3" json:"object,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HubsCreateResponse) Reset() {
	*x = HubsCreateResponse{}
	mi := &file_private_v1_hubs_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HubsCreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HubsCreateResponse) ProtoMessage() {}

func (x *HubsCreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_private_v1_hubs_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *HubsCreateResponse) GetObject() *Hub {
	if x != nil {
		return x.Object
	}
	return nil
}

func (x *HubsCreateResponse) SetObject(v *Hub) {
	x.Object = v
}

func (x *HubsCreateResponse) HasObject() bool {
	if x == nil {
		return false
	}
	return x.Object != nil
}

func (x *HubsCreateResponse) ClearObject() {
	x.Object = nil
}

type HubsCreateResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Object *Hub
}

func (b0 HubsCreateResponse_builder) Build() *HubsCreateResponse {
	m0 := &HubsCreateResponse{}
	b, x := &b0, m0
	_, _ = b, x
	x.Object = b.Object
	return m0
}

type HubsDeleteRequest struct {
	state         protoimpl.MessageState `protogen:"hybrid.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HubsDeleteRequest) Reset() {
	*x = HubsDeleteRequest{}
	mi := &file_private_v1_hubs_service_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HubsDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HubsDeleteRequest) ProtoMessage() {}

func (x *HubsDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_private_v1_hubs_service_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *HubsDeleteRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *HubsDeleteRequest) SetId(v string) {
	x.Id = v
}

type HubsDeleteRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Id string
}

func (b0 HubsDeleteRequest_builder) Build() *HubsDeleteRequest {
	m0 := &HubsDeleteRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.Id = b.Id
	return m0
}

type HubsDeleteResponse struct {
	state         protoimpl.MessageState `protogen:"hybrid.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HubsDeleteResponse) Reset() {
	*x = HubsDeleteResponse{}
	mi := &file_private_v1_hubs_service_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HubsDeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HubsDeleteResponse) ProtoMessage() {}

func (x *HubsDeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_private_v1_hubs_service_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

type HubsDeleteResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

}

func (b0 HubsDeleteResponse_builder) Build() *HubsDeleteResponse {
	m0 := &HubsDeleteResponse{}
	b, x := &b0, m0
	_, _ = b, x
	return m0
}

type HubsUpdateRequest struct {
	state         protoimpl.MessageState `protogen:"hybrid.v1"`
	Object        *Hub                   `protobuf:"bytes,1,opt,name=object,proto3" json:"object,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HubsUpdateRequest) Reset() {
	*x = HubsUpdateRequest{}
	mi := &file_private_v1_hubs_service_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HubsUpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HubsUpdateRequest) ProtoMessage() {}

func (x *HubsUpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_private_v1_hubs_service_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *HubsUpdateRequest) GetObject() *Hub {
	if x != nil {
		return x.Object
	}
	return nil
}

func (x *HubsUpdateRequest) SetObject(v *Hub) {
	x.Object = v
}

func (x *HubsUpdateRequest) HasObject() bool {
	if x == nil {
		return false
	}
	return x.Object != nil
}

func (x *HubsUpdateRequest) ClearObject() {
	x.Object = nil
}

type HubsUpdateRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Object *Hub
}

func (b0 HubsUpdateRequest_builder) Build() *HubsUpdateRequest {
	m0 := &HubsUpdateRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.Object = b.Object
	return m0
}

type HubsUpdateResponse struct {
	state         protoimpl.MessageState `protogen:"hybrid.v1"`
	Object        *Hub                   `protobuf:"bytes,1,opt,name=object,proto3" json:"object,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HubsUpdateResponse) Reset() {
	*x = HubsUpdateResponse{}
	mi := &file_private_v1_hubs_service_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HubsUpdateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HubsUpdateResponse) ProtoMessage() {}

func (x *HubsUpdateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_private_v1_hubs_service_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *HubsUpdateResponse) GetObject() *Hub {
	if x != nil {
		return x.Object
	}
	return nil
}

func (x *HubsUpdateResponse) SetObject(v *Hub) {
	x.Object = v
}

func (x *HubsUpdateResponse) HasObject() bool {
	if x == nil {
		return false
	}
	return x.Object != nil
}

func (x *HubsUpdateResponse) ClearObject() {
	x.Object = nil
}

type HubsUpdateResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Object *Hub
}

func (b0 HubsUpdateResponse_builder) Build() *HubsUpdateResponse {
	m0 := &HubsUpdateResponse{}
	b, x := &b0, m0
	_, _ = b, x
	x.Object = b.Object
	return m0
}

var File_private_v1_hubs_service_proto protoreflect.FileDescriptor

var file_private_v1_hubs_service_proto_rawDesc = string([]byte{
	0x0a, 0x1d, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x75, 0x62,
	0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0a, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x19, 0x70, 0x72, 0x69,
	0x76, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x68, 0x75, 0x62, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x86, 0x01, 0x0a, 0x0f, 0x48, 0x75, 0x62, 0x73, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x06, 0x6f, 0x66,
	0x66, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x06, 0x6f, 0x66,
	0x66, 0x73, 0x65, 0x74, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x48, 0x01, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x88,
	0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x02, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x88, 0x01, 0x01, 0x42,
	0x09, 0x0a, 0x07, 0x5f, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x22,
	0x80, 0x01, 0x0a, 0x10, 0x48, 0x75, 0x62, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x48, 0x00, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a,
	0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x48, 0x01, 0x52, 0x05,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x88, 0x01, 0x01, 0x12, 0x25, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x42,
	0x07, 0x0a, 0x05, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x22, 0x20, 0x0a, 0x0e, 0x48, 0x75, 0x62, 0x73, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x22, 0x3a, 0x0a, 0x0f, 0x48, 0x75, 0x62, 0x73, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x52, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x22, 0x3c, 0x0a, 0x11, 0x48, 0x75, 0x62, 0x73, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x52, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x22, 0x3d,
	0x0a, 0x12, 0x48, 0x75, 0x62, 0x73, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x48, 0x75, 0x62, 0x52, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x22, 0x23, 0x0a,
	0x11, 0x48, 0x75, 0x62, 0x73, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x22, 0x14, 0x0a, 0x12, 0x48, 0x75, 0x62, 0x73, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3c, 0x0a, 0x11, 0x48, 0x75, 0x62, 0x73,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a,
	0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x52, 0x06,
	0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x22, 0x3d, 0x0a, 0x12, 0x48, 0x75, 0x62, 0x73, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x06,
	0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70,
	0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x52, 0x06, 0x6f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x32, 0xee, 0x02, 0x0a, 0x04, 0x48, 0x75, 0x62, 0x73, 0x12, 0x43,
	0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x48, 0x75, 0x62, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x40, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x69,
	0x76, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x73, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x73, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x49, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12,
	0x1d, 0x2e, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62,
	0x73, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e,
	0x2e, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x73,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x49, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x70, 0x72, 0x69,
	0x76, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x73, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x70, 0x72, 0x69, 0x76,
	0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x73, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x49, 0x0a, 0x06, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x48, 0x75, 0x62, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x48, 0x75, 0x62, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0xb3, 0x01, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x2e, 0x70,
	0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x10, 0x48, 0x75, 0x62, 0x73, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x44, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x6e, 0x61, 0x62, 0x6f,
	0x78, 0x2f, 0x66, 0x75, 0x6c, 0x66, 0x69, 0x6c, 0x6c, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x63, 0x6c,
	0x69, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70,
	0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74,
	0x65, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x0a, 0x50, 0x72, 0x69, 0x76,
	0x61, 0x74, 0x65, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0b, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65,
	0x5f, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x17, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x5c,
	0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x0b, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
})

var file_private_v1_hubs_service_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_private_v1_hubs_service_proto_goTypes = []any{
	(*HubsListRequest)(nil),    // 0: private.v1.HubsListRequest
	(*HubsListResponse)(nil),   // 1: private.v1.HubsListResponse
	(*HubsGetRequest)(nil),     // 2: private.v1.HubsGetRequest
	(*HubsGetResponse)(nil),    // 3: private.v1.HubsGetResponse
	(*HubsCreateRequest)(nil),  // 4: private.v1.HubsCreateRequest
	(*HubsCreateResponse)(nil), // 5: private.v1.HubsCreateResponse
	(*HubsDeleteRequest)(nil),  // 6: private.v1.HubsDeleteRequest
	(*HubsDeleteResponse)(nil), // 7: private.v1.HubsDeleteResponse
	(*HubsUpdateRequest)(nil),  // 8: private.v1.HubsUpdateRequest
	(*HubsUpdateResponse)(nil), // 9: private.v1.HubsUpdateResponse
	(*Hub)(nil),                // 10: private.v1.Hub
}
var file_private_v1_hubs_service_proto_depIdxs = []int32{
	10, // 0: private.v1.HubsListResponse.items:type_name -> private.v1.Hub
	10, // 1: private.v1.HubsGetResponse.object:type_name -> private.v1.Hub
	10, // 2: private.v1.HubsCreateRequest.object:type_name -> private.v1.Hub
	10, // 3: private.v1.HubsCreateResponse.object:type_name -> private.v1.Hub
	10, // 4: private.v1.HubsUpdateRequest.object:type_name -> private.v1.Hub
	10, // 5: private.v1.HubsUpdateResponse.object:type_name -> private.v1.Hub
	0,  // 6: private.v1.Hubs.List:input_type -> private.v1.HubsListRequest
	2,  // 7: private.v1.Hubs.Get:input_type -> private.v1.HubsGetRequest
	4,  // 8: private.v1.Hubs.Create:input_type -> private.v1.HubsCreateRequest
	6,  // 9: private.v1.Hubs.Delete:input_type -> private.v1.HubsDeleteRequest
	8,  // 10: private.v1.Hubs.Update:input_type -> private.v1.HubsUpdateRequest
	1,  // 11: private.v1.Hubs.List:output_type -> private.v1.HubsListResponse
	3,  // 12: private.v1.Hubs.Get:output_type -> private.v1.HubsGetResponse
	5,  // 13: private.v1.Hubs.Create:output_type -> private.v1.HubsCreateResponse
	7,  // 14: private.v1.Hubs.Delete:output_type -> private.v1.HubsDeleteResponse
	9,  // 15: private.v1.Hubs.Update:output_type -> private.v1.HubsUpdateResponse
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_private_v1_hubs_service_proto_init() }
func file_private_v1_hubs_service_proto_init() {
	if File_private_v1_hubs_service_proto != nil {
		return
	}
	file_private_v1_hub_type_proto_init()
	file_private_v1_hubs_service_proto_msgTypes[0].OneofWrappers = []any{}
	file_private_v1_hubs_service_proto_msgTypes[1].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_private_v1_hubs_service_proto_rawDesc), len(file_private_v1_hubs_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_private_v1_hubs_service_proto_goTypes,
		DependencyIndexes: file_private_v1_hubs_service_proto_depIdxs,
		MessageInfos:      file_private_v1_hubs_service_proto_msgTypes,
	}.Build()
	File_private_v1_hubs_service_proto = out.File
	file_private_v1_hubs_service_proto_goTypes = nil
	file_private_v1_hubs_service_proto_depIdxs = nil
}
