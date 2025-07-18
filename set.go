package deimosclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) Set(ctx context.Context, key, value string) (*Response, error) {
	fullURL := c.buildURL(key)

	formData := url.Values{}
	formData.Set("value", value)
	body := strings.NewReader(formData.Encode())

	req, err := http.NewRequestWithContext(ctx, "PUT", fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("create deimos request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.doRequest(req)
}
