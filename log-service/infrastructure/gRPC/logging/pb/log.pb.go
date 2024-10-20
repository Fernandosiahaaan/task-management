// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.2
// source: proto/log.proto

package pb

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

// Enum to define the various actions for task and user logs
type TaskAction int32

const (
	TaskAction_CREATE_TASK TaskAction = 0
	TaskAction_UPDATE_TASK TaskAction = 1
	TaskAction_DELETE_TASK TaskAction = 2
)

// Enum value maps for TaskAction.
var (
	TaskAction_name = map[int32]string{
		0: "CREATE_TASK",
		1: "UPDATE_TASK",
		2: "DELETE_TASK",
	}
	TaskAction_value = map[string]int32{
		"CREATE_TASK": 0,
		"UPDATE_TASK": 1,
		"DELETE_TASK": 2,
	}
)

func (x TaskAction) Enum() *TaskAction {
	p := new(TaskAction)
	*p = x
	return p
}

func (x TaskAction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TaskAction) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_log_proto_enumTypes[0].Descriptor()
}

func (TaskAction) Type() protoreflect.EnumType {
	return &file_proto_log_proto_enumTypes[0]
}

func (x TaskAction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TaskAction.Descriptor instead.
func (TaskAction) EnumDescriptor() ([]byte, []int) {
	return file_proto_log_proto_rawDescGZIP(), []int{0}
}

type UserAction int32

const (
	UserAction_LOGIN       UserAction = 0
	UserAction_LOGOUT      UserAction = 1
	UserAction_CREATE_USER UserAction = 2
	UserAction_UPDATE_USER UserAction = 3
)

// Enum value maps for UserAction.
var (
	UserAction_name = map[int32]string{
		0: "LOGIN",
		1: "LOGOUT",
		2: "CREATE_USER",
		3: "UPDATE_USER",
	}
	UserAction_value = map[string]int32{
		"LOGIN":       0,
		"LOGOUT":      1,
		"CREATE_USER": 2,
		"UPDATE_USER": 3,
	}
)

func (x UserAction) Enum() *UserAction {
	p := new(UserAction)
	*p = x
	return p
}

func (x UserAction) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserAction) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_log_proto_enumTypes[1].Descriptor()
}

func (UserAction) Type() protoreflect.EnumType {
	return &file_proto_log_proto_enumTypes[1]
}

func (x UserAction) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserAction.Descriptor instead.
func (UserAction) EnumDescriptor() ([]byte, []int) {
	return file_proto_log_proto_rawDescGZIP(), []int{1}
}

// Request message for logging task-related events
type LogTaskRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId    string       `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`                             // ID of the user performing the action
	TaskId    int64        `protobuf:"varint,2,opt,name=taskId,proto3" json:"taskId,omitempty"`                            // ID of the task being affected
	Action    TaskAction   `protobuf:"varint,3,opt,name=action,proto3,enum=logservice.TaskAction" json:"action,omitempty"` // The action being logged (create, update, delete)
	Timestamp string       `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`                       // Timestamp of when the action happened
	Before    *TaskDetails `protobuf:"bytes,5,opt,name=before,proto3" json:"before,omitempty"`                             // Optional: task details before the change (for update)
	After     *TaskDetails `protobuf:"bytes,6,opt,name=after,proto3" json:"after,omitempty"`                               // Optional: task details after the change (for update)
}

func (x *LogTaskRequest) Reset() {
	*x = LogTaskRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_log_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogTaskRequest) ProtoMessage() {}

func (x *LogTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_log_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogTaskRequest.ProtoReflect.Descriptor instead.
func (*LogTaskRequest) Descriptor() ([]byte, []int) {
	return file_proto_log_proto_rawDescGZIP(), []int{0}
}

func (x *LogTaskRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *LogTaskRequest) GetTaskId() int64 {
	if x != nil {
		return x.TaskId
	}
	return 0
}

func (x *LogTaskRequest) GetAction() TaskAction {
	if x != nil {
		return x.Action
	}
	return TaskAction_CREATE_TASK
}

func (x *LogTaskRequest) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *LogTaskRequest) GetBefore() *TaskDetails {
	if x != nil {
		return x.Before
	}
	return nil
}

func (x *LogTaskRequest) GetAfter() *TaskDetails {
	if x != nil {
		return x.After
	}
	return nil
}

// Message to log details of a task
type TaskDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title       string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	DueDate     string `protobuf:"bytes,3,opt,name=dueDate,proto3" json:"dueDate,omitempty"` // Optional due date for the task
	Status      string `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`   // Status of the task (e.g., pending, completed)
}

func (x *TaskDetails) Reset() {
	*x = TaskDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_log_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskDetails) ProtoMessage() {}

func (x *TaskDetails) ProtoReflect() protoreflect.Message {
	mi := &file_proto_log_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskDetails.ProtoReflect.Descriptor instead.
func (*TaskDetails) Descriptor() ([]byte, []int) {
	return file_proto_log_proto_rawDescGZIP(), []int{1}
}

func (x *TaskDetails) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *TaskDetails) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *TaskDetails) GetDueDate() string {
	if x != nil {
		return x.DueDate
	}
	return ""
}

func (x *TaskDetails) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

// Request message for logging user-related events
type LogUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId    string       `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`                             // ID of the user being affected
	Action    UserAction   `protobuf:"varint,2,opt,name=action,proto3,enum=logservice.UserAction" json:"action,omitempty"` // The action being logged (login, logout, create, update)
	Timestamp string       `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`                       // Timestamp of when the action happened
	Before    *UserDetails `protobuf:"bytes,4,opt,name=before,proto3" json:"before,omitempty"`                             // Optional: user details before the change (for update)
	After     *UserDetails `protobuf:"bytes,5,opt,name=after,proto3" json:"after,omitempty"`                               // Optional: user details after the change (for update)
}

func (x *LogUserRequest) Reset() {
	*x = LogUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_log_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogUserRequest) ProtoMessage() {}

func (x *LogUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_log_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogUserRequest.ProtoReflect.Descriptor instead.
func (*LogUserRequest) Descriptor() ([]byte, []int) {
	return file_proto_log_proto_rawDescGZIP(), []int{2}
}

func (x *LogUserRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *LogUserRequest) GetAction() UserAction {
	if x != nil {
		return x.Action
	}
	return UserAction_LOGIN
}

func (x *LogUserRequest) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *LogUserRequest) GetBefore() *UserDetails {
	if x != nil {
		return x.Before
	}
	return nil
}

func (x *LogUserRequest) GetAfter() *UserDetails {
	if x != nil {
		return x.After
	}
	return nil
}

// Message to log details of a user
type UserDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   string `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`     // ID of the user
	Email    string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`       // User email
	Username string `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"` // Username
	Role     string `protobuf:"bytes,4,opt,name=role,proto3" json:"role,omitempty"`         // role of the user
}

func (x *UserDetails) Reset() {
	*x = UserDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_log_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserDetails) ProtoMessage() {}

func (x *UserDetails) ProtoReflect() protoreflect.Message {
	mi := &file_proto_log_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserDetails.ProtoReflect.Descriptor instead.
func (*UserDetails) Descriptor() ([]byte, []int) {
	return file_proto_log_proto_rawDescGZIP(), []int{3}
}

func (x *UserDetails) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserDetails) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UserDetails) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *UserDetails) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

// Response message for both task and user logs
type LogResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"` // Whether the log was successfully recorded
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`  // Additional info or error message
}

func (x *LogResponse) Reset() {
	*x = LogResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_log_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogResponse) ProtoMessage() {}

func (x *LogResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_log_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogResponse.ProtoReflect.Descriptor instead.
func (*LogResponse) Descriptor() ([]byte, []int) {
	return file_proto_log_proto_rawDescGZIP(), []int{4}
}

func (x *LogResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *LogResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_proto_log_proto protoreflect.FileDescriptor

var file_proto_log_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0a, 0x6c, 0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0xee, 0x01,
	0x0a, 0x0e, 0x4c, 0x6f, 0x67, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x73, 0x6b,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64,
	0x12, 0x2e, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x16, 0x2e, 0x6c, 0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x61,
	0x73, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x2f,
	0x0a, 0x06, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x6c, 0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x61, 0x73, 0x6b,
	0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x06, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x12,
	0x2d, 0x0a, 0x05, 0x61, 0x66, 0x74, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x6c, 0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x54, 0x61, 0x73, 0x6b,
	0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x05, 0x61, 0x66, 0x74, 0x65, 0x72, 0x22, 0x77,
	0x0a, 0x0b, 0x54, 0x61, 0x73, 0x6b, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x75, 0x65, 0x44, 0x61, 0x74, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x75, 0x65, 0x44, 0x61, 0x74, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0xd6, 0x01, 0x0a, 0x0e, 0x4c, 0x6f, 0x67, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x2e, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x16, 0x2e, 0x6c, 0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x12, 0x2f, 0x0a, 0x06, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x6c, 0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x06, 0x62, 0x65, 0x66, 0x6f, 0x72,
	0x65, 0x12, 0x2d, 0x0a, 0x05, 0x61, 0x66, 0x74, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x6c, 0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x05, 0x61, 0x66, 0x74, 0x65, 0x72,
	0x22, 0x6b, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12,
	0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a,
	0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x22, 0x41, 0x0a,
	0x0b, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2a, 0x3f, 0x0a, 0x0a, 0x54, 0x61, 0x73, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0f,
	0x0a, 0x0b, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x10, 0x00, 0x12,
	0x0f, 0x0a, 0x0b, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x10, 0x01,
	0x12, 0x0f, 0x0a, 0x0b, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x5f, 0x54, 0x41, 0x53, 0x4b, 0x10,
	0x02, 0x2a, 0x45, 0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x09, 0x0a, 0x05, 0x4c, 0x4f, 0x47, 0x49, 0x4e, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x4c, 0x4f,
	0x47, 0x4f, 0x55, 0x54, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45,
	0x5f, 0x55, 0x53, 0x45, 0x52, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x50, 0x44, 0x41, 0x54,
	0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x10, 0x03, 0x32, 0x98, 0x01, 0x0a, 0x0a, 0x4c, 0x6f, 0x67,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x44, 0x0a, 0x0d, 0x4c, 0x6f, 0x67, 0x54, 0x61,
	0x73, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x2e, 0x6c, 0x6f, 0x67, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x67, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x6c, 0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a,
	0x0d, 0x4c, 0x6f, 0x67, 0x55, 0x73, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a,
	0x2e, 0x6c, 0x6f, 0x67, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x67, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x6c, 0x6f, 0x67,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x2f, 0x6c, 0x6f, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_log_proto_rawDescOnce sync.Once
	file_proto_log_proto_rawDescData = file_proto_log_proto_rawDesc
)

func file_proto_log_proto_rawDescGZIP() []byte {
	file_proto_log_proto_rawDescOnce.Do(func() {
		file_proto_log_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_log_proto_rawDescData)
	})
	return file_proto_log_proto_rawDescData
}

var file_proto_log_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_proto_log_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_log_proto_goTypes = []any{
	(TaskAction)(0),        // 0: logservice.TaskAction
	(UserAction)(0),        // 1: logservice.UserAction
	(*LogTaskRequest)(nil), // 2: logservice.LogTaskRequest
	(*TaskDetails)(nil),    // 3: logservice.TaskDetails
	(*LogUserRequest)(nil), // 4: logservice.LogUserRequest
	(*UserDetails)(nil),    // 5: logservice.UserDetails
	(*LogResponse)(nil),    // 6: logservice.LogResponse
}
var file_proto_log_proto_depIdxs = []int32{
	0, // 0: logservice.LogTaskRequest.action:type_name -> logservice.TaskAction
	3, // 1: logservice.LogTaskRequest.before:type_name -> logservice.TaskDetails
	3, // 2: logservice.LogTaskRequest.after:type_name -> logservice.TaskDetails
	1, // 3: logservice.LogUserRequest.action:type_name -> logservice.UserAction
	5, // 4: logservice.LogUserRequest.before:type_name -> logservice.UserDetails
	5, // 5: logservice.LogUserRequest.after:type_name -> logservice.UserDetails
	2, // 6: logservice.LogService.LogTaskAction:input_type -> logservice.LogTaskRequest
	4, // 7: logservice.LogService.LogUserAction:input_type -> logservice.LogUserRequest
	6, // 8: logservice.LogService.LogTaskAction:output_type -> logservice.LogResponse
	6, // 9: logservice.LogService.LogUserAction:output_type -> logservice.LogResponse
	8, // [8:10] is the sub-list for method output_type
	6, // [6:8] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_proto_log_proto_init() }
func file_proto_log_proto_init() {
	if File_proto_log_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_log_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*LogTaskRequest); i {
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
		file_proto_log_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*TaskDetails); i {
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
		file_proto_log_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*LogUserRequest); i {
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
		file_proto_log_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*UserDetails); i {
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
		file_proto_log_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*LogResponse); i {
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
			RawDescriptor: file_proto_log_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_log_proto_goTypes,
		DependencyIndexes: file_proto_log_proto_depIdxs,
		EnumInfos:         file_proto_log_proto_enumTypes,
		MessageInfos:      file_proto_log_proto_msgTypes,
	}.Build()
	File_proto_log_proto = out.File
	file_proto_log_proto_rawDesc = nil
	file_proto_log_proto_goTypes = nil
	file_proto_log_proto_depIdxs = nil
}
