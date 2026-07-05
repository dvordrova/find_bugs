package main

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/dvordrova/find_bugs/synctest/context_timeout_without_wall_clock/internal/leases"
)

func TestOrdinaryTestOnlyChecksLeaseStartsActive(t *testing.T) {
	manager := leases.NewManager(50 * time.Millisecond)
	lease := manager.Start(context.Background(), "order-42")
	defer lease.Stop()

	if !lease.Active() {
		t.Fatal("new lease is already expired")
	}
}

func TestSynctestChecksLeaseDeadline(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const ttl = 5 * time.Second
		manager := leases.NewManager(ttl)
		lease := manager.Start(t.Context(), "order-42")
		defer lease.Stop()

		time.Sleep(ttl - time.Nanosecond)
		synctest.Wait()
		if !lease.Active() {
			t.Fatal("before ttl: lease expired too early")
		}

		time.Sleep(time.Nanosecond)
		synctest.Wait()
		if lease.Active() {
			t.Fatalf("after ttl: lease is still active, want expired")
		}
	})
}
