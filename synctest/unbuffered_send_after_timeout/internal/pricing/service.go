package pricing

import (
	"context"
	"time"
)

type Price struct {
	SKU   string
	Cents int
}

type Client interface {
	Lookup(sku string) (Price, error)
}

type FixedDelayClient struct {
	Delay time.Duration
	Price Price
	Err   error
}

func (c FixedDelayClient) Lookup(string) (Price, error) {
	time.Sleep(c.Delay)
	return c.Price, c.Err
}

type Service struct {
	client  Client
	timeout time.Duration
}

func NewService(client Client, timeout time.Duration) Service {
	return Service{
		client:  client,
		timeout: timeout,
	}
}

func (s Service) Lookup(ctx context.Context, sku string) (Price, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	results := make(chan lookupResult)
	go func() {
		price, err := s.client.Lookup(sku)
		results <- lookupResult{price: price, err: err}
	}()

	select {
	case result := <-results:
		return result.price, result.err
	case <-ctx.Done():
		return Price{}, ctx.Err()
	}
}

type lookupResult struct {
	price Price
	err   error
}
