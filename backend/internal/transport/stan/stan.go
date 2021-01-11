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
func (s *Stan) Serve(service domain.CacheService) error {
	subscriptionOptions := []stan.SubscriptionOption{
		stan.SetManualAckMode(),
		stan.AckWait(time.Second * 1),
		stan.MaxInflight(maxInFlightMsg),
	}

	compSub, err := s.client.Subscribe("friends", makeFriendsSub(service, s.logger),
		append(subscriptionOptions, stan.DurableName("friends-actions"))...)

	if err != nil {
		return err
	}

	s.subscriptions = append(s.subscriptions, compSub)

	forecSub, err := s.client.Subscribe("news", makeNewsSub(service, s.logger),
		append(subscriptionOptions, stan.DurableName("news-actions"))...)

	if err != nil {
		return err
	}

	s.subscriptions = append(s.subscriptions, forecSub)

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
