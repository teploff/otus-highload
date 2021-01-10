package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"social-network/internal/config"
)

// Pool connection pool for Redis.
type Pool struct {
	dbPool map[int]*redis.Client
}

// NewPool gets Pool instance.
func NewPool(cfg []config.CacheConfig) (*Pool, error) {
	dbPool := make(map[int]*redis.Client, len(cfg))

	for _, cacheConfig := range cfg {
		conn := redis.NewClient(&redis.Options{
			Addr:     cacheConfig.Addr,
			Password: cacheConfig.Password,
			DB:       cacheConfig.DB,
		})

		if _, err := conn.Ping(context.TODO()).Result(); err != nil {
			return nil, fmt.Errorf("redis ping fail, %w", err)
		}

		dbPool[cacheConfig.DB] = conn
	}

	return &Pool{dbPool: dbPool}, nil
}

func (c *Pool) GetConnByIndexDB(index int) (*redis.Client, error) {
	conn, exist := c.dbPool[index]
	if !exist {
		return nil, fmt.Errorf("index db absent")
	}

	return conn, nil
}

// Close connection.
func (c *Pool) Close() {
	for _, conn := range c.dbPool {
		conn.Close()
	}
}
