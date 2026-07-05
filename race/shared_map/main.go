package main

import "fmt"

type Metrics struct {
	counts map[string]*Counter
}

type Counter struct {
	Value int
}

func NewMetrics() *Metrics {
	return &Metrics{counts: make(map[string]*Counter)}
}

func (m *Metrics) Record(route string) {
	counter, ok := m.counts[route]
	if !ok {
		counter = &Counter{}
		m.counts[route] = counter
	}
	counter.Value++
}

func (m *Metrics) Snapshot() map[string]int {
	snapshot := make(map[string]int, len(m.counts))
	for route, counter := range m.counts {
		snapshot[route] = counter.Value
	}
	return snapshot
}

func main() {
	metrics := NewMetrics()
	metrics.Record("/checkout")

	fmt.Printf("checkout requests: %d\n", metrics.Snapshot()["/checkout"])
}
