package deimosclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (c *Client) Set(ctx context.Context, key, value string, opts ...SetOption) (*Response, error) {
	options := &SetOptions{}

	for _, opt := range opts {
		opt(options)
	}

	fullURL := c.buildURL(key)

	formData := url.Values{}
	formData.Set("value", value)

	// handle ttl option
	if options.TTL > 0 {
		ttlInSeconds := int64(options.TTL.Seconds())
		formData.Set("ttl", fmt.Sprintf("%d", ttlInSeconds))
	}

	body := strings.NewReader(formData.Encode())

	req, err := http.NewRequestWithContext(ctx, "PUT", fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("create deimos request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.doRequest(req)
}

type SetOptions struct {
	TTL time.Duration
}

type SetOption func(*SetOptions)

func WithTTL(ttl time.Duration) SetOption {
	return func(opts *SetOptions) {
		opts.TTL = ttl
	}
}
