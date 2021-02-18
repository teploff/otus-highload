package app

import (
	"context"
	"gateway/internal/config"
	"gateway/internal/domain"
	"gateway/internal/implementation"
	"gateway/internal/infrastructure/consul"
	kitgrpc "gateway/internal/transport/grpc"
	httptransport "gateway/internal/transport/http"
	"net/http"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const httpTimeoutClose = 5 * time.Second

type Option func(app *App)

func WithLogger(logger *zap.Logger) Option {
	return func(app *App) {
		app.logger = logger
	}
}

type App struct {
	cfg          *config.Config
	httpSrv      *http.Server
	srvList      *domain.ServerAvailableList
	consulClient *consul.Client
	logger       *zap.Logger
}

func New(cfg *config.Config, opts ...Option) *App {
	app := &App{
		cfg:     cfg,
		srvList: domain.NewServerAvailableList(),
		logger:  zap.NewNop(),
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (a *App) Run(messengerConn *grpc.ClientConn) {
	consulClient, err := consul.NewClient(a.cfg.Consul, a.srvList, a.logger)
	if err != nil {
		a.logger.Fatal("fail consul init")
	}
	a.consulClient = consulClient

	go a.consulClient.HealthCheck()

	a.httpSrv = httptransport.NewHTTPServer(a.cfg.Addr,
		httptransport.MakeEndpoints(a.cfg,
			implementation.NewGRPCMessengerProxyService(kitgrpc.MakeMessengerProxyEndpoints(messengerConn)),
			a.srvList,
		),
	)

	go func() {
		if err = a.httpSrv.ListenAndServe(); err != nil {
			a.logger.Fatal("http serve error", zap.Error(err))
		}
	}()
}

func (a *App) Stop() {
	a.consulClient.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeoutClose)
	defer cancel()

	if err := a.httpSrv.Shutdown(ctx); err != nil {
		a.logger.Fatal("http closing error", zap.Error(err))
	}
}
