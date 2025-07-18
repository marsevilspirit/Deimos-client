package deimosclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type GetOptions struct {
	recursive bool
}

func newGetOptions(options []GetOption) *GetOptions {
	getOpts := GetOptions{}

	for _, opt := range options {
		opt.applyToGet(&getOpts)
	}

	return &getOpts
}

func (c *Client) Get(ctx context.Context, key string, opts ...GetOption) (*Response, error) {
	getOpts := newGetOptions(opts)

	URL := c.buildURL(key)
	query := url.Values{}

	if getOpts.recursive {
		query.Set("recursive", "true")
	}

	if len(query) > 0 {
		URL += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		return nil, fmt.Errorf("create deimos request failed: %w", err)
	}

	return c.doRequest(req)
}
