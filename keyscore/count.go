package keyscore

import (
	"context"
	"net/http"
)

// Count performs a basic count request.
func (c *Client) Count(ctx context.Context, body CountRequest) (*CountResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/count", nil, body)
	if err != nil {
		return nil, err
	}
	var out CountResponse
	if err := c.doJSON(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CountDetailed performs a detailed count request.
func (c *Client) CountDetailed(ctx context.Context, body CountRequest) (*DetailedCountResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/count/detailed", nil, body)
	if err != nil {
		return nil, err
	}
	var out DetailedCountResponse
	if err := c.doJSON(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
