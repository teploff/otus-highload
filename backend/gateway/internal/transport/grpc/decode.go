package grpc

import (
	"context"

	pbmessenger "gateway/internal/transport/grpc/messenger"
)

func decodeGetChatsProxyResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response.(*pbmessenger.GetChatsResponse), nil
}

func decodeGetMessagesResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response.(*pbmessenger.GetMessagesResponse), nil
}
