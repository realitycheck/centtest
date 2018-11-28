package centtest

import (
	"github.com/centrifugal/centrifuge-go"
)

type Client struct {
	id        string
	cli       *centrifuge.Client
	connected bool
}

func NewClient(wsURL string, config centrifuge.Config) *Client {
	c := &Client{
		cli: centrifuge.New(wsURL, config),
	}
	return c
}

func (c *Client) String() string {
	if c.connected {
		return c.id
	}
	return "<not connected>"
}

func (c *Client) disconnect() error {
	if err := c.cli.Disconnect(); err != nil {
		return err
	}
	if err := c.cli.Close(); err != nil {
		return err
	}
	return nil
}
