package stan

import (
	"context"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
	"social/internal/domain"
)

func makeFriendsSub(service domain.CacheService, logger *zap.Logger) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		if err := msg.Ack(); err != nil {
			logger.Error("fail to acknowledge a message: ", zap.Error(err))
		}

		var request FriendsActionRequest
		if err := request.UnmarshalJSON(msg.Data); err != nil {
			logger.Error("fail to unmarshal msg: ", zap.Error(err))
		}

		switch request.Action {
		case "persist":
			if err := service.AddFriends(context.TODO(), request.UserID, request.FriendsID); err != nil {
				logger.Error("fail persist friends: ", zap.Error(err))
			}
		case "delete":
			if err := service.DeleteFriends(context.TODO(), request.UserID, request.FriendsID); err != nil {
				logger.Error("fail to persist friends: ", zap.Error(err))
			}
		default:
			logger.Error("unknown action to friends")
		}
	}
}

func makeNewsSub(service domain.CacheService, logger *zap.Logger) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		if err := msg.Ack(); err != nil {
			logger.Error("fail to acknowledge a message: ", zap.Error(err))
		}

		var request NewsPersistRequest
		if err := request.UnmarshalJSON(msg.Data); err != nil {
			logger.Error("fail to unmarshal msg: ", zap.Error(err))
		}

		if err := service.AddNews(context.TODO(), request.OwnerID, request.News); err != nil {
			logger.Error("fail to persist news: ", zap.Error(err))
		}
	}
}

func makeNewsWSSub(service domain.WSService, logger *zap.Logger) func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		if err := msg.Ack(); err != nil {
			logger.Error("fail to acknowledge a message: ", zap.Error(err))
		}

		var request NewsPersistRequest
		if err := request.UnmarshalJSON(msg.Data); err != nil {
			logger.Error("fail to unmarshal msg: ", zap.Error(err))
		}

		if err := service.SendNews(context.TODO(), request.OwnerID, request.News); err != nil {
			logger.Error("fail to send news to friends via socket: ", zap.Error(err))
		}
	}
}
