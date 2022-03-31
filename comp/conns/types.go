package conns

import (
	"net"
)

// Connection is a core.Message containing a new connection.
type Connection struct {
	Conn net.Conn
}

// incoming represents an incoming event from a connection
type incoming struct {
	// cid is the connection ID for this connection
	cid int

	// if close is true, the connection is closed
	close bool

	// if close is false, line contains the line from the connection (without newline)
	line string
}

// outgoing represents an outgoing event to a connection
type outgoing struct {
	// cid is the connection ID for this connection
	cid int

	// the line to sendt o the connection
	line string
}
