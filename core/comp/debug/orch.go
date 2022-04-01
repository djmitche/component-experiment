package debug

import (
	"comps/core"
	"context"
	"fmt"
	"net/http"
)

var Orchestrator = core.ComponentImpl{
	Path:         componentPath("Orchestrator"),
	Dependencies: []core.ComponentPath{"core/comp/debug.Main"},
	Start: func(orch *core.Orchestrator, deps map[core.ComponentPath]core.ComponentReference) core.Component {
		ctx := context.Background()
		o := &orchestrator{orch: orch}
		deps["core/comp/debug.Main"].RequestAsync(
			ctx,
			RegisterHandler{
				Name:    "Orchestrator",
				Pattern: "/orchestrator",
				Handler: http.HandlerFunc(o.handler),
			})
		return o
	},
}

type orchestrator struct {
	core.BaseComponent
	orch *core.Orchestrator
}

func (o *orchestrator) handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	for comp, deps := range o.orch.DependencyGraph() {
		fmt.Fprintf(w, "%s depends on:\n", string(comp))
		for _, d := range deps {
			fmt.Fprintf(w, "  %s\n", string(d))
		}
	}
}
