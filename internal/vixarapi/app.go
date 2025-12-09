package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	oas "github.com/keenywheels/backend/internal/api/v1"
	"github.com/keenywheels/backend/internal/pkg/client/vk"
	securityApi "github.com/keenywheels/backend/internal/vixarapi/delivery/http/security"
	api "github.com/keenywheels/backend/internal/vixarapi/delivery/http/v1"
	pgRepository "github.com/keenywheels/backend/internal/vixarapi/repository/postgres"
	redisRepository "github.com/keenywheels/backend/internal/vixarapi/repository/redis"
	"github.com/keenywheels/backend/internal/vixarapi/service"
	"github.com/keenywheels/backend/pkg/cors"
	"github.com/keenywheels/backend/pkg/httpserver"
	"github.com/keenywheels/backend/pkg/httputils"
	"github.com/keenywheels/backend/pkg/logger"
	"github.com/keenywheels/backend/pkg/logger/zap"
	mw "github.com/keenywheels/backend/pkg/middleware"
	"github.com/keenywheels/backend/pkg/postgres"
	"github.com/keenywheels/backend/pkg/redis"
	"golang.org/x/sync/errgroup"
)

// middleware allias for middleware funcs
type middleware func(http.Handler) http.Handler

// App represent app environment
type App struct {
	opts *Options

	cfg    *Config
	logger logger.Logger
}

// New creates new application instance with options
func New() *App {
	opts := NewDefaultOpts()
	opts.LoadEnv()
	opts.LoadFlags()

	return &App{
		opts: opts,
	}
}

// Run starts app
func (app *App) Run() error {
	// read config
	cfg, err := LoadConfig(app.opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load app config: %w", err)
	}

	app.cfg = cfg
	app.initLogger()
	defer func() {
		if err := app.logger.Close(); err != nil {
			log.Printf("failed to close logger: %v", err)
		}
	}()

	// create postgres connection
	db, err := app.getPostgresConn()
	if err != nil {
		return fmt.Errorf("failed to create postgres connection: %w", err)
	}
	defer db.Close()

	// create postgres repository
	pgRepo, err := pgRepository.New(db)
	if err != nil {
		return fmt.Errorf("failed to create interest repository: %w", err)
	}

	// create redis repository
	redisClient, err := redis.New(&app.cfg.RedisCfg)
	if err != nil {
		return fmt.Errorf("failed to create redis client: %w", err)
	}

	redisRepo, err := redisRepository.New(redisClient)
	if err != nil {
		return fmt.Errorf("failed to create redis repository: %w", err)
	}

	// create service
	vkClient := vk.New(&cfg.AppCfg.VKConfig)
	svc := service.New(pgRepo, redisRepo, vkClient, &cfg.AppCfg.Service)

	// create handlers
	tokenHandler := api.New(svc)
	securityHandler := securityApi.New(svc)

	// create router
	mux, err := app.initRouter(tokenHandler, securityHandler)
	if err != nil {
		return fmt.Errorf("failed to create http ogen server: %v", err)
	}

	// create main http server
	apiSrv := app.createHttpServer(context.Background(), mux)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// run repository scheduler
	app.logger.Infof("starting repository scheduler with cfg=%+v", app.cfg.AppCfg.SchedulerConfig)

	if err := pgRepo.StartScheduler(ctx, app.logger, &app.cfg.AppCfg.SchedulerConfig); err != nil {
		return fmt.Errorf("failed to start repository scheduler: %w", err)
	}

	g.Go(func() error {
		<-ctx.Done()
		app.logger.Infof("shutting down repository scheduler...")
		return pgRepo.CloseScheduler()
	})

	// run http server
	g.Go(func() error {
		app.logger.Infof("http api server is running on %s", apiSrv.GetAddr())
		return apiSrv.Run(ctx)
	})

	// wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		app.logger.Error("server error: %v", err)
		return err
	}

	return nil
}

// initLogger create new Logger based on config
func (app *App) initLogger() {
	logCfg := app.cfg.AppCfg.LoggerCfg
	opts := []zap.Option{}

	// setup opts
	if len(logCfg.LogLevel) != 0 {
		opts = append(opts, zap.LogLvl(logCfg.LogLevel))
	}

	if len(logCfg.Mode) != 0 {
		opts = append(opts, zap.Mode(logCfg.Mode))
	}

	if len(logCfg.Encoding) != 0 {
		opts = append(opts, zap.Encoding(logCfg.Encoding))
	}

	if len(logCfg.LogPath) != 0 {
		opts = append(opts, zap.LogPath(logCfg.LogPath))
	}

	if logCfg.MaxLogSize != 0 {
		opts = append(opts, zap.MaxLogSize(logCfg.MaxLogSize))
	}

	if logCfg.MaxLogBackups != 0 {
		opts = append(opts, zap.MaxLogBackups(logCfg.MaxLogBackups))
	}

	if logCfg.MaxLogAge != 0 {
		opts = append(opts, zap.MaxLogAge(logCfg.MaxLogAge))
	}

	// set logger
	app.logger = zap.New(opts...)
}

// initRouter creates router using ogen
func (app *App) initRouter(handler oas.Handler, securityHandler oas.SecurityHandler) (http.Handler, error) {
	// create custom handlers
	notFoundHandler := func(w http.ResponseWriter, r *http.Request) {
		httputils.NotFoundJSON(w)
	}

	errorHandler := func(_ context.Context, w http.ResponseWriter, r *http.Request, err error) {
		app.logger.Errorf("API ERROR: %v", err)
		httputils.BadRequestJSON(w)
	}

	// create ogen http server
	srv, err := oas.NewServer(
		handler,
		securityHandler,
		oas.WithNotFound(notFoundHandler),
		oas.WithErrorHandler(errorHandler),
	)
	if err != nil {
		return nil, err
	}

	// apply middlewares
	middlewares := app.prepareMiddlewares()

	var mux http.Handler = srv
	for _, m := range middlewares {
		mux = m(mux)
	}

	return mux, nil
}

// prepareMiddlewares generates all middleware to be used in app
func (app *App) prepareMiddlewares() []middleware {
	// prepare recover middleware
	recoverMw := func(next http.Handler) http.Handler {
		return mw.WithRecover(app.logger, next)
	}

	// prepare log mw
	logMw := func(next http.Handler) http.Handler {
		return mw.WithLogging(app.logger, next)
	}

	// prepare cors mw
	cc := app.cfg.AppCfg.CORSConfig

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowCredentials = cc.AllowCredentials // ok, default value is false in both cases
	corsCfg.AllowOrigins = cc.AllowOrigins         // also ok, in both cases defautl is an empty slice

	if len(cc.AllowHeaders) != 0 {
		corsCfg.AllowHeaders = cc.AllowHeaders
	}

	if len(cc.AllowMethods) != 0 {
		corsCfg.AllowMethods = cc.AllowMethods
	}

	if cc.MaxAge != 0 {
		corsCfg.MaxAge = cc.MaxAge
	}

	corsMw := func(next http.Handler) http.Handler {
		return cors.WithCORS(corsCfg, next)
	}

	return []middleware{
		corsMw,
		mw.WithContentTypeJSON,
		logMw,
		recoverMw,
	}
}

// createHttpServer creates new httpserver.Server instance using app config
func (app *App) createHttpServer(ctx context.Context, mux http.Handler) *httpserver.Server {
	opts := []httpserver.Option{}

	httpCfg := app.cfg.AppCfg.HttpCfg

	if httpCfg.ReadTimeout != 0 {
		opts = append(opts, httpserver.ReadTimeout(httpCfg.ReadTimeout))
	}

	if httpCfg.WriteTimeout != 0 {
		opts = append(opts, httpserver.WriteTimeout(httpCfg.WriteTimeout))
	}

	if httpCfg.ShutdownTimeout != 0 {
		opts = append(opts, httpserver.ShutdownTimeout(httpCfg.ShutdownTimeout))
	}

	if httpCfg.Port != "" || httpCfg.Host != "" {
		port := "8000"
		if httpCfg.Port != "" {
			port = httpCfg.Port
		}

		host := ""
		if httpCfg.Host != "" {
			host = httpCfg.Host
		}

		opts = append(opts, httpserver.Addr(host, port))
	}

	return httpserver.New(ctx, mux, opts...)
}

// getPostgresConn creates and returns a new Postgres connection
func (app *App) getPostgresConn() (*postgres.Postgres, error) {
	var opts []postgres.Option

	if app.cfg.PostgresCfg.MaxPoolSize != 0 {
		opts = append(opts, postgres.MaxPoolSize(app.cfg.PostgresCfg.MaxPoolSize))
	}

	if app.cfg.PostgresCfg.ConnAttempts != 0 {
		opts = append(opts, postgres.ConnAttempts(app.cfg.PostgresCfg.ConnAttempts))
	}

	if app.cfg.PostgresCfg.ConnTimeout != 0 {
		opts = append(opts, postgres.ConnTimeout(app.cfg.PostgresCfg.ConnTimeout))
	}

	return postgres.New(app.cfg.PostgresCfg.DSN(), opts...)
}
