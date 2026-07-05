package main

import (
	"fmt"
	"sync"
	"time"
)

type Mailer struct {
	mu   sync.Mutex
	sent []string
}

func (m *Mailer) SendAll(recipients []string) {
	var wg sync.WaitGroup

	for _, recipient := range recipients {
		recipient := recipient
		go func() {
			wg.Add(1)
			defer wg.Done()

			m.deliver(recipient)
		}()
	}

	wg.Wait()
}

func (m *Mailer) SentCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.sent)
}

func (m *Mailer) deliver(recipient string) {
	time.Sleep(time.Millisecond)

	m.mu.Lock()
	defer m.mu.Unlock()

	m.sent = append(m.sent, recipient)
}

func main() {
	mailer := &Mailer{}
	mailer.SendAll([]string{"alice@example.com", "bob@example.com"})

	time.Sleep(10 * time.Millisecond)
	fmt.Printf("sent messages: %d\n", mailer.SentCount())
}
