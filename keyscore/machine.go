package keyscore

import (
	"context"
	"net/http"
	"net/url"
)

// MachineInfo fetches machine info by UUID.
func (c *Client) MachineInfo(ctx context.Context, uuid string) (*MachineInfo, error) {
	q := url.Values{"uuid": {uuid}}
	req, err := c.newRequest(ctx, http.MethodGet, "/machineinfo", q, nil)
	if err != nil {
		return nil, err
	}
	var out MachineInfo
	if err := c.doJSON(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Download retrieves a file by UUID. If filePath is empty, downloads the full archive.
// Caller MUST close the returned Body.
func (c *Client) Download(ctx context.Context, uuid string, filePath string) (*DownloadResult, error) {
	q := url.Values{"uuid": {uuid}}
	if filePath != "" {
		q.Set("file", filePath)
	}
	req, err := c.newRequest(ctx, http.MethodGet, "/download", q, nil)
	if err != nil {
		return nil, err
	}
	body, ct, cl, cd, err := c.doRaw(req)
	if err != nil {
		return nil, err
	}
	return &DownloadResult{Body: body, ContentType: ct, ContentLength: cl, ContentDisposition: cd}, nil
}
