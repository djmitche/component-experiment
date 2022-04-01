package core

import (
	"fmt"
	"sync"
)

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
		orch.registered[ci.Path] = ci
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

			comp = compImpl.Start(deps)
			orch.active[path] = comp
		}

		return comp.NewReference(), nil
	}

	return recur([]ComponentPath{}, path)
}
