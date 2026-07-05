package main

import (
	"context"
	"testing"
	"testing/synctest"
	"time"
)

func TestOrdinaryTestOnlyChecksAfterCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sink := NewAuditSink()

	NotifyWhenCanceled(ctx, sink, "tenant-42")
	cancel()

	waitCtx, stop := context.WithTimeout(context.Background(), time.Second)
	defer stop()

	if !sink.Wait(waitCtx) {
		t.Fatal("after cancel: audit record was not written")
	}
}

func TestSynctestChecksBeforeCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		sink := NewAuditSink()

		NotifyWhenCanceled(ctx, sink, "tenant-42")
		synctest.Wait()
		if got := sink.Count(); got != 0 {
			t.Fatalf("before cancel: audit records = %d, want 0", got)
		}

		cancel()
		synctest.Wait()
		if got := sink.Count(); got != 1 {
			t.Fatalf("after cancel: audit records = %d, want 1", got)
		}
	})
}
