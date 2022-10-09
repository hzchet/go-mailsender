package application

import (
	"context"

	"gitlab.com/golang-hse-2022/team1/mailsender/pkg/infra/logger"
	"gitlab.com/golang-hse-2022/team1/mailsender/pkg/infra/metrics"
	"gitlab.com/golang-hse-2022/team1/mailsender/pkg/infra/probes"
)

type App struct {
	l             logger.Logger
	shutdownFuncs []func(ctx context.Context) error
}

func New(l logger.Logger) *App {
	return &App{
		l: l,
	}
}

func (app *App) Start() error {
	probes, _ := probes.New(app.l)
	probes.SetStarted()
	err := probes.Start()
	if err != nil {
		app.l.Sugar().Fatalf("probes start failed: %s", err.Error())
	}

	me, err := metrics.New()
	if err != nil {
		app.l.Sugar().Fatalf("Metrics init failed: %s", err.Error())
	}
	app.shutdownFuncs = append(app.shutdownFuncs, me.Stop)

	probes.SetReady()

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	var err error
	for i := len(a.shutdownFuncs) - 1; i >= 0; i-- {
		err = a.shutdownFuncs[i](ctx)
		if err != nil {
			a.l.Sugar().Error(err)
		}
	}

	a.l.Info("app stopped")

	return nil
}
