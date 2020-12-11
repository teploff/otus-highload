package main

import (
	"backend/internal/app"
	"backend/internal/config"
	"context"
	"database/sql"
	"flag"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

func main() {
	configFile := flag.String("config", "./configs/config.yaml", "configuration file path")
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	cfg, err := config.Load(*configFile)
	if err != nil {
		logger.Fatal("error reading config file", zap.Error(err))
	}

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

	application := app.NewApp(cfg,
		app.WithLogger(logger),
	)
	go application.Run(mysqlConn, chConn, redisConn)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	application.Stop()
}
