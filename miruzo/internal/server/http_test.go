package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/mntone/miruzo-core/miruzo/internal/config"
)

func pickTestAddr(t testing.TB) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen test address: %v", err)
	}
	addr := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("close test listener: %v", err)
	}
	return addr
}

func waitServerReady(url string) error {
	client := &http.Client{Timeout: 50 * time.Millisecond}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}

	return fmt.Errorf("server did not become ready: %s", url)
}

func TestNewHTTPServerAppliesConfig(t *testing.T) {
	conf := config.ServerConfig{
		Port:              8080,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      4 * time.Second,
		IdleTimeout:       5 * time.Second,
		MaxHeaderBytes:    1024,
	}

	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	server := NewHTTPServer(handler, conf)

	if server.Addr != ":8080" {
		t.Fatalf("Addr = %s, want :8080", server.Addr)
	}
	if server.Handler == nil {
		t.Fatal("Handler = nil, want non-nil")
	}
	if server.ReadHeaderTimeout != conf.ReadHeaderTimeout {
		t.Fatalf("ReadHeaderTimeout = %s, want %s", server.ReadHeaderTimeout, conf.ReadHeaderTimeout)
	}
	if server.ReadTimeout != conf.ReadTimeout {
		t.Fatalf("ReadTimeout = %s, want %s", server.ReadTimeout, conf.ReadTimeout)
	}
	if server.WriteTimeout != conf.WriteTimeout {
		t.Fatalf("WriteTimeout = %s, want %s", server.WriteTimeout, conf.WriteTimeout)
	}
	if server.IdleTimeout != conf.IdleTimeout {
		t.Fatalf("IdleTimeout = %s, want %s", server.IdleTimeout, conf.IdleTimeout)
	}
	if server.MaxHeaderBytes != conf.MaxHeaderBytes {
		t.Fatalf("MaxHeaderBytes = %d, want %d", server.MaxHeaderBytes, conf.MaxHeaderBytes)
	}
}

func TestRunWithContextReturnsListenError(t *testing.T) {
	httpServer := &http.Server{
		Addr:    "127.0.0.1:invalid",
		Handler: http.NewServeMux(),
	}

	err := runWithContext(context.Background(), httpServer, 100*time.Millisecond)
	if err == nil {
		t.Fatal("runWithContext() error = nil, want non-nil")
	}
}

func TestRunWithContextShutsDownOnContextCancel(t *testing.T) {
	addr := pickTestAddr(t)
	httpServer := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
	}

	notifyContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	readyError := make(chan error, 1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		readyError <- waitServerReady("http://" + addr + "/ready")
		cancel()
	}()

	err := runWithContext(notifyContext, httpServer, time.Second)
	if err != nil {
		t.Fatalf("runWithContext() error = %v, want nil", err)
	}
	if err := <-readyError; err != nil {
		t.Fatal(err)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("cancel goroutine did not finish")
	}
}

func TestRunWithContextReturnsShutdownDeadlineExceeded(t *testing.T) {
	addr := pickTestAddr(t)
	requestStarted := make(chan struct{})
	release := make(chan struct{})

	httpServer := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/ready" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			close(requestStarted)
			<-release
			w.WriteHeader(http.StatusNoContent)
		}),
	}

	notifyContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	readyError := make(chan error, 1)
	go func() {
		readyError <- waitServerReady("http://" + addr + "/ready")

		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get("http://" + addr + "/block")
		if err == nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}
	}()

	go func() {
		select {
		case <-requestStarted:
			cancel()
		case <-time.After(2 * time.Second):
			cancel()
		}
	}()

	err := runWithContext(notifyContext, httpServer, 10*time.Millisecond)
	close(release)

	if errReady := <-readyError; errReady != nil {
		t.Fatal(errReady)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("runWithContext() error = %v, want context deadline exceeded", err)
	}
}
