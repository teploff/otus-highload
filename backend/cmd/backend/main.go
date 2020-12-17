package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"os/signal"
	"social-network/internal/app"
	"social-network/internal/config"
	zaplogger "social-network/internal/infrastructure/logger"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

func main() {
	configFile := flag.String("config", "./configs/config.yaml", "configuration file path")
	flag.Parse()

	cfg, err := config.Load(*configFile)
	if err != nil {
		panic(fmt.Sprintf("error reading config file %s", err.Error()))
	}

	logger := zaplogger.NewLogger(&cfg.Logger)

	mysqlConn, err := sql.Open("mysql", cfg.Storage.DSN)
	if err != nil {
		logger.Fatal("mysql connection fail", zap.Error(err))
	}
	defer mysqlConn.Close()

	// See "Important settings" section.
	mysqlConn.SetConnMaxLifetime(cfg.Storage.ConnMaxLifetime)
	mysqlConn.SetMaxOpenConns(cfg.Storage.MaxOpenConns)
	mysqlConn.SetMaxIdleConns(cfg.Storage.MaxIdleConns)

	logger.Info("try establish connection with MySQL...")
	if err = establishConnection(mysqlConn, cfg.Storage.AttemptCount); err != nil {
		logger.Fatal("mysql ping fail, ", zap.Error(err))
	}
	logger.Info("connection with MySQL is established")

	redisConn := redis.NewClient(&redis.Options{
		Addr:     cfg.Cache.Addr,
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.DB,
	})

	logger.Info("try establish connection with Redis...")
	if _, err = redisConn.Ping(context.TODO()).Result(); err != nil {
		logger.Fatal("redis ping fail, ", zap.Error(err))
	}
	logger.Info("connection with Redis is established")
	defer redisConn.Close()

	application := app.NewApp(cfg,
		app.WithLogger(logger),
	)
	go application.Run(mysqlConn)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	application.Stop()
}

func establishConnection(conn *sql.DB, attemptCount int) error {
	var err error

	for i := 0; i < attemptCount; i++ {
		if err = conn.Ping(); err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return err
}
