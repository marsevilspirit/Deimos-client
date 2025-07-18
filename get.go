package deimosclient

import (
	"context"
	"fmt"
	"net/http"
)

// TODO: add get options
type GetOptions struct{}

func (c *Client) Get(ctx context.Context, key string) (*Response, error) {
	URL := c.buildURL(key)
	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		return nil, fmt.Errorf("create deimos request failed: %w", err)
	}

	return c.doRequest(req)
}
