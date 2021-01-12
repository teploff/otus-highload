package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"social-network/internal/app"
	"social-network/internal/config"
	"social-network/internal/infrastructure/cache"
	zaplogger "social-network/internal/infrastructure/logger"
	"social-network/internal/infrastructure/stan"
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

	logger.Info("try establish connection with Redis...")
	redisPool, err := cache.NewPool(cfg.Cache)
	if err != nil {
		logger.Fatal("redis connection fail, ", zap.Error(err))
	}
	logger.Info("connection with Redis is established")
	defer redisPool.Close()

	stanClient, err := stan.NewClient(cfg.Stan, logger)
	if err != nil {
		logger.Fatal("stan connection fail", zap.Error(err))
	}
	defer stanClient.Close()

	logger.Info("cache heater is starting...")
	cacheHeater := cache.NewHeater(redisPool, mysqlConn, stanClient, logger)
	if err = cacheHeater.Heat(); err != nil {
		logger.Fatal("fail to start cache heater", zap.Error(err))
	}
	logger.Info("cache heater work is over")

	logger.Info("cache heater is starting listening actions")
	go func() {
		if err = cacheHeater.Listening(); err != nil {
			logger.Fatal("fail to start listeting cache heater", zap.Error(err))
		}
	}()

	application := app.NewApp(cfg,
		app.WithLogger(logger),
	)
	go application.Run(mysqlConn, redisPool, stanClient)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	application.Stop()
	cacheHeater.Stop()
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
