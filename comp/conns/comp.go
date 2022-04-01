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
	Start: func(orch *core.Orchestrator, ctx context.Context, deps map[core.ComponentPath]core.ComponentReference) core.Component {
		c := &component{
			logger:        logger.Wrap(deps),
			users:         deps["comp/users.Main"],
			newConnection: make(chan net.Conn, 5),
			incoming:      make(chan incoming, 5),
			outgoing:      make(chan outgoing, 5),
			ctx:           ctx,
			done:          make(chan struct{}),
		}
		go c.run()
		return c
	},
}

type component struct {
	core.BaseComponent
	logger logger.Wrapper
	users  core.ComponentReference

	newConnection chan net.Conn
	incoming      chan incoming
	outgoing      chan outgoing

	ctx  context.Context
	done chan struct{}
}

var _ core.Component = &component{}

// NewReference implements core.Component#NewReference.
func (c *component) NewReference() core.ComponentReference {
	return &connReference{newConnection: c.newConnection}
}

// Done implements core.Component#Done.
func (c *component) Done() <-chan struct{} {
	return c.done
}

func (c *component) run() {
	defer close(c.done)
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
		case <-c.ctx.Done():
			// close all net.Conn's
			for _, c := range conns {
				c.conn.Close()
			}
			// wait for all to close
			for len(conns) > 0 {
				inc := <-c.incoming
				if inc.close {
					c.logger.Output(fmt.Sprintf("Got close from %d", inc.cid))
					c.users.RequestAsync(context.Background(), users.UserGone{Cid: inc.cid})
					delete(conns, inc.cid)
					if len(conns) == 0 {
						break
					}
				}
			}
			return
		}
	}
}
