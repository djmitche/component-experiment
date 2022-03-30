package listen

import (
	"comps/comp/conn"
	"comps/comp/logger"
	"comps/core"
	"context"
	"fmt"
	"net"
)

type listenImpl struct{}

// Main is the component implementation for this package (`comp/listen.Main`).
//
// On requests with messages of type `comp/listen.Run`, it listens for new connections
// and hands them to the `comp/conn.Main` component.
var Main core.ComponentImpl = &listenImpl{}

// Path implements core.ComponentImpl#Path.
func (*listenImpl) Path() core.ComponentPath {
	return (core.ComponentPath)("comp/listen.Main")
}

// Dependencies implements core.ComponentImpl#Dependencies.
func (*listenImpl) Dependencies() []core.ComponentPath {
	return []core.ComponentPath{"comp/logger.Main", "comp/conn.Main"}
}

// Start implements core.ComponentImpl#Start.
func (*listenImpl) Start(deps map[core.ComponentPath]core.ComponentReference) core.Component {
	l := &listen{
		logger: deps["comp/logger.Main"],
		conn:   deps["comp/conn.Main"],
	}
	return l
}

// Run is a core.Message that indicates the component should run
type Run struct{}

type listen struct {
	logger core.ComponentReference
	conn   core.ComponentReference
}

var _ core.Component = &listen{}
var _ core.ComponentReference = &listen{}

// NewReference implements core.Component#NewReference.
func (l *listen) NewReference() core.ComponentReference {
	return l
}

// Request implements core.ComponentReference#Request.
func (l *listen) Request(ctx context.Context, msg core.Message) (core.Message, error) {
	switch msg.(type) {
	case Run:
		err := l.run(ctx)
		return nil, err
	default:
		return nil, fmt.Errorf("Unrecognized message type %T", msg)
	}
}

// RequestAsync implements core.ComponentReference#RequestAsync.
func (l *listen) RequestAsync(ctx context.Context, msg core.Message) {
	go l.Request(ctx, msg)
}

func (l *listen) run(ctx context.Context) error {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9000")
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	l.logger.Request(ctx, logger.Output{Message: fmt.Sprintf("Listening on port %d", 9000)})

	// stupid workaround to stop listening when the context expires
	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		c, err := listener.Accept()
		if err != nil {
			return err
		}
		_, err = l.conn.Request(ctx, conn.Connection{Conn: c})
		if err != nil {
			return err
		}
	}

	l.logger.Request(ctx, logger.Output{Message: fmt.Sprintf("Done listening on port %d", 9000)})
	return nil
}
