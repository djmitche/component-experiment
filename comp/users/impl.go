package users

import (
	"comps/comp/logger"
	"comps/core"
)

var componentName core.ComponentPath = "comp/users.Main"

type impl struct{}

// Main is the component implementation for this package (`comp/users.Main`).
//
// This component accepts NewUser, UserGone, and UserMessage messages to handle
// user traffic.  It accepts a SetConnsComponent message at startup to identify
// the `comp/conns.Main` component, on which it has a weak dependency.
var Main core.ComponentImpl = &impl{}

// Path implements core.ComponentImpl#Path.
func (*impl) Path() core.ComponentPath {
	return (core.ComponentPath)(componentName)
}

// Dependencies implements core.ComponentImpl#Dependencies.
func (*impl) Dependencies() []core.ComponentPath {
	return []core.ComponentPath{"comp/logger.Main"}
}

// Start implements core.ComponentImpl#Start.
func (*impl) Start(deps map[core.ComponentPath]core.ComponentReference) core.Component {
	c := &component{
		logger: logger.Wrap(deps),
		users:  map[int]user{},
	}
	return c
}
