package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "messenger/api/swagger"
	"messenger/internal/app"
	"messenger/internal/config"
	zaplogger "messenger/internal/infrastructure/logger"
	"messenger/internal/infrastructure/tracer"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

// @title Messenger API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:10003
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
		panic(fmt.Sprintf("error reading config file %s", err))
	}

	logger := zaplogger.NewLogger(&cfg.Logger)

	closer, err := tracer.InitGlobalTracer(cfg.Jaeger.ServiceName, cfg.Jaeger.AgentAddr)
	if err != nil {
		logger.Fatal("fail to connect jaeger", zap.Error(err))
	}
	defer closer.Close()

	chConn, err := sql.Open("clickhouse", cfg.Clickhouse.DSN)
	if err != nil {
		logger.Fatal("clickhouse connection fail", zap.Error(err))
	}

	if err = chConn.Ping(); err != nil {
		logger.Fatal("clickhouse ping fail, ", zap.Error(err))
	}
	defer chConn.Close()

	application := app.NewApp(cfg,
		app.WithLogger(logger),
	)
	go application.Run(chConn)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	application.Stop()
}
