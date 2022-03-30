package conns

import (
	"bufio"
	"net"
)

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
