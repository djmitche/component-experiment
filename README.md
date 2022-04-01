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

## Why References?

ComponentReference is separate from the Component to allow remote references (by proxying requests an responses) and to support the common case where the Component is an actor and reqests/responses are handled with channels.

In cases where this is not required, the BaseComponent and BaseComponentRef types remove most of the boilerplate.

## Wrapped References

Calling a method on a ComponentReference isn't very ergonomic.
`comp/logger.Main` shows an alternative, with a type that wraps a ComponentReference and provides a more ergonomic interface.

## Circular Dependencies

The comps/conns.Main and comps/users.Main components have a circular dependency: comps/conns.Main must provide incoming messages to comps/users.Main, while comps/users.Main must provide outgoing messages to comps/conns.Main.
This is accomplished by replacing one of those dependencies with a simple callback.

## Debug Output

The core/comps/debug.* components provide a debug server containing useful debugging information about the running system.
This follows the Go expvar pattern, but has pluggable handlers and can include lots of other useful output, defined by other components.
They need only depend on `core/comps/debug.Main` and send it a `RegisterHandler` message.

## Shutdown

Shutodwn occurs in the opposite order of startup.
Components are given a context which will cancel when they should stop -- this makes it easy to pass that context to other operations in the component.
Components signal that they are complete with a Done() method similar to that in the context package.
The Base types make all of this invisible to components that do not have any need to do anything special when shutting down (such as comps/logging.Main).

# TODO

 - health monitoring
 - api
 - telemetry
