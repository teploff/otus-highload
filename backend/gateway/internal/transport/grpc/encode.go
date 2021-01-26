package grpc

import (
	"context"

	pbmessenger "gateway/internal/transport/grpc/messenger"
)

func encodeGetChatsProxyRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request.(*pbmessenger.GetChatsRequest), nil
}

func encodeGetMessagesProxyRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request.(*pbmessenger.GetMessagesRequest), nil
}
