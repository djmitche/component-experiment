package logger

import (
	"comps/core"
	"context"
	"fmt"
)

type loggerImpl struct{}

// Main is the component implementation for this package (`comp/logger.Main`).
//
// On requests with messages of type `comp/logger.Output`, it logs the message
// and returns nil.
var Main core.ComponentImpl = &loggerImpl{}

// Path implements core.ComponentImpl#Path.
func (*loggerImpl) Path() core.ComponentPath {
	return (core.ComponentPath)("comp/logger.Main")
}

// Dependencies implements core.ComponentImpl#Dependencies.
func (*loggerImpl) Dependencies() []core.ComponentPath {
	return []core.ComponentPath{}
}

// Start implements core.ComponentImpl#Start.
func (*loggerImpl) Start(map[core.ComponentPath]core.ComponentReference) core.Component {
	return &logger{}
}

// Output is a core.Message containing the string to be logged.
type Output struct {
	Message string
}

type logger struct{}

var _ core.Component = &logger{}

// NewReference implements core.Component#NewReference.
func (l *logger) NewReference() core.ComponentReference {
	return &reference{}
}

type reference struct{}

var _ core.ComponentReference = &reference{}

// Request implements core.ComponentReference#Request.
func (l *reference) Request(ctx context.Context, msg core.Message) (core.Message, error) {
	switch v := msg.(type) {
	case Output:
		fmt.Println(v.Message)
		return nil, nil
	default:
		return nil, fmt.Errorf("Unrecognized message type %T", msg)
	}
}

// RequestAsync implements core.ComponentReference#RequestAsync.
func (l *reference) RequestAsync(ctx context.Context, msg core.Message) {
	l.Request(ctx, msg)
}
