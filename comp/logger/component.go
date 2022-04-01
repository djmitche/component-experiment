package logger

import (
	"comps/core"
	"context"
	"fmt"
)

var componentPath core.ComponentPath = "comp/logger.Main"

// Main is the component implementation for this package (`comp/logger.Main`).
//
// On requests with messages of type `comp/logger.Output`, it logs the message
// and returns nil.
var Main = core.ComponentImpl{
	Path:         componentPath,
	Dependencies: []core.ComponentPath{},
	Start: func(map[core.ComponentPath]core.ComponentReference) core.Component {
		return &logger{}
	},
}

type logger struct {
	core.BaseComponent
}

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
