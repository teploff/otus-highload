package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"social-network/internal/config"
	"social-network/internal/domain"
	"social-network/internal/implementation"
	"social-network/internal/infrastructure/cache"
	"social-network/internal/infrastructure/stan"
	httptransport "social-network/internal/transport/http"
	stantransport "social-network/internal/transport/stan"
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
	stanSrv *stantransport.Stan
	wsSvc   domain.WSService
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
func (a *App) Run(mysqlConn *sql.DB, redisPool *cache.Pool, stanClient *stan.Client) {
	a.stanSrv = stantransport.NewStan(stanClient, a.logger)
	wsPoolRep := implementation.NewWSPoolRepository()
	a.wsSvc = implementation.NewWSService(
		implementation.NewUserRepository(mysqlConn),
		implementation.NewCacheRepository(redisPool),
		wsPoolRep,
		stanClient,
		a.logger)

	authSvc := implementation.NewAuthService(implementation.NewUserRepository(mysqlConn), a.cfg.JWT)
	profileSvc := implementation.NewProfileService(implementation.NewUserRepository(mysqlConn))
	socialSvc := implementation.NewSocialService(
		implementation.NewUserRepository(mysqlConn),
		implementation.NewSocialRepository(mysqlConn),
		implementation.NewCacheRepository(redisPool),
		stanClient,
	)
	messengerSvc := implementation.NewMessengerService(
		implementation.NewUserRepository(mysqlConn),
		implementation.NewMessengerRepository(mysqlConn))
	cacheSvc := implementation.NewCacheService(implementation.NewCacheRepository(redisPool))

	a.httpSrv = httptransport.NewHTTPServer(a.cfg.Addr,
		httptransport.MakeEndpoints(authSvc, profileSvc, socialSvc, messengerSvc, a.wsSvc))

	go func() {
		if err := a.httpSrv.ListenAndServe(); err != nil {
			a.logger.Fatal("http serve error, ", zap.Error(err))
		}
	}()

	go func() {
		if err := a.stanSrv.Serve(cacheSvc, a.wsSvc); err != nil {
			a.logger.Fatal("stan server error", zap.Error(err))
		}
	}()
}

func (a *App) Stop() {
	a.wsSvc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), httpTimeoutClose)
	defer cancel()

	if err := a.httpSrv.Shutdown(ctx); err != nil {
		log.Fatal("http closing error, ", zap.Error(err))
	}

	a.stanSrv.Stop()
}
