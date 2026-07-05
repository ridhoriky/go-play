package app

import (
	"ne-project/src/internal/config/appconfig"
	"ne-project/src/internal/config/grace"
	"ne-project/src/internal/config/logger"
	"ne-project/src/internal/config/resource"
	"ne-project/src/internal/config/router"
	"ne-project/src/internal/config/server"
	"ne-project/src/internal/config/token"

	"github.com/rs/zerolog"
)

type App struct {
	cfg      *appconfig.Config
	log      *zerolog.Logger
	res      *resource.Resources
	tokenSvc token.TokenServiceItf
}

func NewApp(cfg *appconfig.Config) (*App, error) {
	log := logger.InitLogger(&cfg.Logger)
	zerolog.DefaultContextLogger = log

	res, err := resource.InitResources(log, cfg)
	if err != nil {
		return nil, err
	}

	tokenSvc, err := token.NewTokenService(&cfg.Token)
	if err != nil {
		res.Close(log)
		return nil, err
	}

	return &App{
		cfg:      cfg,
		log:      log,
		res:      res,
		tokenSvc: tokenSvc,
	}, nil
}

func (a *App) Start() {
	r := router.SetupRouter(a.log, a.cfg, a.res, a.tokenSvc)
	srv := server.InitServer(a.cfg, r)

	grace.WaitForShutdown(a.log, srv, &a.cfg.App)
}

func (a *App) Close() {
	a.res.Close(a.log)
}

func Run(cfg *appconfig.Config) {
	app, err := NewApp(cfg)
	if err != nil {
		panic(err)
	}
	defer app.Close()

	app.Start()
}
