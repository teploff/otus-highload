package app

import (
	"auth/internal/config"
	"auth/internal/implementation"
	httptransport "auth/internal/transport/http"
	"context"
	"database/sql"
	"log"
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
	cfg     *config.Config
	httpSrv *http.Server
	logger  *zap.Logger
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
func (a *App) Run(mysqlConn *sql.DB) {
	authSvc := implementation.NewAuthService(implementation.NewUserRepository(mysqlConn), a.cfg.JWT)
	profileSvc := implementation.NewProfileService(implementation.NewUserRepository(mysqlConn))

	a.httpSrv = httptransport.NewHTTPServer(a.cfg.Addr,
		httptransport.MakeEndpoints(authSvc, profileSvc))

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
}
