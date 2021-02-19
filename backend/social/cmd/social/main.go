package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"social/internal/app"
	"social/internal/config"
	"social/internal/infrastructure/cache"
	"social/internal/infrastructure/consul"
	zaplogger "social/internal/infrastructure/logger"
	"social/internal/infrastructure/stan"
	"social/internal/infrastructure/tracer"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

// @title Social API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:10002
// @BasePath /
// @query.collection.format multi

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	configFile := flag.String("config", "./configs/config.yaml", "configuration file path")
	flag.Parse()

	cfg, err := config.Load(*configFile)
	if err != nil {
		panic(fmt.Sprintf("error reading config file %s", err.Error()))
	}

	logger := zaplogger.NewLogger(&cfg.Logger)

	closer, err := tracer.InitGlobalTracer(cfg.Jaeger.ServiceName, cfg.Jaeger.AgentAddr)
	if err != nil {
		logger.Fatal("fail to connect jaeger", zap.Error(err))
	}
	defer closer.Close()

	consulClient, err := consul.NewClient(cfg.Consul)
	if err != nil {
		logger.Fatal("fail to connect Consul", zap.Error(err))
	}

	if err = consulClient.Register(); err != nil {
		logger.Fatal("", zap.Error(err))
	}

	defer func() {
		if err = consulClient.Deregister(); err != nil {
			logger.Fatal("", zap.Error(err))
		}

		logger.Info("service auth deregister in consul")
	}()

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

	logger.Info("cache heater is initializing...")
	cacheHeater, err := cache.NewHeater(cfg.Heater, redisPool, stanClient, logger)
	if err != nil {
		logger.Fatal("fail to initialize cache heater", zap.Error(err))
	}

	logger.Info("cache heater is starting...")
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
