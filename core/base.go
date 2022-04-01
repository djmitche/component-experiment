package core

import (
	"context"
	"fmt"
)

// BaseComponent is an empty type that can be embedded in a component to
// implement the Component interface.  By default, it returns a
// BaseComponentReference (itself, in fact) that accepts no requests.
type BaseComponent struct{ BaseComponentReference }

// NewReference implements Component#NewReference.
func (bc *BaseComponent) NewReference() ComponentReference {
	return bc
}

// BaseComponentReference is an empty type implementing ComponentReference and
// failing on all requests.  It is a useful shortcut for components which do not
// accept incoming messages.
type BaseComponentReference struct{}

// Request implements ComponentReference#Request.
func (bcr *BaseComponentReference) Request(context.Context, Message) (Message, error) {
	return nil, fmt.Errorf("%T does not accept requests", bcr)
}

// Request implements ComponentReference#RequestAsync.
func (bcr *BaseComponentReference) RequestAsync(context.Context, Message) {
}
