package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/config"
)

func NewHTTPServer(handler http.Handler, conf config.ServerConfig) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", conf.Port),
		Handler:           handler,
		ReadHeaderTimeout: conf.ReadHeaderTimeout,
		ReadTimeout:       conf.ReadTimeout,
		WriteTimeout:      conf.WriteTimeout,
		IdleTimeout:       conf.IdleTimeout,
		MaxHeaderBytes:    conf.MaxHeaderBytes,
	}
}

func runWithContext(
	notifyContext context.Context,
	httpServer *http.Server,
	shutdownTimeout time.Duration,
) error {
	serverError := make(chan error, 1)
	go func() {
		serverError <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverError:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil

	case <-notifyContext.Done():
	}

	shutdownContext, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownContext); err != nil {
		closeError := httpServer.Close()
		if closeError != nil {
			return errors.Join(err, closeError)
		}
		return err
	}

	if err := <-serverError; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func Run(
	httpServer *http.Server,
	shutdownTimeout time.Duration,
) error {
	shutdownSignalContext, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	return runWithContext(shutdownSignalContext, httpServer, shutdownTimeout)
}
