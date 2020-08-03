package dockerop

import (
	"context"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// OpClient the client operating docker engine
type OpClient struct {
	client *client.Client
	// others
	Auth          types.AuthConfig
	Time          time.Time // set in RegistryLogin
	IdentityToken string    // returned in RegistryLogin
	Status        string    // returned in RegistryLogin
}

// NewOpClient return a pointer to OpClient
func NewOpClient(host string, httpcli *http.Client) (*OpClient, error) {
	cli, err := client.NewClient(host, "", httpcli, nil)
	if err != nil {
		return nil, err
	}
	var opCli = &OpClient{
		client: cli,
	}
	return opCli, nil
}

// Close close the connection to docker engine
func (c *OpClient) Close() error {
	return c.client.Close()
}

// Ping pings the server
func (c *OpClient) Ping(ctx context.Context) (types.Ping, error) {
	return c.client.Ping(ctx)
}

// Info returns the information about the docker engine
func (c *OpClient) Info(ctx context.Context) (types.Info, error) {
	return c.client.Info(ctx)
}

// RegistryLogin authenticates the docker server with a given docker registry.
func (c *OpClient) RegistryLogin(ctx context.Context, auth types.AuthConfig) error {
	body, err := c.client.RegistryLogin(ctx, auth)
	if err == nil {
		c.Time = time.Now()
		c.IdentityToken = body.IdentityToken
		c.Status = body.Status
	}
	return err
}
