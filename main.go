package main

import (
	"comps/comp/conns"
	"comps/comp/listen"
	"comps/comp/logger"
	"comps/comp/users"
	"comps/core"
	"comps/core/comp/debug"
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	orch := core.NewOrchestrator(
		Main,
		logger.Main,
		listen.Main,
		conns.Main,
		users.Main,
		debug.Main,
		debug.Expvar,
		debug.Orchestrator,
	)
	_, err := orch.Start("Main")
	if err != nil {
		fmt.Printf("Uhoh: %s\n", err)
		os.Exit(1)
	}

	// TODO: shutdown support
	for {
		time.Sleep(60 * time.Second)
	}
}

var componentPath core.ComponentPath = "Main"

var Main = core.ComponentImpl{
	Path: componentPath,
	Dependencies: []core.ComponentPath{
		"comp/logger.Main",
		"comp/listen.Main",
		"core/comp/debug.Main",
		"core/comp/debug.Expvar",
		"core/comp/debug.Orchestrator",
	},
	Start: func(orch *core.Orchestrator, deps map[core.ComponentPath]core.ComponentReference) core.Component {
		ctx := context.Background()
		deps["comp/logger.Main"].RequestAsync(ctx, logger.Output{Message: "Debug on http://127.0.0.1:8080"})
		deps["core/comp/debug.Main"].RequestAsync(ctx, debug.Serve{Port: 8080})
		deps["comp/listen.Main"].RequestAsync(ctx, listen.Run{})
		return &comp{}
	},
}

type comp struct{}

var _ core.Component = &comp{}
var _ core.ComponentReference = &comp{}

// NewReference implements core.Component#NewReference.
func (l *comp) NewReference() core.ComponentReference {
	return l
}

// Request implements core.ComponentReference#Request.
func (l *comp) Request(ctx context.Context, msg core.Message) (core.Message, error) {
	switch msg.(type) {
	default:
		return nil, fmt.Errorf("Unrecognized message type %T", msg)
	}
}

// RequestAsync implements core.ComponentReference#RequestAsync.
func (l *comp) RequestAsync(ctx context.Context, msg core.Message) {
	go l.Request(ctx, msg)
}
