package rest

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/longwavee/effective-mobile-test/internal/config"
)

type (
	HTTPServer struct {
		srv             *http.Server
		shutdownTimeout time.Duration

		once sync.Once

		errCh  chan error
		doneCh chan struct{}
	}
)

func NewHTTPServer(cfg *config.HTTPServer, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		srv: &http.Server{
			Addr:              ":" + cfg.Port,
			Handler:           handler,
			ReadTimeout:       cfg.ReadTimeout,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			MaxHeaderBytes:    cfg.MaxHeaderBytes,
		},
		shutdownTimeout: cfg.ShutdownTimeout,

		errCh:  make(chan error, 1),
		doneCh: make(chan struct{}),
	}
}

func (s *HTTPServer) Start() error {
	l, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen %s: %w", s.srv.Addr, err)
	}

	s.once.Do(func() {
		go func() {
			err := s.srv.Serve(l)
			if errors.Is(err, http.ErrServerClosed) {
				err = nil
			}
			s.errCh <- err

			close(s.doneCh)
		}()
	})

	return nil
}

func (s *HTTPServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	<-s.doneCh

	return <-s.errCh
}
