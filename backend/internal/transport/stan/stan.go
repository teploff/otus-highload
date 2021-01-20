package stan

import (
	"social-network/internal/domain"
	staninfrastructure "social-network/internal/infrastructure/stan"
	"time"

	"go.uber.org/zap"

	"github.com/nats-io/stan.go"
)

const maxInFlightMsg = 100

// Stan is client to the NATS server.
type Stan struct {
	client        *staninfrastructure.Client
	logger        *zap.Logger
	subscriptions []stan.Subscription
	doneCh        chan struct{}
}

// NewStan create instance of Stan.
func NewStan(client *staninfrastructure.Client, logger *zap.Logger) *Stan {
	return &Stan{
		client: client,
		logger: logger,
		doneCh: make(chan struct{}),
	}
}

// Serve starts listening stan.
func (s *Stan) Serve(cacheSvc domain.CacheService, wsSvc domain.WSService) error {
	subscriptionOptions := []stan.SubscriptionOption{
		stan.SetManualAckMode(),
		stan.AckWait(time.Second * 1),
		stan.MaxInflight(maxInFlightMsg),
	}

	friendsActionsSub, err := s.client.Subscribe("friends", makeFriendsSub(cacheSvc, s.logger),
		append(subscriptionOptions, stan.DurableName("friends-actions"))...)
	if err != nil {
		return err
	}

	s.subscriptions = append(s.subscriptions, friendsActionsSub)

	newsActionsSub, err := s.client.Subscribe("news", makeNewsSub(cacheSvc, s.logger),
		append(subscriptionOptions, stan.DurableName("news-actions"))...)
	if err != nil {
		return err
	}

	s.subscriptions = append(s.subscriptions, newsActionsSub)

	wsSub, err := s.client.Subscribe("news", makeNewsWSSub(wsSvc, s.logger),
		append(subscriptionOptions, stan.DurableName("news-ws"))...)
	if err != nil {
		return err
	}

	s.subscriptions = append(s.subscriptions, wsSub)

	for _, subscribe := range s.subscriptions {
		if err = subscribe.SetPendingLimits(-1, -1); err != nil {
			return err
		}
	}

	s.logger.Info("stan serving is starting")

	<-s.doneCh
	close(s.doneCh)

	s.logger.Info("stan serving is over")

	return nil
}

// Stop closes all subscribe connections.
func (s *Stan) Stop() {
	for _, sub := range s.subscriptions {
		if err := sub.Close(); err != nil {
			s.logger.Error("can't close stan connection: ", zap.Error(err))
		}
	}
	s.doneCh <- struct{}{}
	<-s.doneCh
}
