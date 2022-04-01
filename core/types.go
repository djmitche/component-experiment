package core

import "context"

// ComponentPath identifies a component.
//
// Format is <package path>.Component, e.g., `pkg/telemetry.Main`
type ComponentPath string

// ComponentState gives a component's state.  It is one of the *State constants.
type ComponentState string

// ComponentState values
const (
	// RunningState identifies a component that is running
	RunningState ComponentState = "running"

	// StoppingState identifies a component that is stopping
	StoppingState ComponentState = "stopping"

	// Stopped identifies a component that is stopped
	StoppedState ComponentState = "stopped"
)

// ComponentImpl defines a component implementation.  These are simple (usually
// empty) objects defining methods to start components.
type ComponentImpl struct {
	// Path is the component path for this component implementation.
	Path ComponentPath

	// Dependencies gives the component paths on which this component relies.
	Dependencies []ComponentPath

	// Start starts an instance of the component.  This will be called on-demand, when
	// the component is needed.
	//
	// When the Context is finishes, the component should finish any ongoing
	// activity and return a channel from its Done method that is closed when
	// the component has fully stopped.
	//
	// The `deps` map will contain an entry for every dependency path given by
	// Dependencies.
	Start func(*Orchestrator, context.Context, map[ComponentPath]ComponentReference) Component
}

// Component represents a running instance of a component implementation.
type Component interface {
	NewReference() ComponentReference

	// Done returns a channel which closes when this component has fully stopped.
	Done() <-chan struct{}
}

// Message defines types that can be used as requests or responses between components.
// Messages are recognized by casting to concrete types.
type Message interface{}

// ComponentReference is a reference to another component, used to communicate with that
// component.  The set of Message types supported by a component are defined by that
// component.
type ComponentReference interface {
	// Send a message to the given component and wait for a response message.  This method
	// may block the caller, but should return if the context is cancelled.
	Request(context.Context, Message) (Message, error)

	// Send a message to the given component and do not wait for a response.  This method
	// may block the caller until the message is enqueued, but does not wait until the
	// message is completely handled. It should return if the context is cancelled.
	RequestAsync(context.Context, Message)
}

// ComponentStatus is part of the return from the Orchestrator#Status method.
type ComponentStatus struct {
	// Dependencies gives the component's dependencies.
	Dependencies []ComponentPath

	// State gives the component's current state.
	State ComponentState
}
