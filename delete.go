package deimosclient

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) Delete(ctx context.Context, key string) (*Response, error) {
	fullURL := c.buildURL(key)
	req, err := http.NewRequestWithContext(ctx, "DELETE", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	return c.doRequest(req)
}
