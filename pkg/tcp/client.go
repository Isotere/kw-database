package tcp

import (
	"context"
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type HandlerClient func(context.Context, QueryClient)

type Client struct {
	log Logger

	host string
	port int
}

func NewClient(log Logger, host string, port int) *Client {
	return &Client{
		log:  log,
		host: host,
		port: port,
	}
}

func (c *Client) ProcessRequest(ctx context.Context, handler HandlerClient) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.host, c.port))
	if err != nil {
		return errors.Wrap(err, "Error processing request")
	}
	defer func() {
		_ = conn.Close()
	}()

	handler(ctx, NewQuery(conn))

	return nil
}
