package main

import (
	"testing"
	"time"
)

func TestSendAllEventuallySendsMessages(t *testing.T) {
	mailer := &Mailer{}
	mailer.SendAll([]string{"alice@example.com", "bob@example.com"})

	deadline := time.After(100 * time.Millisecond)
	for {
		if mailer.SentCount() == 2 {
			return
		}

		select {
		case <-deadline:
			t.Fatalf("expected 2 sent messages, got %d", mailer.SentCount())
		default:
			time.Sleep(time.Millisecond)
		}
	}
}
