package main

import "testing"

func TestTopic(t *testing.T) {
	consumer := NewStreamConsumer("payments")

	if got := consumer.Topic(); got != "payments" {
		t.Fatalf("expected payments topic, got %q", got)
	}
}
