package users

import (
	"comps/comp/logger"
	"comps/core"
	"context"
	"fmt"
	"sync"
)

type component struct {
	mu     sync.Mutex
	logger logger.Wrapper
	users  map[int]user
}

var _ core.Component = &component{}
var _ core.ComponentReference = &component{}

// NewReference implements core.Component#NewReference.
func (c *component) NewReference() core.ComponentReference {
	return c
}

// Request implements core.ComponentReference#Request.
func (c *component) Request(ctx context.Context, msg core.Message) (core.Message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch v := msg.(type) {
	case NewUser:
		u := user{cid: v.Cid, sendMessage: v.SendMessage}
		c.users[v.Cid] = u
		u.sendMessage("welcome!")
	case UserGone:
		delete(c.users, v.Cid)
	case UserMessage:
		for cid, u := range c.users {
			message := fmt.Sprintf("%d: %s", v.Cid, v.Message)
			if cid != v.Cid {
				u.sendMessage(message)
			}
		}
	default:
		return nil, fmt.Errorf("Unrecognized message type %T", msg)
	}

	return nil, nil
}

// RequestAsync implements core.ComponentReference#RequestAsync.
func (c *component) RequestAsync(ctx context.Context, msg core.Message) {
	c.Request(ctx, msg)
}
