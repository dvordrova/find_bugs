package catalog

import (
	"context"
	"fmt"
	"net/http"
)

func FetchStatus(ctx context.Context, client *http.Client, url string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("fetch catalog: %w", err)
	}

	return resp.StatusCode, nil
}
