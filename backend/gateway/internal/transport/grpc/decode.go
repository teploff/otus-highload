package grpc

import (
	"context"

	pbmessenger "gateway/internal/transport/grpc/messenger"
)

func decodeCreateChatProxyResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response.(*pbmessenger.CreateChatResponse), nil
}

func decodeGetChatsProxyResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response.(*pbmessenger.GetChatsResponse), nil
}

func decodeGetMessagesResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response.(*pbmessenger.GetMessagesResponse), nil
}
