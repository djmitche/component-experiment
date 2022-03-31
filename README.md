# Component Experiment

This is an experiment in building a well-isolated component architecture in Go.

There are three primary types involved:

 * ComponentImpl -- this defines a component implementation.
   It has a single name (ComponentPath) and is instantiated on demand.
   This is similar to a class in OOP.

 * Component -- this defines an instance of a component, created when another component depends on it.
   The system currently instantiates a component only once, but other strategies might be applied.
   This is similar to an instance of a class in OOP.
 
 * ComponentReference -- this defines a reference to a (possibly remote) component, mediating communication with that component.
   This is similar to an instance pointer in OOP.

# Interesting bits

## Why References

ComponentReference is separate from the Component to allow remote references (by proxying requests an responses) and to support the common case where the Component is an actor and reqests/responses are handled with channels.

## Wrapped References

Calling a method on a ComponentReference isn't very ergonomic.
`comp/logger.Main` shows an alternative, with a type that wraps a ComponentReference and provides a more ergonomic interface.

## Circular Dependencies

The comps/conns.Main and comps/users.Main components have a circular dependency: comps/conns.Main must provide incoming messages to comps/users.Main, while comps/users.Main must provide outgoing messages to comps/conns.Main.
This is accomplished by replacing one of those dependencies with a simple callback.

# TODO

 - health monitoring
 - shutdown
 - status reporting
 - api
 - telemetry
 - describe components and links (explicit "weak" deps?)
