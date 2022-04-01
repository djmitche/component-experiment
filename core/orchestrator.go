package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Orchestrator orchestrates multiple components.
type Orchestrator struct {
	mu sync.Mutex

	// registered contains all registered components (passed to the constructor)
	registered map[ComponentPath]ComponentImpl

	// active contains all active components (those in the dependency graph of Root)
	active map[ComponentPath]activeComponent

	// Root is a ComponentReference to the root component (set after Start)
	Root ComponentReference

	// RootPath is the path of the root component (the first argument to NewOrchestrator)
	RootPath ComponentPath
}

type activeComponent struct {
	stop  context.CancelFunc
	comp  Component
	state ComponentState
}

// NewOrchestrator creates a new orchestrator, containing the given component
// implementations.  The first argument is considered the "root" component, and
// only components on which it depends (directly or indirectly) will be
// instantiated.
func NewOrchestrator(componentImpls ...ComponentImpl) *Orchestrator {
	orch := &Orchestrator{
		registered: make(map[ComponentPath]ComponentImpl),
		active:     make(map[ComponentPath]activeComponent),
	}
	for _, ci := range componentImpls {
		orch.registered[ci.Path] = ci
	}
	orch.RootPath = componentImpls[0].Path
	return orch
}

// Start starts an orchestrator by starting the component with the named path
// (and all of its dependencies) and returning a ComponentReference.  Typically
// the next step is to call `compRef.Request(componentpkg.StartMessage{..})` to
// pass information to the component and cause it to start.
func (orch *Orchestrator) Start() error {
	orch.mu.Lock()
	defer orch.mu.Unlock()

	if orch.Root != nil {
		return errors.New("Orchestrator has already been started")
	}

	root, err := orch.getComponentReference(orch.RootPath)
	orch.Root = root
	return err
}

// Stop stops a running orchestrator, in an orderly fashion.  Components are
// stopped only after everything depending on them are stopped.  This method will
// block until all components are stopped, or the passed context expires.
func (orch *Orchestrator) Stop(stopCtx context.Context) error {
	// components are stopped in the reverse of the order in which they were started
	order := make([]ComponentPath, len(orch.active))
	i := len(orch.active) - 1
	seen := map[ComponentPath]struct{}{}
	var recur func(path ComponentPath)
	recur = func(path ComponentPath) {
		_, found := seen[path]
		if !found {
			seen[path] = struct{}{}
			for _, dep := range orch.registered[path].Dependencies {
				recur(dep)
			}
			order[i] = path
			i--
		}
	}
	recur(orch.RootPath)
	if i != -1 {
		panic("not zero")
	}

	for _, path := range order {
		acomp := orch.active[path]
		acomp.stop()
		select {
		case <-acomp.comp.Done():
		case <-stopCtx.Done():
			return stopCtx.Err()
		}
	}
	return nil
}

// Status returns the status of the orchestrator, in the form of a map from
// component path to information about that component.
func (orch *Orchestrator) Status() map[ComponentPath]ComponentStatus {
	orch.mu.Lock()
	defer orch.mu.Unlock()

	rv := map[ComponentPath]ComponentStatus{}
	for path, acomp := range orch.active {
		deps := orch.registered[path].Dependencies
		rv[path] = ComponentStatus{
			Dependencies: deps,
			State:        acomp.state,
		}
	}
	return rv
}

// getComponentReference loads the given component, if it is not already loaded, and returns
// a reference to it.  This assumes that orch.mu is held.
func (orch *Orchestrator) getComponentReference(path ComponentPath) (ComponentReference, error) {
	bkgnd := context.Background()

	var recur func(seen []ComponentPath, path ComponentPath) (ComponentReference, error)
	recur = func(seen []ComponentPath, path ComponentPath) (ComponentReference, error) {
		acomp, found := orch.active[path]
		if !found {
			compImpl := orch.registered[path]
			if compImpl.Path == "" {
				err := fmt.Errorf("No component with path %s", path)
				return nil, err
			}

			seen = append(seen, path)

			deps := map[ComponentPath]ComponentReference{}
			for _, depPath := range compImpl.Dependencies {
				ref, err := recur(seen, depPath)
				if err != nil {
					return nil, err
				}
				deps[depPath] = ref
			}

			ctx, stop := context.WithCancel(bkgnd)
			acomp = activeComponent{
				comp:  compImpl.Start(orch, ctx, deps),
				stop:  stop,
				state: RunningState,
			}
			orch.active[path] = acomp
		}

		return acomp.comp.NewReference(), nil
	}

	return recur([]ComponentPath{}, path)
}
