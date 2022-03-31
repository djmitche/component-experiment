package conns

import (
	"bufio"
	"net"
)

type connection struct {
	conn     net.Conn
	cid      int
	outgoing chan string
}

func (c *connection) run(incomingChan chan<- incoming) {
	// send outgoing messages to the remote end until error
	go func() {
		for {
			msg := append([]byte(<-c.outgoing), '\n')
			for len(msg) > 0 {
				n, err := c.conn.Write(msg)
				if err != nil {
					return // assume this is EOF
				}
				msg = msg[n:]
			}
		}
	}()

	// read from the connection and send to incoming, until EOF
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		incomingChan <- incoming{
			cid:  c.cid,
			line: scanner.Text(),
		}
	}

	// close the conn, for good measure
	_ = c.conn.Close()

	incomingChan <- incoming{
		cid:   c.cid,
		close: true,
	}
}
