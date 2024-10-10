package fetcher

import (
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
	return &Fetcher{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 10 * time.Minute,
		},
	}
}

func (f *Fetcher) Fetch(page string) ([]models.Item, error) {
	url := fmt.Sprintf("%s%s", f.BaseURL, page)
	r, err := f.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	defer r.Body.Close()

	b := &models.HHResponse{}
	w, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	err = json.Unmarshal(w, &b)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return b.HHItems, nil
}
