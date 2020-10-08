package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"os/signal"
	"social-network/internal/config"
	"social-network/internal/implementation"
	"social-network/internal/transport/http"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

const httpTimeoutClose = 5 * time.Second

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

	authSvc := implementation.NewAuthService(implementation.NewUserRepository(mysqlConn), cfg.JWTSecret)
	socialSvc := implementation.NewSocialService(implementation.NewUserRepository(mysqlConn))
	srv := http.NewHTTPServer(cfg.Addr, http.MakeEndpoints(authSvc, socialSvc))

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			logger.Fatal("http serve error, ", zap.Error(err))
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), httpTimeoutClose)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal("http closing error, ", zap.Error(err))
	}
}
