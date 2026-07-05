package main

import "fmt"

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type StreamConsumer struct {
	noCopy noCopy
	topic  string
}

func NewStreamConsumer(topic string) StreamConsumer {
	return StreamConsumer{topic: topic}
}

func (c StreamConsumer) Topic() string {
	return c.topic
}

func main() {
	consumer := NewStreamConsumer("payments")

	fmt.Printf("consumer topic: %s\n", consumer.Topic())
}
