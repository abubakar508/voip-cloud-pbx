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

func (c *Client) URL() string {
	if c.nc == nil {
		return ""
	}
	return fmt.Sprintf("%s", c.nc.ConnectedUrl())
}

func (c *Client) Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	return c.nc.Subscribe(subject, handler)
}

func (c *Client) Unmarshal(msg *nats.Msg, v interface{}) error {
	return json.Unmarshal(msg.Data, v)
}
