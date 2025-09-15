package keyscore

import (
	"context"
	"net/http"
)

// HashLookup performs a hash lookup for the provided terms.
func (c *Client) HashLookup(ctx context.Context, reqBody HashLookupRequest) (*HashLookupResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/hashlookup", nil, reqBody)
	if err != nil {
		return nil, err
	}
	var out HashLookupResponse
	if err := c.doJSON(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// IPLookup performs an IP lookup for the provided terms.
func (c *Client) IPLookup(ctx context.Context, reqBody IPLookupRequest) (*IPLookupResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/iplookup", nil, reqBody)
	if err != nil {
		return nil, err
	}
	var out IPLookupResponse
	if err := c.doJSON(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
