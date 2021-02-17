package main

import (
	"flag"
	"fmt"
	"gateway/internal/app"
	"gateway/internal/config"
	zapLogger "gateway/internal/infrastructure/logger"
	"gateway/internal/infrastructure/tracer"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	consulapi "github.com/hashicorp/consul/api"
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

	// fixme
	consulCfg := consulapi.DefaultConfig()
	consulCfg.Address = cfg.Consul.Addr

	consul, err := consulapi.NewClient(consulCfg)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}

	health, _, err := consul.Health().Service(cfg.Consul.ServiceName, "", false, nil)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}

	var servers []string
	for _, item := range health {
		addr := item.Service.Address + ":" + strconv.Itoa(item.Service.Port)
		servers = append(servers, addr)
	}
	fmt.Println(servers)
	// fixme

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
