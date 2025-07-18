package deimosclient

import "fmt"

func (c *Client) buildURL(key string) string {
	endpoint := c.cluster.pick()
	return fmt.Sprintf("%s/keys%s", endpoint, key)
}
