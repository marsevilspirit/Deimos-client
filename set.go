package deimosclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SetOptions struct {
	ttl       time.Duration
	dir       bool
	prevExist *bool
}

func newSetOptions(options []SetOption) *SetOptions {
	setOpts := SetOptions{}

	for _, opt := range options {
		opt.applyToSet(&setOpts)
	}

	return &setOpts
}

func (c *Client) Set(ctx context.Context, key, value string, opts ...SetOption) (*Response, error) {
	setOpts := newSetOptions(opts)

	URL := c.buildURL(key)
	query := url.Values{}

	if setOpts.dir {
		query.Set("dir", "true")
	} else {
		query.Set("value", value)
	}

	if setOpts.ttl > 0 {
		query.Set("ttl", fmt.Sprintf("%d", int64(setOpts.ttl.Seconds())))
	}

	if setOpts.prevExist != nil {
		if *setOpts.prevExist {
			query.Set("prevExists", "true")
		} else {
			query.Set("prevExists", "false")
		}
	}

	body := strings.NewReader(query.Encode())
	req, err := http.NewRequestWithContext(ctx, "PUT", URL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.doRequest(req)
}
