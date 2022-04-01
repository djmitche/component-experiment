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
	Start: func(orch *core.Orchestrator, ctx context.Context, deps map[core.ComponentPath]core.ComponentReference) core.Component {
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
	for comp, status := range o.orch.Status() {
		fmt.Fprintf(w, "%s: %s\n", string(comp), status.State)
		fmt.Fprintf(w, "  Depends on:\n")
		for _, d := range status.Dependencies {
			fmt.Fprintf(w, "    %s\n", string(d))
		}
	}
}
