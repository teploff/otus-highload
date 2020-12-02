package main

import (
	"backend/internal/app"
	"backend/internal/config"
	"database/sql"
	"flag"
	"os"
	"os/signal"
	"syscall"

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


	application := app.NewApp(cfg,
		app.WithLogger(logger),
	)
	go application.Run(mysqlConn)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	application.Stop()
}
