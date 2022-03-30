package logger

import "comps/core"

var componentPath core.ComponentPath = "comp/logger.Main"

type loggerImpl struct{}

// Main is the component implementation for this package (`comp/logger.Main`).
//
// On requests with messages of type `comp/logger.Output`, it logs the message
// and returns nil.
var Main core.ComponentImpl = &loggerImpl{}

// Path implements core.ComponentImpl#Path.
func (*loggerImpl) Path() core.ComponentPath {
	return (core.ComponentPath)(componentPath)
}

// Dependencies implements core.ComponentImpl#Dependencies.
func (*loggerImpl) Dependencies() []core.ComponentPath {
	return []core.ComponentPath{}
}

// Start implements core.ComponentImpl#Start.
func (*loggerImpl) Start(map[core.ComponentPath]core.ComponentReference) core.Component {
	return &logger{}
}
