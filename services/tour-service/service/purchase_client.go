package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type PurchaseClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewPurchaseClient() *PurchaseClient {
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("PURCHASE_SERVICE_URL")), "/")
	if baseURL == "" {
		baseURL = "http://localhost:8087"
	}
	return &PurchaseClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *PurchaseClient) IsPurchased(ctx context.Context, touristID string, tourID string) (bool, error) {
	url := fmt.Sprintf("%s/api/purchases/tokens/%s/tour/%s", c.baseURL, touristID, tourID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return false, fmt.Errorf("purchase service returned status %d", res.StatusCode)
	}

	var payload struct {
		Purchased bool `json:"purchased"`
	}
	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return false, err
	}
	return payload.Purchased, nil
}
