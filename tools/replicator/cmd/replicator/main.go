package main

import (
	"flag"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"replicator/internal/config"
	"replicator/internal/implementation"
	"replicator/internal/infrastructure/tarantool"
	"syscall"
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

	conn, err := tarantool.NewConn(cfg.Tarantool)
	if err != nil {
		logger.Fatal("Failed connect to tarantool", zap.Error(err))
	}
	defer conn.Close()

	syncer, err := implementation.NewMySQLSyncer(cfg, conn, logger)
	if err != nil {
		log.Fatal("failed syncer launch", zap.Error(err))
	}

	go syncer.Run()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	syncer.Stop()
}
