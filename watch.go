package deimosclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// WatchOptions contains all optional parameters for a Watch operation.
type WatchOptions struct {
	recursive bool   // Whether to recursively watch a directory.
	waitIndex uint64 // The index to start waiting from.
}

func newWatchOptions(options []WatchOption) *WatchOptions {
	watchOpts := WatchOptions{}

	for _, opt := range options {
		opt.applyToWatch(&watchOpts)
	}

	return &watchOpts
}

// Watch monitors a key for changes.
// It returns a read-only Response channel. When a change occurs, the response is sent through the channel.
// The caller must use a context to control the Watcher's lifecycle. When the context is canceled, the watcher will stop and close the channel.
func (c *Client) Watch(ctx context.Context, key string, opts ...WatchOption) <-chan *Response {
	respChan := make(chan *Response, 1)

	watchOpts := newWatchOptions(opts)

	go c.watcher(ctx, key, watchOpts, respChan)

	return respChan
}

// watcher is the long-polling loop that runs in the background.
func (c *Client) watcher(ctx context.Context, key string, opts *WatchOptions, respChan chan<- *Response) {
	// Ensure the channel is closed on goroutine exit, which is how the caller is notified that the watch has ended.
	defer close(respChan)

	for {
		query := url.Values{}
		query.Set("wait", "true")

		if opts.recursive {
			query.Set("recursive", "true")
		}
		// If waitIndex is greater than 0, add it to the query.
		// This is the core of the loop: after each event, we use the new index + 1 to make the next request.
		if opts.waitIndex > 0 {
			query.Set("waitIndex", fmt.Sprintf("%d", opts.waitIndex))
		}

		URL := c.buildURL(key) + "?" + query.Encode()

		// Create the request and pass the parent context into it.
		// If the parent context is canceled, the request here will fail immediately,
		// allowing for a graceful exit from the loop.
		req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
		if err != nil {
			// This should not normally happen, but if it does, we cannot continue.
			// We could send an error to an error channel here, but for simplicity, we just return.
			fmt.Printf("watcher: failed to create request: %v\n", err)
			return
		}

		// Execute the long-polling request.
		resp, err := c.httpClient.Do(req)
		if err != nil {
			// When the context is canceled, an error will be received here. This is the expected way to exit.
			// Check the context's error; if it's canceled, exit silently.
			select {
			case <-ctx.Done():
				// Context was canceled, this is the intended way to stop the watcher.
			default:
				fmt.Printf("watcher: failed to execute HTTP request: %v\n", err)
			}
			return
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("watcher: failed to read response body: %v\n", err)
			return
		}

		var deimosResp Response
		if err := json.Unmarshal(respBody, &deimosResp); err != nil {
			fmt.Printf("watcher: failed to unmarshal JSON response: %v\n", err)
			continue // It might be a temporary issue, so try the next iteration.
		}

		// If deimos returns an API error (e.g., key not found), stop the watcher.
		if deimosResp.ErrorCode != 0 {
			fmt.Printf("watcher: deimos API error, stopping watch: [%d] %s\n", deimosResp.ErrorCode, deimosResp.Message)
			return
		}

		// Update waitIndex so the next request can get the next event.
		// This is the key to achieving continuous watching!
		opts.waitIndex = deimosResp.Node.ModifiedIndex + 1

		// Send the response to the channel.
		// At the same time, check if the context has been canceled in case the caller
		// has already exited while we are trying to send.
		select {
		case respChan <- &deimosResp:
			// Sent successfully.
		case <-ctx.Done():
			// The context was canceled while we were trying to send.
			return
		}
	}
}
