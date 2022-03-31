package main

import (
	"comps/comp/conns"
	"comps/comp/listen"
	"comps/comp/logger"
	"comps/comp/users"
	"comps/core"
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	orch := core.NewOrchestrator(logger.Main, listen.Main, conns.Main, users.Main)
	l, err := orch.Start("comp/listen.Main")
	if err != nil {
		fmt.Printf("Uhoh: %s\n", err)
		return
	}
	l.Request(ctx, listen.Run{})
}
