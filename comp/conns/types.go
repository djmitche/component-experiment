package conns

import (
	"net"
)

// Output is a core.Message containing the string to be logged.
type Connection struct {
	Conn net.Conn
}

type userLine struct {
	uid  int
	line string
}
