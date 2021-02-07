package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"social/internal/config"
	"social/internal/domain"
	staninfrastructure "social/internal/infrastructure/stan"
	"time"

	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

const (
	maxInFlightMsg = 100
	friendsDBIndex = 1
	NewsDBIndex    = 2
)

type Heater struct {
	authMySQLConn   *sql.DB
	socialMySQLConn *sql.DB
	redisPool       *Pool
	stanClient      *staninfrastructure.Client
	subscription    stan.Subscription
	logger          *zap.Logger
	doneCh          chan struct{}
}

func NewHeater(cfg config.HeaterConfig, redisPool *Pool, stanClient *staninfrastructure.Client, logger *zap.Logger) (*Heater, error) {
	authConn, err := sql.Open("mysql", cfg.AuthDSN)
	if err != nil {
		return nil, fmt.Errorf("fail to auth connection db")
	}

	socialConn, err := sql.Open("mysql", cfg.SocialDSN)
	if err != nil {
		return nil, fmt.Errorf("fail to social connection db")
	}

	// See "Important settings" section.
	authConn.SetConnMaxLifetime(10)
	socialConn.SetConnMaxLifetime(10)
	authConn.SetMaxOpenConns(10)
	socialConn.SetMaxOpenConns(10)
	authConn.SetMaxIdleConns(10)
	socialConn.SetMaxIdleConns(10)

	return &Heater{
		authMySQLConn:   authConn,
		socialMySQLConn: socialConn,
		redisPool:       redisPool,
		stanClient:      stanClient,
		logger:          logger,
		doneCh:          make(chan struct{}),
	}, err
}

func (h *Heater) Heat() error {
	ctx := context.Background()

	if err := h.heatFriends(ctx); err != nil {
		return err
	}

	if err := h.heatNews(ctx); err != nil {

		return err
	}

	return nil
}

func (h *Heater) Listening() (err error) {
	subscriptionOptions := []stan.SubscriptionOption{
		stan.SetManualAckMode(),
		stan.AckWait(time.Second * 1),
		stan.MaxInflight(maxInFlightMsg),
	}

	h.subscription, err = h.stanClient.Subscribe("cache-reload", h.makeCacheReloadSub(),
		append(subscriptionOptions, stan.DurableName("cache-heater"))...)
	if err != nil {
		return err
	}

	<-h.doneCh
	close(h.doneCh)

	h.logger.Info("cache heater stan subscribing is over")

	return nil
}

func (h *Heater) heatFriends(ctx context.Context) error {
	tx, err := h.socialMySQLConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	friends := make(map[string][]string, 1)

	rows, err := tx.Query(`
		SELECT
			master_user_id, slave_user_id
		FROM
			friendship
		WHERE
		    friendship.status = ?`, "accepted")
	if err != nil {
		tx.Rollback()

		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id1, id2 string

		if err = rows.Scan(&id1, &id2); err != nil {
			tx.Rollback()

			return err
		}

		if _, ok := friends[id1]; !ok {
			friends[id1] = make([]string, 0, 1)
		}
		friends[id1] = append(friends[id1], id2)

		if _, ok := friends[id2]; !ok {
			friends[id2] = make([]string, 0, 1)
		}
		friends[id2] = append(friends[id2], id1)
	}

	conn, err := h.redisPool.GetConnByIndexDB(friendsDBIndex)
	if err != nil {
		tx.Rollback()

		return err
	}

	for key, value := range friends {
		data, err := json.Marshal(value)
		if err != nil {
			tx.Rollback()

			return fmt.Errorf("cannot marshal friends id, %w", err)
		}

		if err = conn.Set(ctx, key, data, 0).Err(); err != nil {
			tx.Rollback()

			return err
		}
	}

	return tx.Commit()
}

func (h *Heater) heatNews(ctx context.Context) error {
	usersNews := make(map[string][]*domain.News, 1)

	socialTx, err := h.socialMySQLConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	rows, err := socialTx.Query(`
		SELECT
			id, owner_id, content, create_time
		FROM
			news
		ORDER BY create_time DESC
		`)
	if err != nil {
		socialTx.Rollback()

		return err
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		news := new(domain.News)

		if err = rows.Scan(&news.ID, &userID, &news.Content, &news.CreateTime); err != nil {
			return err
		}

		if _, ok := usersNews[userID]; !ok {
			usersNews[userID] = make([]*domain.News, 0, 1)
		}
		usersNews[userID] = append(usersNews[userID], news)
	}

	if err = socialTx.Commit(); err != nil {
		return err
	}

	if len(usersNews) == 0 {
		return nil
	}

	userIDs := make([]string, 0, len(usersNews))
	for id := range usersNews {
		userIDs = append(userIDs, id)
	}

	authTx, err := h.authMySQLConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	sqlStr := "SELECT id, name, surname, sex FROM user WHERE id IN ("
	vals := make([]interface{}, 0, len(userIDs))

	for _, id := range userIDs {
		sqlStr += "?,"
		vals = append(vals, id)
	}

	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	// add ) with limit and offset
	sqlStr += ")"

	//prepare the statement
	stmt, err := authTx.Prepare(sqlStr)
	if err != nil {
		authTx.Rollback()

		return err
	}

	rows, err = stmt.Query(vals...)
	if err != nil {
		authTx.Rollback()

		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, name, surname, sex string

		if err = rows.Scan(&id, &name, &surname, &sex); err != nil {
			authTx.Rollback()

			return err
		}

		for _, news := range usersNews[id] {
			news.Owner.Name = name
			news.Owner.Surname = surname
			news.Owner.Sex = sex
		}
	}

	if err = authTx.Commit(); err != nil {
		return err
	}

	conn, err := h.redisPool.GetConnByIndexDB(NewsDBIndex)
	if err != nil {
		return err
	}

	for key, value := range usersNews {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("cannot marshal friends id, %w", err)
		}

		if err = conn.Set(ctx, key, data, 0).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (h *Heater) Stop() {
	h.authMySQLConn.Close()
	h.socialMySQLConn.Close()

	if err := h.subscription.Close(); err != nil {
		h.logger.Error("cache heater can't close stan connection: ", zap.Error(err))
	}
}

func (h *Heater) makeCacheReloadSub() func(msg *stan.Msg) {
	return func(msg *stan.Msg) {
		if err := msg.Ack(); err != nil {
			h.logger.Error("cache heater fail to acknowledge a message: ", zap.Error(err))
		}

		ctx := context.Background()

		if err := h.heatFriends(ctx); err != nil {
			h.logger.Error("cache heater fail to heat friends", zap.Error(err))
		}

		if err := h.heatNews(ctx); err != nil {
			h.logger.Error("cache heater fail to heat news", zap.Error(err))
		}
	}
}
