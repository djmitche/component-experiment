package core

import (
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
	active map[ComponentPath]Component

	// Root is a ComponentReference to the root component (set after Start)
	Root ComponentReference

	// RootPath is the path of the root component (the first argument to NewOrchestrator)
	RootPath ComponentPath
}

// NewOrchestrator creates a new orchestrator, containing the given component
// implementations.  The first argument is considered the "root" component, and
// only components on which it depends (directly or indirectly) will be
// instantiated.
func NewOrchestrator(componentImpls ...ComponentImpl) *Orchestrator {
	orch := &Orchestrator{
		registered: make(map[ComponentPath]ComponentImpl),
		active:     make(map[ComponentPath]Component),
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

// DependencyGraph returns the dependency graph, in the form of a map from node
// name to its dependencies, each given by its ComponentPath.
func (orch *Orchestrator) DependencyGraph() map[ComponentPath][]ComponentPath {
	orch.mu.Lock()
	defer orch.mu.Unlock()

	rv := map[ComponentPath][]ComponentPath{}
	for path, _ := range orch.active {
		deps := orch.registered[path].Dependencies
		rv[path] = deps
	}
	return rv
}

// getComponentReference loads the given component, if it is not already loaded, and returns
// a reference to it.  This assumes that orch.mu is held.
func (orch *Orchestrator) getComponentReference(path ComponentPath) (ComponentReference, error) {
	var recur func(seen []ComponentPath, path ComponentPath) (ComponentReference, error)
	recur = func(seen []ComponentPath, path ComponentPath) (ComponentReference, error) {
		comp := orch.active[path]
		if comp == nil {
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

			comp = compImpl.Start(orch, deps)
			orch.active[path] = comp
		}

		return comp.NewReference(), nil
	}

	return recur([]ComponentPath{}, path)
}
