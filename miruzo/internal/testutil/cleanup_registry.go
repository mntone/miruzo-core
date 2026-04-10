package testutil

import (
	"errors"
	"sync"
)

// CleanupRegistry stores cleanup callbacks and runs them in reverse registration
// order.
type CleanupRegistry struct {
	mtx     sync.Mutex
	closed  bool
	closers []func() error
}

func (reg *CleanupRegistry) Register(closer func() error) {
	if closer == nil {
		return
	}

	reg.mtx.Lock()
	defer reg.mtx.Unlock()
	if reg.closed {
		return
	}

	reg.closers = append(reg.closers, closer)
}

func (reg *CleanupRegistry) CloseAll() error {
	reg.mtx.Lock()
	if reg.closed {
		reg.mtx.Unlock()
		return nil
	}

	reg.closed = true
	closers := reg.closers
	reg.closers = nil
	reg.mtx.Unlock()

	var joined error
	for i := len(closers) - 1; i >= 0; i-- {
		joined = errors.Join(joined, closers[i]())
	}
	return joined
}
