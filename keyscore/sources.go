package keyscore

import (
	"context"
	"net/http"
)

func (c *Client) Sources(ctx context.Context) (*SourcesResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/sources", nil, nil)
	if err != nil {
		return nil, err
	}
	var out SourcesResponse
	if err := c.doJSON(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
