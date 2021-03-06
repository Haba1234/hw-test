// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package internalgrpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CalendarClient is the client API for Calendar service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CalendarClient interface {
	CreateEvent(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*CreateEventResponse, error)
	UpdateEvent(ctx context.Context, in *EventUpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteEvent(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetListEvents(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetListEventsResponse, error)
	GetListEventsDay(ctx context.Context, in *GetListEventsDayRequest, opts ...grpc.CallOption) (*GetListEventsDayResponse, error)
	GetListEventsWeek(ctx context.Context, in *GetListEventsWeekRequest, opts ...grpc.CallOption) (*GetListEventsWeekResponse, error)
	GetListEventsMonth(ctx context.Context, in *GetListEventsMonthRequest, opts ...grpc.CallOption) (*GetListEventsMonthResponse, error)
}

type calendarClient struct {
	cc grpc.ClientConnInterface
}

func NewCalendarClient(cc grpc.ClientConnInterface) CalendarClient {
	return &calendarClient{cc}
}

func (c *calendarClient) CreateEvent(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*CreateEventResponse, error) {
	out := new(CreateEventResponse)
	err := c.cc.Invoke(ctx, "/event.Calendar/CreateEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) UpdateEvent(ctx context.Context, in *EventUpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/event.Calendar/UpdateEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) DeleteEvent(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/event.Calendar/DeleteEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) GetListEvents(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetListEventsResponse, error) {
	out := new(GetListEventsResponse)
	err := c.cc.Invoke(ctx, "/event.Calendar/GetListEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) GetListEventsDay(ctx context.Context, in *GetListEventsDayRequest, opts ...grpc.CallOption) (*GetListEventsDayResponse, error) {
	out := new(GetListEventsDayResponse)
	err := c.cc.Invoke(ctx, "/event.Calendar/GetListEventsDay", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) GetListEventsWeek(ctx context.Context, in *GetListEventsWeekRequest, opts ...grpc.CallOption) (*GetListEventsWeekResponse, error) {
	out := new(GetListEventsWeekResponse)
	err := c.cc.Invoke(ctx, "/event.Calendar/GetListEventsWeek", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calendarClient) GetListEventsMonth(ctx context.Context, in *GetListEventsMonthRequest, opts ...grpc.CallOption) (*GetListEventsMonthResponse, error) {
	out := new(GetListEventsMonthResponse)
	err := c.cc.Invoke(ctx, "/event.Calendar/GetListEventsMonth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalendarServer is the server API for Calendar service.
// All implementations should embed UnimplementedCalendarServer
// for forward compatibility
type CalendarServer interface {
	CreateEvent(context.Context, *CreateEventRequest) (*CreateEventResponse, error)
	UpdateEvent(context.Context, *EventUpdateRequest) (*emptypb.Empty, error)
	DeleteEvent(context.Context, *DeleteEventRequest) (*emptypb.Empty, error)
	GetListEvents(context.Context, *emptypb.Empty) (*GetListEventsResponse, error)
	GetListEventsDay(context.Context, *GetListEventsDayRequest) (*GetListEventsDayResponse, error)
	GetListEventsWeek(context.Context, *GetListEventsWeekRequest) (*GetListEventsWeekResponse, error)
	GetListEventsMonth(context.Context, *GetListEventsMonthRequest) (*GetListEventsMonthResponse, error)
}

// UnimplementedCalendarServer should be embedded to have forward compatible implementations.
type UnimplementedCalendarServer struct {
}

func (UnimplementedCalendarServer) CreateEvent(context.Context, *CreateEventRequest) (*CreateEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEvent not implemented")
}
func (UnimplementedCalendarServer) UpdateEvent(context.Context, *EventUpdateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEvent not implemented")
}
func (UnimplementedCalendarServer) DeleteEvent(context.Context, *DeleteEventRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEvent not implemented")
}
func (UnimplementedCalendarServer) GetListEvents(context.Context, *emptypb.Empty) (*GetListEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListEvents not implemented")
}
func (UnimplementedCalendarServer) GetListEventsDay(context.Context, *GetListEventsDayRequest) (*GetListEventsDayResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListEventsDay not implemented")
}
func (UnimplementedCalendarServer) GetListEventsWeek(context.Context, *GetListEventsWeekRequest) (*GetListEventsWeekResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListEventsWeek not implemented")
}
func (UnimplementedCalendarServer) GetListEventsMonth(context.Context, *GetListEventsMonthRequest) (*GetListEventsMonthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListEventsMonth not implemented")
}

// UnsafeCalendarServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CalendarServer will
// result in compilation errors.
type UnsafeCalendarServer interface {
	mustEmbedUnimplementedCalendarServer()
}

func RegisterCalendarServer(s grpc.ServiceRegistrar, srv CalendarServer) {
	s.RegisterService(&Calendar_ServiceDesc, srv)
}

func _Calendar_CreateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).CreateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Calendar/CreateEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).CreateEvent(ctx, req.(*CreateEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_UpdateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).UpdateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Calendar/UpdateEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).UpdateEvent(ctx, req.(*EventUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_DeleteEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).DeleteEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Calendar/DeleteEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).DeleteEvent(ctx, req.(*DeleteEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_GetListEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).GetListEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Calendar/GetListEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).GetListEvents(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_GetListEventsDay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListEventsDayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).GetListEventsDay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Calendar/GetListEventsDay",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).GetListEventsDay(ctx, req.(*GetListEventsDayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_GetListEventsWeek_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListEventsWeekRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).GetListEventsWeek(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Calendar/GetListEventsWeek",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).GetListEventsWeek(ctx, req.(*GetListEventsWeekRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calendar_GetListEventsMonth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListEventsMonthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalendarServer).GetListEventsMonth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Calendar/GetListEventsMonth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalendarServer).GetListEventsMonth(ctx, req.(*GetListEventsMonthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Calendar_ServiceDesc is the grpc.ServiceDesc for Calendar service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Calendar_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "event.Calendar",
	HandlerType: (*CalendarServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateEvent",
			Handler:    _Calendar_CreateEvent_Handler,
		},
		{
			MethodName: "UpdateEvent",
			Handler:    _Calendar_UpdateEvent_Handler,
		},
		{
			MethodName: "DeleteEvent",
			Handler:    _Calendar_DeleteEvent_Handler,
		},
		{
			MethodName: "GetListEvents",
			Handler:    _Calendar_GetListEvents_Handler,
		},
		{
			MethodName: "GetListEventsDay",
			Handler:    _Calendar_GetListEventsDay_Handler,
		},
		{
			MethodName: "GetListEventsWeek",
			Handler:    _Calendar_GetListEventsWeek_Handler,
		},
		{
			MethodName: "GetListEventsMonth",
			Handler:    _Calendar_GetListEventsMonth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "EventService.proto",
}
