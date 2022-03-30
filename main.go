package main

import (
	"comps/comp/echo"
	"comps/comp/listen"
	"comps/comp/logger"
	"comps/core"
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	orch := core.NewOrchestrator(logger.Main, listen.Main, echo.Main)
	l, err := orch.Start("comp/listen.Main")
	if err != nil {
		fmt.Printf("Uhoh: %s\n", err)
	}
	l.Request(ctx, listen.Run{})
}
