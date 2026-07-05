package main

import (
	"context"
	"fmt"
	"time"
)

type Profile struct {
	ID   string
	Name string
}

type ProfileGateway struct{}

func (ProfileGateway) Lookup(ctx context.Context, id string) (Profile, error) {
	select {
	case <-time.After(time.Millisecond):
		return Profile{ID: id, Name: "Alice"}, nil
	case <-ctx.Done():
		return Profile{}, ctx.Err()
	}
}

func LoadProfile(ctx context.Context, gateway ProfileGateway, id string) (Profile, error) {
	ctx, _ = context.WithTimeout(ctx, 500*time.Millisecond)

	return gateway.Lookup(ctx, id)
}

func main() {
	profile, err := LoadProfile(context.Background(), ProfileGateway{}, "profile-001")
	if err != nil {
		fmt.Printf("load profile failed: %v\n", err)
		return
	}

	fmt.Printf("loaded profile %s for %s\n", profile.ID, profile.Name)
}
