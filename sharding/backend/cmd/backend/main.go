package main

import (
	"backend/internal/app"
	"backend/internal/config"
	zaplogger "backend/internal/infrastructure/logger"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/nats-io/stan.go"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

const (
	pingInterval = 60
	pingAttempts = 2 * 60
)

func main() {
	configFile := flag.String("config", "./configs/config.yaml", "configuration file path")
	flag.Parse()

	cfg, err := config.Load(*configFile)
	if err != nil {
		panic(fmt.Sprintf("error reading config file %s", err))
	}

	logger := zaplogger.NewZapLogger(cfg.Logger)

	mysqlConn, err := sql.Open("mysql", cfg.Storage.DSN)
	if err != nil {
		logger.Fatal("mysql connection fail", zap.Error(err))
	}
	defer mysqlConn.Close()

	// See "Important settings" section.
	mysqlConn.SetConnMaxLifetime(cfg.Storage.ConnMaxLifetime)
	mysqlConn.SetMaxOpenConns(cfg.Storage.MaxOpenConns)
	mysqlConn.SetMaxIdleConns(cfg.Storage.MaxIdleConns)

	if err = mysqlConn.Ping(); err != nil {
		logger.Fatal("mysql ping fail, ", zap.Error(err))
	}

	chConn, err := sql.Open("clickhouse", cfg.Clickhouse.DSN)
	if err != nil {
		logger.Fatal("clickhouse connection fail", zap.Error(err))
	}

	if err = chConn.Ping(); err != nil {
		logger.Fatal("clickhouse ping fail, ", zap.Error(err))
	}
	defer chConn.Close()

	redisConn := redis.NewClient(&redis.Options{
		Addr:     cfg.Cache.Addr,
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.DB,
	})

	if _, err = redisConn.Ping(context.TODO()).Result(); err != nil {
		logger.Fatal("redis ping fail, ", zap.Error(err))
	}
	defer redisConn.Close()

	stanConn, err := stan.Connect(cfg.Stan.ClusterID, "backend-sharding",
		stan.Pings(pingInterval, pingAttempts),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			logger.Fatal("Connection lost, reason: ", zap.Error(reason))
		}), stan.NatsURL(cfg.Stan.Addr))
	if err != nil {
		logger.Fatal("stan transport", zap.Error(err))
	}
	defer stanConn.Close()

	application := app.NewApp(cfg,
		app.WithLogger(logger),
	)
	go application.Run(mysqlConn, chConn, redisConn, stanConn)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	application.Stop()
}
