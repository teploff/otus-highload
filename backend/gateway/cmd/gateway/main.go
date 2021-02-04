package main

import (
	"flag"
	"gateway/internal/app"
	"gateway/internal/config"
	zapLogger "gateway/internal/infrastructure/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configFile := flag.String("config", "./configs/config.yaml", "configuration file path")
	flag.Parse()

	cfg, err := config.Load(*configFile)
	if err != nil {
		panic(err)
	}

	logger := zapLogger.New(&cfg.Logger)

	messengerConn, err := grpc.Dial(cfg.Messenger.GRPCAddr, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("gRPC auth connection", zap.Error(err))
	}
	defer messengerConn.Close()

	application := app.New(cfg,
		app.WithLogger(logger),
	)
	go application.Run(messengerConn)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	application.Stop()
}
