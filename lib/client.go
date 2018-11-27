package centtest

import (
	"github.com/centrifugal/centrifuge-go"
)

type Client struct {
	id        string
	cli       *centrifuge.Client
	u         *User
	connected bool
}

func NewClient(wsURL string, config centrifuge.Config) *Client {
	c := &Client{
		cli: centrifuge.New(wsURL, config),
	}
	return c
}

func (c *Client) String() string {
	if c.id != "" {
		return c.id
	} else {
		return "<NOT CONNECTED>"
	}
}
