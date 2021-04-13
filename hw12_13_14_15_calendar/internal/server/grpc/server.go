package internalgrpc

//go:generate protoc -I ../../../api EventService.proto --go_out=. --go-grpc_out=require_unimplemented_servers=false:.
import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/logger"
	"github.com/Haba1234/hw-test/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	srv  *grpc.Server
	app  *app.App
	logg *logger.Logger
}

//nolint: golint
func NewServer(logg *logger.Logger, app *app.App) *server {
	return &server{
		app:  app,
		logg: logg,
	}
}

func (s *server) Start(ctx context.Context, addr string) error {
	s.logg.Info("gRPC server " + addr + " starting...")
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.srv = grpc.NewServer(grpc.UnaryInterceptor(loggingServerInterceptor(s.logg)))
	if err := s.srv.Serve(lsn); err != nil {
		return err
	}

	return nil
}

func loggingServerInterceptor(logger app.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		logger.Info(fmt.Sprintf("method: %s, duration: %s, request: %+v", info.FullMethod, time.Since(time.Now()), req))
		h, err := handler(ctx, req)
		return h, err
	}
}

func (s *server) Stop(ctx context.Context) error {
	s.logg.Info("gRPC server stopping...")
	s.srv.GracefulStop()
	return nil
}

func (s server) CreateEvent(ctx context.Context, event *Event) (*CreateEventResponse, error) {
	userID, err := uuid.Parse(event.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("error to create event: %v", err))
	}

	e := app.Event{
		Title:        event.Title,
		DateTime:     event.StartDate.AsTime(),
		Duration:     event.Duration.AsDuration(),
		Description:  event.Description,
		UserID:       userID,
		NotifyBefore: event.NotifyBefore,
	}

	id, err := s.app.CreateEvent(ctx, &e)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error to create event: %v", err))
	}

	return &CreateEventResponse{Id: id.String()}, nil
}

func (s server) UpdateEvent(ctx context.Context, event *Event) (*emptypb.Empty, error) {
	ID, err := uuid.Parse(event.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("error to update event: %v", err))
	}

	userID, err := uuid.Parse(event.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("error to update event: %v", err))
	}

	e := app.Event{
		ID:           ID,
		Title:        event.Title,
		DateTime:     event.StartDate.AsTime(),
		Duration:     event.Duration.AsDuration(),
		Description:  event.Description,
		UserID:       userID,
		NotifyBefore: event.NotifyBefore,
	}

	err = s.app.UpdateEvent(ctx, &e)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error to update event: %v", err))
	}

	return &emptypb.Empty{}, nil
}

func (s server) DeleteEvent(ctx context.Context, req *DeleteEventRequest) (*emptypb.Empty, error) {
	ID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("error to delete event: %v", err))
	}

	err = s.app.DeleteEvent(ctx, ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error to delete event: %v", err))
	}

	return &emptypb.Empty{}, nil
}

func (s server) GetListEvents(ctx context.Context) (*GetListEventsResponse, error) {
	listEvents, err := s.app.GetListEvents(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error to get evens list: %v", err))
	}

	return &GetListEventsResponse{Events: convertStorageEvToGrpcEv(listEvents)}, nil
}

func (s server) GetListEventsDay(ctx context.Context, req *GetListEventsDayRequest) (*GetListEventsDayResponse, error) {
	listEvents, err := s.app.GetListEventsDay(ctx, req.StartDate.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error to get evens list day: %v", err))
	}

	return &GetListEventsDayResponse{Events: convertStorageEvToGrpcEv(listEvents)}, nil
}

func (s server) GetListEventsWeek(ctx context.Context, req *GetListEventsWeekRequest) (*GetListEventsWeekResponse, error) {
	listEvents, err := s.app.GetListEventsWeek(ctx, req.StartDate.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error to get evens list week: %v", err))
	}

	return &GetListEventsWeekResponse{Events: convertStorageEvToGrpcEv(listEvents)}, nil
}

func (s server) GetListEventsMonth(ctx context.Context, req *GetListEventsMonthRequest) (*GetListEventsMonthResponse, error) {
	listEvents, err := s.app.GetListEventsMonth(ctx, req.StartDate.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error to get evens list month: %v", err))
	}

	return &GetListEventsMonthResponse{Events: convertStorageEvToGrpcEv(listEvents)}, nil
}

func convertStorageEvToGrpcEv(events []storage.Event) []*Event {
	resultEvents := make([]*Event, 0, len(events))
	for _, event := range events {
		resultEvent := &Event{
			Id:           event.ID.String(),
			Title:        event.Title,
			StartDate:    timestamppb.New(event.DateTime),
			Duration:     durationpb.New(event.Duration),
			Description:  event.Description,
			UserId:       event.UserID.String(),
			NotifyBefore: event.NotifyBefore,
		}
		resultEvents = append(resultEvents, resultEvent)
	}
	return resultEvents
}
