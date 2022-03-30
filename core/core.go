package core

import (
	"context"
	"fmt"
	"sync"
)

// ComponentPath identifies a component.
//
// Format is <package path>.Component, e.g., `pkg/telemetry.Main`
type ComponentPath string

// ComponentImpl defines a component implementation.  These are simple (usually
// empty) objects defining methods to start components.
type ComponentImpl interface {
	// Path returns the component path for this component implementation. This value
	// should not change.
	Path() ComponentPath

	// Dependencies returns the component paths on which this component relies.  This
	// value should not change.
	Dependencies() []ComponentPath

	// Start starts an instance of the component.  This will be called on-demand, when
	// the component is needed.  The given map will contain an entry for every dependency
	// path given by Dependencies().
	Start(map[ComponentPath]ComponentReference) Component
}

// Component represents a running instance of a component implementation.
type Component interface {
	NewReference() ComponentReference
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

// Orchestrator orchestrates multiple components.
type Orchestrator struct {
	mu         sync.Mutex
	registered map[ComponentPath]ComponentImpl
	active     map[ComponentPath]Component
}

// NewOrchestrator creates a new orchestrator, containing the given component implementations.
func NewOrchestrator(componentImpls ...ComponentImpl) *Orchestrator {
	orch := &Orchestrator{
		registered: make(map[ComponentPath]ComponentImpl),
		active:     make(map[ComponentPath]Component),
	}
	for _, ci := range componentImpls {
		orch.registered[ci.Path()] = ci
	}
	return orch
}

// Start starts an orchestrator by starting the component with the named path
// (and all of its dependencies) and returning a ComponentReference.  Typically
// the next step is to call `compRef.Request(componentpkg.StartMessage{..})` to
// pass information to the component and cause it to start.
func (orch *Orchestrator) Start(path ComponentPath) (ComponentReference, error) {
	return orch.getComponentReference(path)
}

func (orch *Orchestrator) getComponentReference(path ComponentPath) (ComponentReference, error) {
	orch.mu.Lock()
	defer orch.mu.Unlock()

	var recur func(seen []ComponentPath, path ComponentPath) (ComponentReference, error)
	recur = func(seen []ComponentPath, path ComponentPath) (ComponentReference, error) {
		comp := orch.active[path]
		if comp == nil {
			compImpl := orch.registered[path]
			if compImpl == nil {
				err := fmt.Errorf("No component with path %s", path)
				return nil, err
			}

			seen = append(seen, path)

			deps := map[ComponentPath]ComponentReference{}
			for _, depPath := range compImpl.Dependencies() {
				ref, err := recur(seen, depPath)
				if err != nil {
					return nil, err
				}
				deps[depPath] = ref
			}

			comp = compImpl.Start(deps)
			orch.active[path] = comp
		}

		return comp.NewReference(), nil
	}

	return recur([]ComponentPath{}, path)
}

// TODO: support partial shutdown/restart -- shutdown all dependent components first
