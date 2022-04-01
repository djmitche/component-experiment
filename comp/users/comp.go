package users

import (
	"comps/comp/logger"
	"comps/core"
	"context"
	"fmt"
	"strings"
	"sync"
)

var componentPath core.ComponentPath = "comp/users.Main"

// Main is the component implementation for this package (`comp/users.Main`).
//
// This component accepts NewUser, UserGone, and UserMessage messages to handle
// user traffic.  It accepts a SetConnsComponent message at startup to identify
// the `comp/conns.Main` component, on which it has a weak dependency.
var Main = core.ComponentImpl{
	Path:         componentPath,
	Dependencies: []core.ComponentPath{"comp/logger.Main"},
	Start: func(deps map[core.ComponentPath]core.ComponentReference) core.Component {
		c := &component{
			logger: logger.Wrap(deps),
			users:  map[int]*user{},
		}
		return c
	},
}

type component struct {
	core.BaseComponent
	mu     sync.Mutex
	logger logger.Wrapper
	users  map[int]*user
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
		u := &user{cid: v.Cid, sendMessage: v.SendMessage}
		c.users[v.Cid] = u
		u.sendMessage("welcome!")
	case UserGone:
		delete(c.users, v.Cid)
	case UserMessage:
		u, found := c.users[v.Cid]
		if !found {
			break
		}
		switch {
		case strings.HasPrefix(v.Message, "/join"):
			room := v.Message[6:]
			if u.room != "" {
				c.sendToRoom(0, u.room, fmt.Sprintf("%d has left %s", v.Cid, room))
			}
			u.room = room
			c.sendToRoom(0, u.room, fmt.Sprintf("%d has joined %s", v.Cid, room))
			c.logger.Output(fmt.Sprintf("%d has joined %s", v.Cid, room))
		default:
			if u.room == "" {
				u.sendMessage("join a room first (/join)")
			} else {
				c.sendToRoom(v.Cid, u.room, fmt.Sprintf("%d: %s", v.Cid, v.Message))
			}
		}
	default:
		return nil, fmt.Errorf("Unrecognized message type %T", msg)
	}

	return nil, nil
}

func (c *component) sendToRoom(senderCid int, room string, message string) {
	for cid, u := range c.users {
		if cid != senderCid && u.room == room {
			u.sendMessage(message)
		}
	}
}

// RequestAsync implements core.ComponentReference#RequestAsync.
func (c *component) RequestAsync(ctx context.Context, msg core.Message) {
	c.Request(ctx, msg)
}
