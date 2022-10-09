package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"gitlab.com/golang-hse-2022/team1/mailsender/internal/ports"
	"gitlab.com/golang-hse-2022/team1/mailsender/pkg/infra/logger"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
)

type Adapter struct {
	s          *http.Server
	l          net.Listener
	mailsender ports.Mailsender
}

type Config struct {
	Port int `env:"HTTP_PORT" envDefault:"3000"`
}

func New(mailsender ports.Mailsender, log logger.Logger) (*Adapter, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse server http adapter configuration failed: %w", err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("server start failed: %w", err)
	}

	router := gin.Default()
	server := http.Server{
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	a := Adapter{
		s:          &server,
		l:          l,
		mailsender: mailsender,
	}
	initRouter(&a, router, log)

	return &a, nil
}

func (a *Adapter) Start() error {
	var err error
	go func() {
		err = a.s.Serve(a.l)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			err = fmt.Errorf("server start failed: %w", err)
		}
		err = nil
	}()

	if err != nil {
		return err
	}
	return nil
}

func (a *Adapter) Stop(ctx context.Context) error {
	var (
		err  error
		once sync.Once
	)
	once.Do(func() {
		err = a.s.Shutdown(ctx)
	})
	return err
}
