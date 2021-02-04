package grpc

import (
	pbmessenger "gateway/internal/transport/grpc/messenger"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type MessengerProxyEndpoints struct {
	CreateChat  endpoint.Endpoint
	GetChats    endpoint.Endpoint
	GetMessages endpoint.Endpoint
}

func MakeMessengerProxyEndpoints(conn *grpc.ClientConn) *MessengerProxyEndpoints {
	return &MessengerProxyEndpoints{
		CreateChat:  makeCreateChatProxyEndpoint(conn),
		GetChats:    makeGetChatsProxyEndpoint(conn),
		GetMessages: makeGetMessagesProxyEndpoint(conn),
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
	).Endpoint()
}
