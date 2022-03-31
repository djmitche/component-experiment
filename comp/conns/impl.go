package conns

import (
	"comps/comp/logger"
	"comps/core"
	"net"
)

var componentName core.ComponentPath = "comp/conns.Main"

type impl struct{}

// Main is the component implementation for this package (`comp/conns.Main`).
//
// The component accepts messages of type Connection, indicating to begin handling
// the given TCP connection.  The component responds immediately with a nil message,
// and handles the connection until EOF.
var Main core.ComponentImpl = &impl{}

// Path implements core.ComponentImpl#Path.
func (*impl) Path() core.ComponentPath {
	return (core.ComponentPath)(componentName)
}

// Dependencies implements core.ComponentImpl#Dependencies.
func (*impl) Dependencies() []core.ComponentPath {
	return []core.ComponentPath{"comp/logger.Main", "comp/users.Main"}
}

// Start implements core.ComponentImpl#Start.
func (*impl) Start(deps map[core.ComponentPath]core.ComponentReference) core.Component {
	c := &component{
		logger:        logger.Wrap(deps),
		users:         deps["comp/users.Main"],
		newConnection: make(chan net.Conn, 5),
		incoming:      make(chan incoming, 5),
		outgoing:      make(chan outgoing, 5),
	}
	go c.run()
	return c
}
