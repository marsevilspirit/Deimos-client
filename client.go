package deimosclient

import (
	"net/http"
	"time"
)

type Client struct {
	cluster    *Cluster
	httpClient *http.Client
}

// NewClient create a basic client that is configured to be used
// with the given machine list.
func NewClient(endpoints []string) *Client {
	return &Client{
		cluster: NewCluster(endpoints),
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}
