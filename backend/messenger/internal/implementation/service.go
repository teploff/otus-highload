package implementation

import (
	"context"
	"encoding/binary"
	"github.com/gobwas/ws/wsutil"
	"github.com/imroc/req"
	"go.uber.org/zap"
	"messenger/internal/config"
	"messenger/internal/domain"
	"net"
)

type authService struct {
	authAddr string
}

func NewAuthService(authAddr string) *authService {
	return &authService{
		authAddr: authAddr,
	}
}

type getUserIDByAccessTokenResponse struct {
	UserID string `json:"user_id"`
}

func (a *authService) Authenticate(_ context.Context, token string) (string, error) {
	header := req.Header{
		"Accept":        "application/json",
		"Authorization": token,
	}

	r, err := req.Get("http://"+a.authAddr+"/auth/user/get-id-by-token", header)
	if err != nil {
		return "", err
	}

	var response getUserIDByAccessTokenResponse
	if err = r.ToJSON(&response); err != nil {
		return "", err
	}

	return response.UserID, nil
}

type messengerService struct {
	messRep     domain.MessengerRepository
	shardingCfg config.ShardingConfig
}

func NewMessengerService(messengerRep domain.MessengerRepository, cfg config.ShardingConfig) *messengerService {
	return &messengerService{
		messRep:     messengerRep,
		shardingCfg: cfg,
	}
}

func (m *messengerService) CreateChat(ctx context.Context, masterID, slaveID string) (string, error) {
	//tx, err := m.userRep.GetTx(ctx)
	//if err != nil {
	//	return "", err
	//}
	//
	//_, err = m.userRep.GetByID(tx, slaveID)
	//switch err {
	//case nil:
	//case sql.ErrNoRows:
	//	return "", fmt.Errorf("chat companion with id=[%s] doesn't exist", slaveID)
	//default:
	//	return "", err
	//}

	chatID, err := m.messRep.CreateChat(masterID, slaveID)
	if err != nil {
		return "", err
	}

	return chatID, nil
}

func (m *messengerService) GetChat(_ context.Context, masterID, slaveID string) (*domain.Chat, error) {
	chat, err := m.messRep.GetChatWithCompanion(masterID, slaveID)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (m *messengerService) GetChats(_ context.Context, userID string, limit, offset int) ([]*domain.Chat, int, error) {
	total, err := m.messRep.GetCountChats(userID)
	if err != nil {
		return nil, 0, err
	}

	chats, err := m.messRep.GetChats(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return chats, total, nil
}

func (m *messengerService) SendMessages(_ context.Context, userID, chatID string, messages []*domain.ShortMessage) error {
	shardID := int(binary.BigEndian.Uint64([]byte(userID)) % uint64(m.shardingCfg.CountNodes))

	return m.messRep.SendMessages(shardID, userID, chatID, messages)
}

func (m *messengerService) GetMessages(_ context.Context, userID, chatID string, limit, offset int) ([]*domain.Message, int, error) {
	_, err := m.messRep.GetChatAsParticipant(userID)
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
	Payload string `json:"payload"`
}

// easyjson:json
type WSMessagesRequest struct {
	Content string `json:"content"`
}

type wsService struct {
	wsPoolRep domain.WSPoolRepository
	logger    *zap.Logger
}

func NewWSService(wsPoolRep domain.WSPoolRepository, logger *zap.Logger) *wsService {
	return &wsService{
		wsPoolRep: wsPoolRep,
		logger:    logger,
	}
}

func (w *wsService) EstablishConn(ctx context.Context, userID string, conn net.Conn) {
	w.wsPoolRep.AddConnection(userID, conn)

	go func(userID string) {
		defer conn.Close()
		defer w.wsPoolRep.RemoveConnection(userID, conn)

		for {
			msg, _, err := wsutil.ReadClientData(conn)
			if err != nil {
				w.logger.Error("fail read from ws", zap.Error(err))

				return
			}

			w.logger.Info(string(msg))

			//if err = w.parseRequest(ctx, user, msg); err != nil {
			//	w.logger.Error("", zap.Error(err))
			//}
		}
	}(userID)
}

//func (w *wsService) parseRequest(ctx context.Context, user *domain.User, msg []byte) error {
//	var request WSRequest
//	if err := request.UnmarshalJSON(msg); err != nil {
//		return err
//	}
//
//	switch request.Topic {
//	case "news":
//		var r WSNewsRequest
//
//		if err := r.UnmarshalJSON([]byte(request.Payload)); err != nil {
//			return err
//		}
//
//		n := &domain.News{
//			ID: uuid.NewV4().String(),
//			Owner: struct {
//				Name    string `json:"name"`
//				Surname string `json:"surname"`
//				Sex     string `json:"sex"`
//			}{
//				user.Name,
//				user.Surname,
//				user.Sex,
//			},
//			Content:    r.Content,
//			CreateTime: time.Now().UTC(),
//		}
//		news := []*domain.News{n}
//
//		tx, err := w.socialRep.GetTx(ctx)
//		if err != nil {
//			return err
//		}
//
//		if err = w.socialRep.PublishNews(tx, user.ID, news); err != nil {
//			return err
//		}
//
//		if err = tx.Commit(); err != nil {
//			return err
//		}
//
//		if err = w.stanClient.Publish("news", &stantransport.NewsPersistRequest{OwnerID: user.ID, News: []*domain.News{n}}); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

//func (w *wsService) SendNews(ctx context.Context, ownerID string, news []*domain.News) error {
//	ids, err := w.socialCacheRep.RetrieveFriendsID(ctx, ownerID)
//	if err != nil {
//		return err
//	}
//	ids = append(ids, ownerID)
//
//	for _, id := range ids {
//		conns := w.wsPoolRep.RetrieveConnByUserID(id)
//		for _, conn := range conns {
//			for _, n := range news {
//				data, err := n.MarshalJSON()
//				if err != nil {
//					return err
//				}
//
//				if err = wsutil.WriteServerMessage(conn, ws.OpText, data); err != nil {
//					return err
//				}
//			}
//		}
//	}
//
//	return nil
//}

func (w *wsService) Close() {
	w.wsPoolRep.FlushConnections()
}
