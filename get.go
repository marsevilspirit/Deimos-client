package deimosclient

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) Get(ctx context.Context, key string) (*Response, error) {
	fullURL := c.buildURL(key)
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create deimos request failed: %w", err)
	}

	return c.doRequest(req)
}
