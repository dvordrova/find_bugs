package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/dvordrova/find_bugs/golangci/http_response_body_leak/internal/catalog"
)

func main() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"items":[]}`))
	}))
	defer server.Close()

	status, err := catalog.FetchStatus(context.Background(), server.Client(), server.URL)
	fmt.Printf("status=%d err=%v\n", status, err)
	fmt.Println("run make lint to see the bodyclose report")
}
