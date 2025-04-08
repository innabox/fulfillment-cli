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
// source: events/v1/events_service.proto

//go:build !protoopaque

package eventsv1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type EventsWatchRequest struct {
	state protoimpl.MessageState `protogen:"hybrid.v1"`
	// Filter criteria.
	//
	// The value of this parameter is a [CEL](https://cel.dev) boolean expression. The `event` variable will contain the
	// fields of the event. If the result of the expression is `true` then the event will be sent by the server. For
	// example, to receive only the events that indicate that a cluster order has been modified and is now in the
	// fulfilled state:
	//
	// ```
	// event.type == EVENT_TYPE_OBJECT_CREATED && event.cluster_order.status.state == CLUSTER_ORDER_STATE_FULFILLED
	// ```
	//
	// If this isn't provided, or if the value is empty, then all the events that the user has permission to see will be
	// sent by the server.
	Filter        *string `protobuf:"bytes,1,opt,name=filter,proto3,oneof" json:"filter,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EventsWatchRequest) Reset() {
	*x = EventsWatchRequest{}
	mi := &file_events_v1_events_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventsWatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventsWatchRequest) ProtoMessage() {}

func (x *EventsWatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_events_v1_events_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *EventsWatchRequest) GetFilter() string {
	if x != nil && x.Filter != nil {
		return *x.Filter
	}
	return ""
}

func (x *EventsWatchRequest) SetFilter(v string) {
	x.Filter = &v
}

func (x *EventsWatchRequest) HasFilter() bool {
	if x == nil {
		return false
	}
	return x.Filter != nil
}

func (x *EventsWatchRequest) ClearFilter() {
	x.Filter = nil
}

type EventsWatchRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	// Filter criteria.
	//
	// The value of this parameter is a [CEL](https://cel.dev) boolean expression. The `event` variable will contain the
	// fields of the event. If the result of the expression is `true` then the event will be sent by the server. For
	// example, to receive only the events that indicate that a cluster order has been modified and is now in the
	// fulfilled state:
	//
	// ```
	// event.type == EVENT_TYPE_OBJECT_CREATED && event.cluster_order.status.state == CLUSTER_ORDER_STATE_FULFILLED
	// ```
	//
	// If this isn't provided, or if the value is empty, then all the events that the user has permission to see will be
	// sent by the server.
	Filter *string
}

func (b0 EventsWatchRequest_builder) Build() *EventsWatchRequest {
	m0 := &EventsWatchRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.Filter = b.Filter
	return m0
}

type EventsWatchResponse struct {
	state         protoimpl.MessageState `protogen:"hybrid.v1"`
	Event         *Event                 `protobuf:"bytes,1,opt,name=event,proto3" json:"event,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EventsWatchResponse) Reset() {
	*x = EventsWatchResponse{}
	mi := &file_events_v1_events_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventsWatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventsWatchResponse) ProtoMessage() {}

func (x *EventsWatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_events_v1_events_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *EventsWatchResponse) GetEvent() *Event {
	if x != nil {
		return x.Event
	}
	return nil
}

func (x *EventsWatchResponse) SetEvent(v *Event) {
	x.Event = v
}

func (x *EventsWatchResponse) HasEvent() bool {
	if x == nil {
		return false
	}
	return x.Event != nil
}

func (x *EventsWatchResponse) ClearEvent() {
	x.Event = nil
}

type EventsWatchResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Event *Event
}

func (b0 EventsWatchResponse_builder) Build() *EventsWatchResponse {
	m0 := &EventsWatchResponse{}
	b, x := &b0, m0
	_, _ = b, x
	x.Event = b.Event
	return m0
}

var File_events_v1_events_service_proto protoreflect.FileDescriptor

var file_events_v1_events_service_proto_rawDesc = string([]byte{
	0x0a, 0x1e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x1a, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3c, 0x0a, 0x12, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x57,
	0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x06, 0x66,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x06, 0x66,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x88, 0x01, 0x01, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x66, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x22, 0x3d, 0x0a, 0x13, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x57, 0x61, 0x74,
	0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x05, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x32, 0x71, 0x0a, 0x06, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x67, 0x0a, 0x05,
	0x57, 0x61, 0x74, 0x63, 0x68, 0x12, 0x1d, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x57, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x57, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x12, 0x15, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x30, 0x01, 0x42, 0xac, 0x01, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x12, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x42, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x6e, 0x61, 0x62, 0x6f,
	0x78, 0x2f, 0x66, 0x75, 0x6c, 0x66, 0x69, 0x6c, 0x6c, 0x6d, 0x65, 0x6e, 0x74, 0x2d, 0x63, 0x6c,
	0x69, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x76,
	0x31, 0xa2, 0x02, 0x03, 0x45, 0x58, 0x58, 0xaa, 0x02, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x2e, 0x56, 0x31, 0xca, 0x02, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x5c, 0x56, 0x31, 0xe2,
	0x02, 0x15, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0a, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var file_events_v1_events_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_events_v1_events_service_proto_goTypes = []any{
	(*EventsWatchRequest)(nil),  // 0: events.v1.EventsWatchRequest
	(*EventsWatchResponse)(nil), // 1: events.v1.EventsWatchResponse
	(*Event)(nil),               // 2: events.v1.Event
}
var file_events_v1_events_service_proto_depIdxs = []int32{
	2, // 0: events.v1.EventsWatchResponse.event:type_name -> events.v1.Event
	0, // 1: events.v1.Events.Watch:input_type -> events.v1.EventsWatchRequest
	1, // 2: events.v1.Events.Watch:output_type -> events.v1.EventsWatchResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_events_v1_events_service_proto_init() }
func file_events_v1_events_service_proto_init() {
	if File_events_v1_events_service_proto != nil {
		return
	}
	file_events_v1_event_type_proto_init()
	file_events_v1_events_service_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_events_v1_events_service_proto_rawDesc), len(file_events_v1_events_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_events_v1_events_service_proto_goTypes,
		DependencyIndexes: file_events_v1_events_service_proto_depIdxs,
		MessageInfos:      file_events_v1_events_service_proto_msgTypes,
	}.Build()
	File_events_v1_events_service_proto = out.File
	file_events_v1_events_service_proto_goTypes = nil
	file_events_v1_events_service_proto_depIdxs = nil
}
