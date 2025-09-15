package keyscore

import (
	"context"
	"net/http"
)

// Health checks API service health.
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/health", nil, nil)
	if err != nil {
		return nil, err
	}
	var out HealthResponse
	if err := c.doJSON(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
