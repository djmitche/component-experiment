package conns

import (
	"comps/comp/logger"
	"comps/core"
	"net"
)

var componentName core.ComponentPath = "comp/conns.Main"

type impl struct{}

// Main is the component implementation for this package (`comp/conn.Main`).
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
	return []core.ComponentPath{"comp/logger.Main"}
}

// Start implements core.ComponentImpl#Start.
func (*impl) Start(deps map[core.ComponentPath]core.ComponentReference) core.Component {
	c := &component{
		logger:        logger.Wrap(deps),
		newConnection: make(chan net.Conn, 5),
		userLines:     make(chan userLine, 5),
	}
	go c.run()
	return c
}
