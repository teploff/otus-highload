package app

import (
	"context"
	"database/sql"
	"log"
	"messenger/internal/config"
	"messenger/internal/implementation"
	"messenger/internal/infrastructure/clickhouse"
	httptransport "messenger/internal/transport/http"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const httpTimeoutClose = 5 * time.Second

type Option func(*App)

// WithLogger adding logger option.
func WithLogger(l *zap.Logger) Option {
	return func(a *App) {
		a.logger = l
	}
}

// App is main application instance.
type App struct {
	cfg       *config.Config
	chStorage *clickhouse.Storage
	httpSrv   *http.Server
	logger    *zap.Logger
}

// NewApp returns instance of app.
func NewApp(cfg *config.Config, opts ...Option) *App {
	app := &App{
		cfg:    cfg,
		logger: zap.NewNop(),
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

// Run lunch application.
func (a *App) Run(chConn *sql.DB) {
	a.chStorage = clickhouse.NewStorage(chConn, a.cfg.Clickhouse.PushTimeout, a.logger)
	go a.chStorage.StartBatching()

	authSvc := implementation.NewAuthService(a.cfg.Auth.Addr)
	//messengerSvc := implementation.NewMessengerService(implementation.NewMessengerRepository(a.chStorage))
	wsSvc := implementation.NewWSService(implementation.NewWSPoolRepository(), a.logger)

	a.httpSrv = httptransport.NewHTTPServer(a.cfg.Addr, httptransport.MakeEndpoints(authSvc, wsSvc))

	go func() {
		if err := a.httpSrv.ListenAndServe(); err != nil {
			a.logger.Fatal("http serve error, ", zap.Error(err))
		}
	}()
}

func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeoutClose)
	defer cancel()

	if err := a.httpSrv.Shutdown(ctx); err != nil {
		log.Fatal("http closing error, ", zap.Error(err))
	}

	a.chStorage.Stop()
}
