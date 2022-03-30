package logger

import (
	"comps/core"
	"context"
	"fmt"
)

type logger struct{}

var _ core.Component = &logger{}
var _ core.ComponentReference = &logger{}

// NewReference implements core.Component#NewReference.
func (l *logger) NewReference() core.ComponentReference {
	return l
}

// Request implements core.ComponentReference#Request.
func (l *logger) Request(ctx context.Context, msg core.Message) (core.Message, error) {
	switch v := msg.(type) {
	case Output:
		fmt.Println(v.Message)
		return nil, nil
	default:
		return nil, fmt.Errorf("Unrecognized message type %T", msg)
	}
}

// RequestAsync implements core.ComponentReference#RequestAsync.
func (l *logger) RequestAsync(ctx context.Context, msg core.Message) {
	l.Request(ctx, msg)
}
