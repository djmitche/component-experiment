package conns

import (
	"comps/core"
	"context"
	"fmt"
	"net"
)

type connReference struct {
	newConnection chan<- net.Conn
}

var _ core.ComponentReference = &connReference{}

// Request implements core.ComponentReference#Request.
func (c *connReference) Request(ctx context.Context, msg core.Message) (core.Message, error) {
	switch v := msg.(type) {
	case Connection:
		select {
		case c.newConnection <- v.Conn:
			return nil, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	default:
		return nil, fmt.Errorf("Unrecognized message type %T", msg)
	}
}

// RequestAsync implements core.ComponentReference#RequestAsync.
func (c *connReference) RequestAsync(ctx context.Context, msg core.Message) {
	c.Request(ctx, msg)
}
