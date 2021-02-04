package grpc

import (
	"context"
	"messenger/internal/domain"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints for messenger.
type Endpoints struct {
	CreateChat endpoint.Endpoint
}

func makeCreateChatEndpoint(svc domain.MessengerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*CreateChatRequest)
		chatID, err := svc.CreateChat(ctx, req.MasterToken, req.SlaveID)

		if err != nil {
			return nil, err
		}

		return &CreateChatResponse{ChatID: chatID}, nil
	}
}

// MakeMessengerEndpoints for messenger functionality.
func MakeMessengerEndpoints(svc domain.MessengerService) *Endpoints {
	return &Endpoints{
		CreateChat: makeCreateChatEndpoint(svc),
	}
}
