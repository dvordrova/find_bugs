package main

import (
	"testing"

	"github.com/dvordrova/find_bugs/concurrency/select_priority_assumption/internal/dispatcher"
)

func TestOrdinaryOnlyHighPriorityReady(t *testing.T) {
	high := make(chan dispatcher.Job, 1)
	low := make(chan dispatcher.Job, 1)
	high <- dispatcher.Job{Queue: "high", ID: "urgent-1"}

	got := dispatcher.Next(high, low)
	if got.Queue != "high" {
		t.Fatalf("selected %s queue, want high", got.Queue)
	}
}

func TestSelectDoesNotGuaranteePriority(t *testing.T) {
	for range 1024 {
		high := make(chan dispatcher.Job, 1)
		low := make(chan dispatcher.Job, 1)
		high <- dispatcher.Job{Queue: "high", ID: "urgent-1"}
		low <- dispatcher.Job{Queue: "low", ID: "batch-1"}

		got := dispatcher.Next(high, low)
		if got.Queue == "low" {
			t.Fatal("select chose low-priority job while high-priority job was ready")
		}
	}
}
