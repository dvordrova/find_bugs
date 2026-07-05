package leases

import (
	"context"
	"time"
)

type Manager struct {
	ttl time.Duration
}

type Lease struct {
	ID     string
	done   <-chan struct{}
	cancel context.CancelFunc
}

func NewManager(ttl time.Duration) Manager {
	return Manager{ttl: ttl}
}

func (m Manager) Start(ctx context.Context, id string) Lease {
	leaseCtx, cancel := context.WithTimeout(ctx, m.ttl*2) // BUG: lease lives twice as long as the configured TTL.
	return Lease{
		ID:     id,
		done:   leaseCtx.Done(),
		cancel: cancel,
	}
}

func (l Lease) Active() bool {
	select {
	case <-l.done:
		return false
	default:
		return true
	}
}

func (l Lease) Done() <-chan struct{} {
	return l.done
}

func (l Lease) Stop() {
	l.cancel()
}
