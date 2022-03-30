package conns

import (
	"comps/comp/logger"
	"comps/core"
	"fmt"
	"net"
)

type component struct {
	logger logger.Wrapper

	newConnection chan net.Conn
	userLines     chan userLine
}

var _ core.Component = &component{}

// NewReference implements core.Component#NewReference.
func (c *component) NewReference() core.ComponentReference {
	return &connReference{newConnection: c.newConnection}
}

func (c *component) run() {
	nextUser := 1
	users := map[int]user{}
	for {
		select {
		case conn := <-c.newConnection:
			uid := nextUser
			nextUser++
			outgoing := make(chan string, 5)
			user := user{uid, outgoing}
			users[uid] = user
			go user.read(conn, uid, c.userLines)
			go user.write(conn, outgoing)

		case userLine := <-c.userLines:
			c.logger.Output(fmt.Sprintf("Got message %#v from %d", userLine.line, userLine.uid))
			output := fmt.Sprintf("%d: %s", userLine.uid, userLine.line)
			for uid, u := range users {
				if uid != userLine.uid {
					u.outgoing <- output
				}
			}
		}
	}
}
