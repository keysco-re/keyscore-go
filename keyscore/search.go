package keyscore

import (
	"context"
	"net/http"
)

func (c *Client) Search(ctx context.Context, body SearchRequest) (*SearchResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/search", nil, body)
	if err != nil {
		return nil, err
	}
	var out SearchResponse
	if err := c.doJSON(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
