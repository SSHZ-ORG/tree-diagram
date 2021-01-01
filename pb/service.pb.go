// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: service.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RenderEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id *string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (x *RenderEventRequest) Reset() {
	*x = RenderEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RenderEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RenderEventRequest) ProtoMessage() {}

func (x *RenderEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RenderEventRequest.ProtoReflect.Descriptor instead.
func (*RenderEventRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{0}
}

func (x *RenderEventRequest) GetId() string {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return ""
}

type RenderEventResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Date               *string                                  `protobuf:"bytes,1,opt,name=date" json:"date,omitempty"`
	Snapshots          []*RenderEventResponse_Snapshot          `protobuf:"bytes,2,rep,name=snapshots" json:"snapshots,omitempty"`
	PlaceStatsTotal    *RenderEventResponse_PlaceNoteCountStats `protobuf:"bytes,3,opt,name=place_stats_total,json=placeStatsTotal" json:"place_stats_total,omitempty"`
	PlaceStatsFinished *RenderEventResponse_PlaceNoteCountStats `protobuf:"bytes,4,opt,name=place_stats_finished,json=placeStatsFinished" json:"place_stats_finished,omitempty"`
}

func (x *RenderEventResponse) Reset() {
	*x = RenderEventResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RenderEventResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RenderEventResponse) ProtoMessage() {}

func (x *RenderEventResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RenderEventResponse.ProtoReflect.Descriptor instead.
func (*RenderEventResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{1}
}

func (x *RenderEventResponse) GetDate() string {
	if x != nil && x.Date != nil {
		return *x.Date
	}
	return ""
}

func (x *RenderEventResponse) GetSnapshots() []*RenderEventResponse_Snapshot {
	if x != nil {
		return x.Snapshots
	}
	return nil
}

func (x *RenderEventResponse) GetPlaceStatsTotal() *RenderEventResponse_PlaceNoteCountStats {
	if x != nil {
		return x.PlaceStatsTotal
	}
	return nil
}

func (x *RenderEventResponse) GetPlaceStatsFinished() *RenderEventResponse_PlaceNoteCountStats {
	if x != nil {
		return x.PlaceStatsFinished
	}
	return nil
}

type RenderPlaceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id *string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (x *RenderPlaceRequest) Reset() {
	*x = RenderPlaceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RenderPlaceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RenderPlaceRequest) ProtoMessage() {}

func (x *RenderPlaceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RenderPlaceRequest.ProtoReflect.Descriptor instead.
func (*RenderPlaceRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{2}
}

func (x *RenderPlaceRequest) GetId() string {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return ""
}

type RenderPlaceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	KnownEventCount *int32 `protobuf:"varint,1,opt,name=known_event_count,json=knownEventCount" json:"known_event_count,omitempty"`
}

func (x *RenderPlaceResponse) Reset() {
	*x = RenderPlaceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RenderPlaceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RenderPlaceResponse) ProtoMessage() {}

func (x *RenderPlaceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RenderPlaceResponse.ProtoReflect.Descriptor instead.
func (*RenderPlaceResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{3}
}

func (x *RenderPlaceResponse) GetKnownEventCount() int32 {
	if x != nil && x.KnownEventCount != nil {
		return *x.KnownEventCount
	}
	return 0
}

type RenderEventResponse_Snapshot struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp     *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=timestamp" json:"timestamp,omitempty"`
	NoteCount     *int32                 `protobuf:"varint,2,opt,name=note_count,json=noteCount" json:"note_count,omitempty"`
	AddedActors   []string               `protobuf:"bytes,3,rep,name=added_actors,json=addedActors" json:"added_actors,omitempty"`
	RemovedActors []string               `protobuf:"bytes,4,rep,name=removed_actors,json=removedActors" json:"removed_actors,omitempty"`
}

func (x *RenderEventResponse_Snapshot) Reset() {
	*x = RenderEventResponse_Snapshot{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RenderEventResponse_Snapshot) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RenderEventResponse_Snapshot) ProtoMessage() {}

func (x *RenderEventResponse_Snapshot) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RenderEventResponse_Snapshot.ProtoReflect.Descriptor instead.
func (*RenderEventResponse_Snapshot) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{1, 0}
}

func (x *RenderEventResponse_Snapshot) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *RenderEventResponse_Snapshot) GetNoteCount() int32 {
	if x != nil && x.NoteCount != nil {
		return *x.NoteCount
	}
	return 0
}

func (x *RenderEventResponse_Snapshot) GetAddedActors() []string {
	if x != nil {
		return x.AddedActors
	}
	return nil
}

func (x *RenderEventResponse_Snapshot) GetRemovedActors() []string {
	if x != nil {
		return x.RemovedActors
	}
	return nil
}

type RenderEventResponse_PlaceNoteCountStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total *int32 `protobuf:"varint,1,opt,name=total" json:"total,omitempty"`
	Rank  *int32 `protobuf:"varint,2,opt,name=rank" json:"rank,omitempty"`
}

func (x *RenderEventResponse_PlaceNoteCountStats) Reset() {
	*x = RenderEventResponse_PlaceNoteCountStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RenderEventResponse_PlaceNoteCountStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RenderEventResponse_PlaceNoteCountStats) ProtoMessage() {}

func (x *RenderEventResponse_PlaceNoteCountStats) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RenderEventResponse_PlaceNoteCountStats.ProtoReflect.Descriptor instead.
func (*RenderEventResponse_PlaceNoteCountStats) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{1, 1}
}

func (x *RenderEventResponse_PlaceNoteCountStats) GetTotal() int32 {
	if x != nil && x.Total != nil {
		return *x.Total
	}
	return 0
}

func (x *RenderEventResponse_PlaceNoteCountStats) GetRank() int32 {
	if x != nil && x.Rank != nil {
		return *x.Rank
	}
	return 0
}

var File_service_proto protoreflect.FileDescriptor

var file_service_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0e, 0x74, 0x72, 0x65, 0x65, 0x64, 0x69, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x2e, 0x70, 0x62, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x24, 0x0a, 0x12, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0xb6, 0x04, 0x0a, 0x13, 0x52, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x65, 0x12, 0x4a, 0x0a, 0x09, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x74, 0x72, 0x65, 0x65, 0x64, 0x69, 0x61, 0x67,
	0x72, 0x61, 0x6d, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53, 0x6e, 0x61, 0x70, 0x73,
	0x68, 0x6f, 0x74, 0x52, 0x09, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x73, 0x12, 0x63,
	0x0a, 0x11, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x73, 0x5f, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x74, 0x72, 0x65, 0x65,
	0x64, 0x69, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x50,
	0x6c, 0x61, 0x63, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x74, 0x61,
	0x74, 0x73, 0x52, 0x0f, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x73, 0x54, 0x6f,
	0x74, 0x61, 0x6c, 0x12, 0x69, 0x0a, 0x14, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x73, 0x5f, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x37, 0x2e, 0x74, 0x72, 0x65, 0x65, 0x64, 0x69, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x2e,
	0x70, 0x62, 0x2e, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x4e, 0x6f, 0x74, 0x65,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x12, 0x70, 0x6c, 0x61, 0x63,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x73, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x1a, 0xad,
	0x01, 0x0a, 0x08, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x12, 0x38, 0x0a, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1d, 0x0a, 0x0a, 0x6e, 0x6f, 0x74, 0x65, 0x5f, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x6e, 0x6f, 0x74, 0x65, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x64, 0x64, 0x65, 0x64, 0x5f, 0x61, 0x63,
	0x74, 0x6f, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x64, 0x64, 0x65,
	0x64, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x72, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x64, 0x5f, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x0d, 0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x64, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x1a, 0x3f,
	0x0a, 0x13, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x4e, 0x6f, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x72,
	0x61, 0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x72, 0x61, 0x6e, 0x6b, 0x22,
	0x24, 0x0a, 0x12, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x41, 0x0a, 0x13, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x50,
	0x6c, 0x61, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x11,
	0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0f, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x32, 0xc8, 0x01, 0x0a, 0x12, 0x54, 0x72, 0x65,
	0x65, 0x44, 0x69, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x58, 0x0a, 0x0b, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x22,
	0x2e, 0x74, 0x72, 0x65, 0x65, 0x64, 0x69, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x2e, 0x70, 0x62, 0x2e,
	0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x23, 0x2e, 0x74, 0x72, 0x65, 0x65, 0x64, 0x69, 0x61, 0x67, 0x72, 0x61, 0x6d,
	0x2e, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x58, 0x0a, 0x0b, 0x52, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x12, 0x22, 0x2e, 0x74, 0x72, 0x65, 0x65, 0x64,
	0x69, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x50, 0x6c, 0x61, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x74,
	0x72, 0x65, 0x65, 0x64, 0x69, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x53, 0x53, 0x48, 0x5a, 0x2d, 0x4f, 0x52, 0x47, 0x2f, 0x74, 0x72, 0x65, 0x65, 0x2d,
	0x64, 0x69, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x2f, 0x70, 0x62,
}

var (
	file_service_proto_rawDescOnce sync.Once
	file_service_proto_rawDescData = file_service_proto_rawDesc
)

func file_service_proto_rawDescGZIP() []byte {
	file_service_proto_rawDescOnce.Do(func() {
		file_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_service_proto_rawDescData)
	})
	return file_service_proto_rawDescData
}

var file_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_service_proto_goTypes = []interface{}{
	(*RenderEventRequest)(nil),                      // 0: treediagram.pb.RenderEventRequest
	(*RenderEventResponse)(nil),                     // 1: treediagram.pb.RenderEventResponse
	(*RenderPlaceRequest)(nil),                      // 2: treediagram.pb.RenderPlaceRequest
	(*RenderPlaceResponse)(nil),                     // 3: treediagram.pb.RenderPlaceResponse
	(*RenderEventResponse_Snapshot)(nil),            // 4: treediagram.pb.RenderEventResponse.Snapshot
	(*RenderEventResponse_PlaceNoteCountStats)(nil), // 5: treediagram.pb.RenderEventResponse.PlaceNoteCountStats
	(*timestamppb.Timestamp)(nil),                   // 6: google.protobuf.Timestamp
}
var file_service_proto_depIdxs = []int32{
	4, // 0: treediagram.pb.RenderEventResponse.snapshots:type_name -> treediagram.pb.RenderEventResponse.Snapshot
	5, // 1: treediagram.pb.RenderEventResponse.place_stats_total:type_name -> treediagram.pb.RenderEventResponse.PlaceNoteCountStats
	5, // 2: treediagram.pb.RenderEventResponse.place_stats_finished:type_name -> treediagram.pb.RenderEventResponse.PlaceNoteCountStats
	6, // 3: treediagram.pb.RenderEventResponse.Snapshot.timestamp:type_name -> google.protobuf.Timestamp
	0, // 4: treediagram.pb.TreeDiagramService.RenderEvent:input_type -> treediagram.pb.RenderEventRequest
	2, // 5: treediagram.pb.TreeDiagramService.RenderPlace:input_type -> treediagram.pb.RenderPlaceRequest
	1, // 6: treediagram.pb.TreeDiagramService.RenderEvent:output_type -> treediagram.pb.RenderEventResponse
	3, // 7: treediagram.pb.TreeDiagramService.RenderPlace:output_type -> treediagram.pb.RenderPlaceResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_service_proto_init() }
func file_service_proto_init() {
	if File_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RenderEventRequest); i {
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
		file_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RenderEventResponse); i {
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
		file_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RenderPlaceRequest); i {
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
		file_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RenderPlaceResponse); i {
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
		file_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RenderEventResponse_Snapshot); i {
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
		file_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RenderEventResponse_PlaceNoteCountStats); i {
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
			RawDescriptor: file_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_proto_goTypes,
		DependencyIndexes: file_service_proto_depIdxs,
		MessageInfos:      file_service_proto_msgTypes,
	}.Build()
	File_service_proto = out.File
	file_service_proto_rawDesc = nil
	file_service_proto_goTypes = nil
	file_service_proto_depIdxs = nil
}
