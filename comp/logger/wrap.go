package logger

import (
	"comps/core"
	"context"
)

type Wrapper struct {
	wrapped core.ComponentReference
}

func (w *Wrapper) Output(message string) {
	w.wrapped.Request(context.Background(), Output{Message: message})
}

func Wrap(deps map[core.ComponentPath]core.ComponentReference) Wrapper {
	return Wrapper{wrapped: deps[componentPath]}
}
