package main

import (
	"context"
	"testing"
)

func TestLoadProfile(t *testing.T) {
	profile, err := LoadProfile(context.Background(), ProfileGateway{}, "profile-001")
	if err != nil {
		t.Fatalf("load profile: %v", err)
	}

	if profile.Name != "Alice" {
		t.Fatalf("expected Alice, got %q", profile.Name)
	}
}
