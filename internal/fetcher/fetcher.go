package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"yandex-lms/internal/models"
)

type Fetcher struct {
	BaseURL string
	Client  *http.Client
}

func New(baseURL string) *Fetcher {
	timeout := time.Second * 10
	return &Fetcher{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (f *Fetcher) Fetch(ctx context.Context, page string) ([]models.Item, error) {
	url := fmt.Sprintf("%s%s", f.BaseURL, page)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer resp.Body.Close()

	w, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	b := &models.HHResponse{}
	err = json.Unmarshal(w, &b)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return b.HHItems, nil
}
