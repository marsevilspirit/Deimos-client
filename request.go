package deimosclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) doRequest(req *http.Request) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do http request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("deimos server error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var deimosResp Response
	if err := json.Unmarshal(respBody, &deimosResp); err != nil {
		// For 4xx errors, the response might be in a different format
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return nil, fmt.Errorf("deimos client error (HTTP %d): %s", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("umarshal json error: %w", err)
	}

	// check errcode
	if deimosResp.ErrorCode != 0 {
		return nil, fmt.Errorf("deimos API err: [%d] %s", deimosResp.ErrorCode, deimosResp.Message)
	}

	return &deimosResp, nil
}
