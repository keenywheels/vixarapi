package app

import (
	"context"
	"fmt"
	"net/http"

	oas "github.com/keenywheels/backend/internal/api/v1"
	httpdomain "github.com/keenywheels/backend/internal/example_domain/delivery/http/v1"
	"github.com/keenywheels/backend/pkg/cors"
	"github.com/keenywheels/backend/pkg/httpserver"
	"github.com/keenywheels/backend/pkg/httputils"
	"github.com/keenywheels/backend/pkg/logger"
	"github.com/keenywheels/backend/pkg/logger/zap"
	mw "github.com/keenywheels/backend/pkg/middleware"
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

	// create mux using ogen
	mux, err := app.initRouter()
	if err != nil {
		return fmt.Errorf("failed to create http ogen server: %v", err)
	}

	// create and run main http server
	apiSrv := app.createHttpServer(context.Background(), mux)

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		app.logger.Infof("http api server is running on %s", apiSrv.GetAddr())
		return apiSrv.Run(ctx)
	})

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
func (app *App) initRouter() (http.Handler, error) {
	// create handler
	domainHandler := httpdomain.NewController()

	// create custom handlers
	notFoundHandler := func(w http.ResponseWriter, r *http.Request) {
		httputils.NotFoundJSON(w)
	}

	errorHandler := func(_ context.Context, w http.ResponseWriter, r *http.Request, err error) {
		httputils.InternalErrorJSON(w)
	}

	// create ogen http server
	srv, err := oas.NewServer(domainHandler,
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

	return srv, nil
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

	if httpCfg.Port != "" && httpCfg.Host != "" {
		port := "8080"
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
