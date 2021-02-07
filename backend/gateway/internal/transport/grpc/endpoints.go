package grpc

import (
	pbmessenger "gateway/internal/transport/grpc/messenger"
	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/endpoint"
	kitopentracing "github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type MessengerProxyEndpoints struct {
	CreateChat  endpoint.Endpoint
	GetChats    endpoint.Endpoint
	GetMessages endpoint.Endpoint
}

func MakeMessengerProxyEndpoints(conn *grpc.ClientConn) *MessengerProxyEndpoints {
	return &MessengerProxyEndpoints{
		CreateChat:  kitopentracing.TraceClient(opentracing.GlobalTracer(), "gateway")(makeCreateChatProxyEndpoint(conn)),
		GetChats:    kitopentracing.TraceClient(opentracing.GlobalTracer(), "gateway")(makeGetChatsProxyEndpoint(conn)),
		GetMessages: kitopentracing.TraceClient(opentracing.GlobalTracer(), "gateway")(makeGetMessagesProxyEndpoint(conn)),
	}
}

func makeCreateChatProxyEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"messenger.Messenger",
		"CreateChat",
		encodeCreateChatProxyRequest,
		decodeCreateChatProxyResponse,
		pbmessenger.CreateChatResponse{},
		grpctransport.ClientBefore(kitopentracing.ContextToGRPC(opentracing.GlobalTracer(), log.NewNopLogger())),
	).Endpoint()
}

func makeGetChatsProxyEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"messenger.Messenger",
		"GetChats",
		encodeGetChatsProxyRequest,
		decodeGetChatsProxyResponse,
		pbmessenger.GetChatsResponse{},
		grpctransport.ClientBefore(kitopentracing.ContextToGRPC(opentracing.GlobalTracer(), log.NewNopLogger())),
	).Endpoint()
}

func makeGetMessagesProxyEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"messenger.Messenger",
		"GetMessages",
		encodeGetMessagesProxyRequest,
		decodeGetMessagesResponse,
		pbmessenger.GetMessagesResponse{},
		grpctransport.ClientBefore(kitopentracing.ContextToGRPC(opentracing.GlobalTracer(), log.NewNopLogger())),
	).Endpoint()
}
