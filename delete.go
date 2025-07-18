package deimosclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type DeleteOptions struct {
	recursive bool
	dir       bool
}

func newDeleteOptions(options []DeleteOption) *DeleteOptions {
	deleteOpts := DeleteOptions{}

	for _, opt := range options {
		opt.applyToDelete(&deleteOpts)
	}

	return &deleteOpts
}

func (c *Client) Delete(ctx context.Context, key string, opts ...DeleteOption) (*Response, error) {
	deleteOpts := newDeleteOptions(opts)

	URL := c.buildURL(key)
	query := url.Values{}

	// Build the query string based on the options.
	if deleteOpts.recursive {
		query.Set("recursive", "true")
	}
	if deleteOpts.dir {
		query.Set("dir", "true")
	}

	// Append the query string to the URL if any options were provided.
	if len(query) > 0 {
		URL += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return c.doRequest(req)
}
