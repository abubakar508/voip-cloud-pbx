package natsclient

import (
	"encoding/json"
	"fmt"

	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/config"
	"github.com/nats-io/nats.go"
)

type Client struct {
	nc *nats.Conn
}

func New(cfg *config.AppConfig) (*Client, error) {
	// Expect NATS_URL in config, fallback to default
	url := cfg.NatsURL
	if url == "" {
		url = "nats://nats:4222"
	}
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &Client{nc: nc}, nil
}

func (c *Client) Close() {
	if c.nc != nil {
		c.nc.Close()
	}
}

func (c *Client) PublishJSON(subject string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.nc.Publish(subject, data)
}

func (c *Client) Flush() error {
	return c.nc.Flush()
}

func (c *Client) URL() string {
	if c.nc == nil {
		return ""
	}
	return fmt.Sprintf("%s", c.nc.ConnectedUrl())
}
