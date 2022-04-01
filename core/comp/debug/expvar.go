package debug

import (
	"comps/core"
	"context"
	expvarPkg "expvar"
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

type expvar struct{ core.BaseComponent }
