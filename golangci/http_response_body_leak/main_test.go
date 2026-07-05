package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dvordrova/find_bugs/golangci/http_response_body_leak/internal/catalog"
)

func TestFetchStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	status, err := catalog.FetchStatus(context.Background(), server.Client(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if status != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", status, http.StatusAccepted)
	}
}
