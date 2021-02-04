package implementation

import (
	"context"
	"gateway/internal/domain"
	"gateway/internal/transport/grpc"
	pbmessenger "gateway/internal/transport/grpc/messenger"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
)

type grpcMessengerProxyService struct {
	messengerEndpoints *grpc.MessengerProxyEndpoints
}

func NewGRPCMessengerProxyService(messengerEndpoints *grpc.MessengerProxyEndpoints) *grpcMessengerProxyService {
	return &grpcMessengerProxyService{messengerEndpoints: messengerEndpoints}
}
func (g *grpcMessengerProxyService) CreateChat(ctx context.Context, userToken, companionID string) (string, error) {
	response, err := g.messengerEndpoints.CreateChat(ctx, &pbmessenger.CreateChatRequest{
		MasterToken: userToken,
		SlaveId:     companionID,
	})
	if err != nil {
		return "", err
	}

	resp := response.(*pbmessenger.CreateChatResponse)

	return resp.ChatId, nil
}

func (g *grpcMessengerProxyService) GetChats(ctx context.Context, userToken string, offset, limit *int32) (*domain.GetChatsResponse, error) {
	response, err := g.messengerEndpoints.GetChats(ctx, &pbmessenger.GetChatsRequest{
		UserToken: userToken,
		Offset:    &wrappers.Int32Value{Value: *offset},
		Limit:     &wrappers.Int32Value{Value: *limit},
	})
	if err != nil {
		return nil, err
	}

	resp := response.(*pbmessenger.GetChatsResponse)

	chats := make([]*domain.Chat, 0, len(resp.Chats))
	for _, c := range resp.Chats {
		participants := make([]*domain.Participant, 0, len(c.Participants))

		for _, p := range participants {
			participants = append(participants, &domain.Participant{
				ID:      p.ID,
				Name:    p.Name,
				Surname: p.Surname,
			})
		}

		chats = append(chats, &domain.Chat{
			ID:           c.Id,
			CreateTime:   time.Unix(0, c.CreateTime),
			Participants: participants,
		})
	}

	return &domain.GetChatsResponse{
		Total:  resp.Total,
		Limit:  resp.Limit,
		Offset: resp.Offset,
		Chats:  chats,
	}, nil
}

func (g *grpcMessengerProxyService) GetMessages(ctx context.Context, userToken, chatID string, offset, limit *int32) (*domain.GetMessagesResponse, error) {
	response, err := g.messengerEndpoints.GetMessages(ctx, &pbmessenger.GetMessagesRequest{
		UserToken: userToken,
		ChatId:    chatID,
		Offset:    &wrappers.Int32Value{Value: *offset},
		Limit:     &wrappers.Int32Value{Value: *limit},
	})
	if err != nil {
		return nil, err
	}

	resp := response.(*pbmessenger.GetMessagesResponse)

	messages := make([]*domain.Message, 0, len(resp.Messages))
	for _, m := range resp.Messages {
		messages = append(messages, &domain.Message{
			ID:         m.Id,
			Text:       m.Text,
			Status:     m.Status,
			CreateTime: time.Unix(0, m.CreateTime),
			UserID:     m.UserId,
			ChatID:     m.ChatId,
		})
	}

	return &domain.GetMessagesResponse{
		Total:    resp.Total,
		Limit:    resp.Limit,
		Offset:   resp.Offset,
		Messages: messages,
	}, nil
}
