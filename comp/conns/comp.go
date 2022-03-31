package conns

import (
	"comps/comp/logger"
	"comps/comp/users"
	"comps/core"
	"context"
	"fmt"
	"net"
)

var componentPath core.ComponentPath = "comp/conns.Main"

// Main is the component implementation for this package (`comp/conns.Main`).
//
// The component accepts messages of type Connection, indicating to begin handling
// the given TCP connection.  The component responds immediately with a nil message,
// and handles the connection until EOF.
var Main = core.ComponentImpl{
	Path:         componentPath,
	Dependencies: []core.ComponentPath{"comp/logger.Main", "comp/users.Main"},
	Start: func(deps map[core.ComponentPath]core.ComponentReference) core.Component {
		c := &component{
			logger:        logger.Wrap(deps),
			users:         deps["comp/users.Main"],
			newConnection: make(chan net.Conn, 5),
			incoming:      make(chan incoming, 5),
			outgoing:      make(chan outgoing, 5),
		}
		go c.run()
		return c
	},
}

type component struct {
	logger logger.Wrapper
	users  core.ComponentReference

	newConnection chan net.Conn
	incoming      chan incoming
	outgoing      chan outgoing
}

var _ core.Component = &component{}

// NewReference implements core.Component#NewReference.
func (c *component) NewReference() core.ComponentReference {
	return &connReference{newConnection: c.newConnection}
}

func (c *component) run() {
	nextUser := 1
	conns := map[int]connection{}
	for {
		select {
		case netconn := <-c.newConnection:
			cid := nextUser
			nextUser++
			outgoingChan := make(chan string, 5)
			conn := connection{netconn, cid, outgoingChan}
			conns[cid] = conn
			go conn.run(c.incoming)
			c.users.RequestAsync(context.Background(),
				users.NewUser{
					Cid: cid,
					SendMessage: func(msg string) {
						c.outgoing <- outgoing{cid: cid, line: msg}
					},
				})

		case out := <-c.outgoing:
			conn, found := conns[out.cid]
			if found {
				conn.outgoing <- out.line
			}

		case inc := <-c.incoming:
			if inc.close {
				c.logger.Output(fmt.Sprintf("Got close from %d", inc.cid))
				c.users.RequestAsync(context.Background(), users.UserGone{Cid: inc.cid})
				delete(conns, inc.cid)
			} else {
				c.logger.Output(fmt.Sprintf("Got message %#v from %d", inc.line, inc.cid))
				c.users.RequestAsync(context.Background(), users.UserMessage{Cid: inc.cid, Message: inc.line})
			}
		}
	}
}
