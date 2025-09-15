package keyscore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func WithBaseURL(base string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(base, "/")
	}
}

func WithAPIKey(key string) Option {
	return func(c *Client) {
		c.apiKey = key
	}
}

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL:    "https://api.keysco.re",
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: 60 * time.Second}
	}
	return c
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("keysco.re API error: status=%d message=%q", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("keysco.re API error: status=%d", e.StatusCode)
}

func (c *Client) newRequest(ctx context.Context, method, path string, q url.Values, body any) (*http.Request, error) {
	u := c.baseURL + path
	if len(q) > 0 {
		u = u + "?" + q.Encode()
	}

	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		r = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, r)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	return req, nil
}

func (c *Client) doJSON(req *http.Request, v any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var er struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&er)
		return &APIError{StatusCode: resp.StatusCode, Message: er.Error}
	}
	if v == nil {
		return nil
	}
	dec := json.NewDecoder(resp.Body)
	// Be lenient to allow forward-compatible responses; ignore unknown fields.
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func (c *Client) doRaw(req *http.Request) (io.ReadCloser, string, int64, string, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", 0, "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		var er struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&er)
		return nil, "", 0, "", &APIError{StatusCode: resp.StatusCode, Message: er.Error}
	}
	ct := resp.Header.Get("Content-Type")
	cl := int64(0)
	if v := resp.Header.Get("Content-Length"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			cl = n
		}
	}
	cd := resp.Header.Get("Content-Disposition")
	return resp.Body, ct, cl, cd, nil
}
