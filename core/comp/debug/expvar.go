package debug

import (
	"comps/core"
	"context"
	expvarPkg "expvar"
	"fmt"
)

// Expvar will include the Go expvar handler at /debug/vars in the
// `core/comp/debug.Main` component's http handler.
var Expvar = core.ComponentImpl{
	Path:         componentPath("Expvar"),
	Dependencies: []core.ComponentPath{"core/comp/debug.Main"},
	Start: func(deps map[core.ComponentPath]core.ComponentReference) core.Component {
		ctx := context.Background()
		deps["core/comp/debug.Main"].RequestAsync(
			ctx,
			RegisterHandler{
				Name:    "expvar",
				Pattern: "/debug/vars",
				Handler: expvarPkg.Handler(),
			})
		return &expvar{}
	},
}

type expvar struct{}

var _ core.Component = &expvar{}
var _ core.ComponentReference = &expvar{}

// NewReference implements core.Component#NewReference.
func (m *expvar) NewReference() core.ComponentReference {
	return m
}

// Request implements core.ComponentReference#Request.
func (m *expvar) Request(ctx context.Context, msg core.Message) (core.Message, error) {
	switch msg.(type) {
	default:
		return nil, fmt.Errorf("Unrecognized message type %T", msg)
	}
}

// RequestAsync implements core.ComponentReference#RequestAsync.
func (m *expvar) RequestAsync(ctx context.Context, msg core.Message) {
	m.Request(ctx, msg)
}
