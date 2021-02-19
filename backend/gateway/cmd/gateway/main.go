package main

import (
	"flag"
	"gateway/internal/app"
	"gateway/internal/config"
	zapLogger "gateway/internal/infrastructure/logger"
	"gateway/internal/infrastructure/tracer"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	configFile := flag.String("config", "./configs/config.yaml", "configuration file path")
	flag.Parse()

	cfg, err := config.Load(*configFile)
	if err != nil {
		panic(err)
	}

	logger := zapLogger.New(&cfg.Logger)

	closer, err := tracer.InitGlobalTracer(cfg.Jaeger.ServiceName, cfg.Jaeger.AgentAddr)
	if err != nil {
		logger.Fatal("fail to connect jaeger", zap.Error(err))
	}
	defer closer.Close()

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
