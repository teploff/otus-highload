package main

import (
	_ "auth/api/swagger"
	"auth/internal/app"
	"auth/internal/config"
	zaplogger "auth/internal/infrastructure/logger"
	"auth/internal/infrastructure/tracer"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

// @title Auth API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:10001
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
