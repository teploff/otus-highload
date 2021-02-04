package grpc

import (
	"context"

	pbmessenger "gateway/internal/transport/grpc/messenger"
)

func encodeCreateChatProxyRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request.(*pbmessenger.CreateChatRequest), nil
}

func encodeGetChatsProxyRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request.(*pbmessenger.GetChatsRequest), nil
}

func encodeGetMessagesProxyRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request.(*pbmessenger.GetMessagesRequest), nil
}
