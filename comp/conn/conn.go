package conn

import (
	"bufio"
	"comps/comp/logger"
	"comps/core"
	"context"
	"fmt"
	"net"
)

type connImpl struct{}

// Main is the component implementation for this package (`comp/conn.Main`).
//
// The component accepts messages of type Connection, indicating to begin handling
// the given TCP connection.  The component responds immediately with a nil message,
// and handles the connection until EOF.
var Main core.ComponentImpl = &connImpl{}

// Path implements core.ComponentImpl#Path.
func (*connImpl) Path() core.ComponentPath {
	return (core.ComponentPath)("comp/conn.Main")
}

// Dependencies implements core.ComponentImpl#Dependencies.
func (*connImpl) Dependencies() []core.ComponentPath {
	return []core.ComponentPath{"comp/logger.Main"}
}

// Start implements core.ComponentImpl#Start.
func (*connImpl) Start(deps map[core.ComponentPath]core.ComponentReference) core.Component {
	c := &conn{
		logger:        logger.Wrap(deps),
		newConnection: make(chan net.Conn, 5),
		userLines:     make(chan userLine, 5),
	}
	go c.run()
	return c
}

// Output is a core.Message containing the string to be logged.
type Connection struct {
	Conn net.Conn
}

type conn struct {
	logger logger.Wrapper

	newConnection chan net.Conn
	userLines     chan userLine
}

var _ core.Component = &conn{}

// NewReference implements core.Component#NewReference.
func (c *conn) NewReference() core.ComponentReference {
	return &connReference{newConnection: c.newConnection}
}

func (c *conn) run() {
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

type userLine struct {
	uid  int
	line string
}

type user struct {
	uid      int
	outgoing chan string
}

func (u *user) read(conn net.Conn, uid int, userLines chan<- userLine) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		userLines <- userLine{
			uid:  uid,
			line: scanner.Text(),
		}
	}
}

func (u *user) write(conn net.Conn, outgoing <-chan string) {
	for {
		msg := append([]byte(<-outgoing), '\n')
		for len(msg) > 0 {
			n, err := conn.Write(msg)
			if err != nil {
				panic(err) // TODO :)
			}
			msg = msg[n:]
		}
	}
}

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
