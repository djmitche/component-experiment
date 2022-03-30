package echo

import (
	"comps/comp/logger"
	"comps/core"
	"context"
	"fmt"
	"io"
	"net"
)

type echoImpl struct{}

// Main is the component implementation for this package (`comp/echo.Main`).
//
// The component accepts messages of type Connection, indicating to begin echoing on
// the given TCP connection.  The component responds immediately with a nil message,
// and echoes until EOF.
var Main core.ComponentImpl = &echoImpl{}

// Path implements core.ComponentImpl#Path.
func (*echoImpl) Path() core.ComponentPath {
	return (core.ComponentPath)("comp/echo.Main")
}

// Dependencies implements core.ComponentImpl#Dependencies.
func (*echoImpl) Dependencies() []core.ComponentPath {
	return []core.ComponentPath{"comp/logger.Main"}
}

// Start implements core.ComponentImpl#Start.
func (*echoImpl) Start(deps map[core.ComponentPath]core.ComponentReference) core.Component {
	return &echo{logger: deps["comp/logger.Main"]}
}

// Output is a core.Message containing the string to be logged.
type Connection struct {
	Conn net.Conn
}

type echo struct {
	logger core.ComponentReference
}

var _ core.Component = &echo{}
var _ core.ComponentReference = &echo{}

// NewReference implements core.Component#NewReference.
func (e *echo) NewReference() core.ComponentReference {
	return e
}

// Request implements core.ComponentReference#Request.
func (e *echo) Request(ctx context.Context, msg core.Message) (core.Message, error) {
	switch v := msg.(type) {
	case Connection:
		go func() {
			e.logger.Request(ctx, logger.Output{Message: fmt.Sprintf("Echoing for %#v", v.Conn)})
			_, _ = io.Copy(v.Conn, v.Conn)
			e.logger.Request(ctx, logger.Output{Message: fmt.Sprintf("Done echoing for %#v", v.Conn)})
		}()
		return nil, nil
	default:
		return nil, fmt.Errorf("Unrecognized message type %T", msg)
	}
}

// RequestAsync implements core.ComponentReference#RequestAsync.
func (e *echo) RequestAsync(ctx context.Context, msg core.Message) {
	e.Request(ctx, msg)
}
