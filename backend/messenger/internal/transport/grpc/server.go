package grpc

import (
	"context"
	pb "messenger/internal/transport/grpc/messenger"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxReceivedMsgSize = 1024 * 1024 * 20

type server struct {
	createChat kitgrpc.Handler
}

func (s *server) CreateChat(ctx context.Context, request *pb.CreateChatRequest) (*pb.CreateChatResponse, error) {
	_, response, err := s.createChat.ServeGRPC(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response.(*pb.CreateChatResponse), nil
}

func (s *server) GetChats(ctx context.Context, request *pb.GetChatsRequest) (*pb.GetChatsResponse, error) {
	panic("implement me")
}

func (s *server) GetMessages(ctx context.Context, request *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	panic("implement me")
}

// NewGRPCServer instance of gRPC server.
func NewGRPCServer(endpoints *Endpoints, errLogger log.Logger) *grpc.Server {
	options := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorHandler(transport.NewLogErrorHandler(errLogger)),
	}

	srv := &server{
		createChat: newRecoveryGRPCHandler(kitgrpc.NewServer(
			endpoints.CreateChat,
			decodeCreateChatRequest,
			encodeSignInResponse,
			options...,
		), errLogger),
	}

	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor), grpc.MaxRecvMsgSize(maxReceivedMsgSize))
	pb.RegisterMessengerServer(baseServer, srv)

	return baseServer
}

func decodeCreateChatRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	request := grpcReq.(*pb.CreateChatRequest)

	return &CreateChatRequest{
		MasterToken: request.MasterToken,
		SlaveID:     request.SlaveId,
	}, nil
}

func encodeSignInResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	response := grpcResp.(CreateChatResponse)

	return &pb.CreateChatResponse{
		ChatId: response.ChatID,
	}, nil
}

//recoveryGRPCHandler wrap gRPC server, recover them if panic was fired.
type recoveryGRPCHandler struct {
	next   kitgrpc.Handler
	logger log.Logger
}

func newRecoveryGRPCHandler(next kitgrpc.Handler, logger log.Logger) *recoveryGRPCHandler {
	return &recoveryGRPCHandler{next: next, logger: logger}
}

func (rh *recoveryGRPCHandler) ServeGRPC(ctx context.Context, req interface{}) (context.Context, interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				_ = rh.logger.Log("msg", "gRPC server panic recover", "text", err.Error())
			}
		}
	}()

	return rh.next.ServeGRPC(ctx, req)
}
