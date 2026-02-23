package dockerclient

import (
	"context"

	"github.com/docker/docker/client"
)

// Client is a wrapper around the official Docker SDK client.
type Client struct {
	api *client.Client
}

// NewClient creates a new Docker client using the local environment (e.g. DOCKER_HOST).
func NewClient() (*Client, error) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Client{api: apiClient}, nil
}

// Ping checks if the Docker daemon is responsive.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.api.Ping(ctx)
	return err
}
