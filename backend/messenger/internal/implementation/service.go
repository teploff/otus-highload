package implementation

import (
	"context"
	"messenger/internal/domain"
	"messenger/internal/transport/http"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"go.uber.org/zap"
)

type authService struct {
	endpoints *http.AuthProxyEndpoints
}

func NewAuthService(endpoints *http.AuthProxyEndpoints) *authService {
	return &authService{
		endpoints: endpoints,
	}
}

func (a *authService) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	request := http.GetUserIDByAccessTokenRequest{Token: token}

	resp, err := a.endpoints.GetUserIDByAccessToken(ctx, request)
	if err != nil {
		return nil, err
	}

	response := resp.(http.GetUserIDByAccessTokenResponse)

	return response.User, nil
}

type messengerService struct {
	authSvc domain.AuthService
	messRep domain.MessengerRepository
}

func NewMessengerService(authSvc domain.AuthService, messengerRep domain.MessengerRepository) *messengerService {
	return &messengerService{
		authSvc: authSvc,
		messRep: messengerRep,
	}
}

func (m *messengerService) CreateChat(ctx context.Context, masterToken, slaveID string) (string, error) {
	user, err := m.authSvc.GetUserByToken(ctx, masterToken)
	if err != nil {
		return "", err
	}

	chatID, err := m.messRep.CreateChat(user.ID, slaveID)
	if err != nil {
		return "", err
	}

	return chatID, nil
}

func (m *messengerService) GetChat(ctx context.Context, masterToken, slaveID string) (*domain.Chat, error) {
	user, err := m.authSvc.GetUserByToken(ctx, masterToken)
	if err != nil {
		return nil, err
	}

	chat, err := m.messRep.GetChatWithCompanion(user.ID, slaveID)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (m *messengerService) GetChats(ctx context.Context, userToken string, offset, limit int) ([]*domain.Chat, int, error) {
	user, err := m.authSvc.GetUserByToken(ctx, userToken)
	if err != nil {
		return nil, 0, err
	}

	total, err := m.messRep.GetCountChats(user.ID)
	if err != nil {
		return nil, 0, err
	}

	chats, err := m.messRep.GetChats(user.ID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return chats, total, nil
}

func (m *messengerService) SendMessages(ctx context.Context, userToken, chatID string, messages []*domain.ShortMessage) error {
	user, err := m.authSvc.GetUserByToken(ctx, userToken)
	if err != nil {
		return err
	}

	return m.messRep.PersistMessages(user.ID, chatID, messages)
}

func (m *messengerService) GetMessages(ctx context.Context, userToken, chatID string, offset, limit int) ([]*domain.Message, int, error) {
	user, err := m.authSvc.GetUserByToken(ctx, userToken)
	if err != nil {
		return nil, 0, err
	}

	_, err = m.messRep.GetChatAsParticipant(user.ID)
	if err != nil {
		return nil, 0, err
	}

	total, err := m.messRep.GetCountMessages(chatID)
	if err != nil {
		return nil, 0, err
	}

	messages, err := m.messRep.GetMessages(chatID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// easyjson:json
type WSRequest struct {
	Topic   string `json:"topic"`
	Action  string `json:"action"`
	Payload string `json:"payload"`
}

// easyjson:json
type WSMessagesRequest struct {
	ChatID   string                 `json:"chat_id"`
	Messages []*domain.ShortMessage `json:"messages"`
}

type wsService struct {
	wsPoolRep domain.WSPoolRepository
	messRep   domain.MessengerRepository
	logger    *zap.Logger
}

func NewWSService(wsPoolRep domain.WSPoolRepository, messRep domain.MessengerRepository, logger *zap.Logger) *wsService {
	return &wsService{
		wsPoolRep: wsPoolRep,
		messRep:   messRep,
		logger:    logger,
	}
}

func (w *wsService) EstablishConn(ctx context.Context, user *domain.User, conn net.Conn) {
	w.wsPoolRep.AddConnection(user.ID, conn)

	go func(user *domain.User) {
		defer conn.Close()
		defer w.wsPoolRep.RemoveConnection(user.ID, conn)

		for {
			msg, _, err := wsutil.ReadClientData(conn)
			if err != nil {
				w.logger.Error("fail read from ws", zap.Error(err))

				return
			}

			if err = w.parseRequest(ctx, conn, user, msg); err != nil {
				w.logger.Error("", zap.Error(err))
			}
		}
	}(user)
}

func (w *wsService) parseRequest(ctx context.Context, conn net.Conn, user *domain.User, msg []byte) error {
	var request WSRequest
	if err := request.UnmarshalJSON(msg); err != nil {
		return err
	}

	switch request.Topic {
	case "messenger":
		var r WSMessagesRequest

		if err := r.UnmarshalJSON([]byte(request.Payload)); err != nil {
			return err
		}

		if err := w.messRep.PersistMessages(user.ID, r.ChatID, r.Messages); err != nil {
			return err
		}

		if err := w.sendMessage(ctx, conn, user.ID, r.ChatID, r.Messages); err != nil {
			return err
		}
	}

	return nil
}

func (w *wsService) sendMessage(_ context.Context, c net.Conn, ownerID, chatID string, messages []*domain.ShortMessage) error {
	ids, err := w.messRep.GetParticipantsByChatID(ownerID, chatID)
	if err != nil {
		return err
	}

	ids = append(ids, ownerID)

	for _, id := range ids {
		conns := w.wsPoolRep.RetrieveConnByUserID(id)

		for _, conn := range conns {
			if conn == c {
				continue
			}

			bytePayload, err := WSMessagesRequest{
				ChatID:   chatID,
				Messages: messages,
			}.MarshalJSON()
			if err != nil {
				return err
			}

			byteRequest, err := WSRequest{
				Topic:   "messenger",
				Action:  "send-message",
				Payload: string(bytePayload),
			}.MarshalJSON()
			if err != nil {
				return err
			}

			if err = wsutil.WriteServerMessage(conn, ws.OpText, byteRequest); err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *wsService) Close() {
	w.wsPoolRep.FlushConnections()
}
