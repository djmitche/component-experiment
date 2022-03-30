# Component Experiment

This is an experiment in building a well-isolated component architecture in Go.

There are three primary types involved:

 * ComponentImpl -- this defines a component implementation.
   It has a single name (ComponentPath) and is instantiated on demand.
   This is similar to a class in OOP.
   Typically this is an empty type, exposed as a public package variable.

 * Component -- this defines an instance of a component, created when another component depends on it.
   The system currently instantiates a component only once, but other strategies might be applied.
   This is similar to an instance of a class in OOP.
 
 * ComponentReference -- this defines a reference to a (possibly remote) component, mediating communication with that component.
   This is similar to an instance pointer in OOP.

# TODO

## Convenience Wrappers

Logging is not particularly ergonomic:
```go
e.logger.Request(ctx, logger.Output{Message: "hello, world"})
```

Maybe a component can provide a convenience wrapper around its ComponentReference type that allows more ergonomic calls:

```go
func (*impl) Start(deps map[core.ComponentPath]core.ComponentReference) core.Component {
	return &comp{
        logger: logger.Wrapper(deps),
    }
}
// ...
e.logger.Output("hello, world")
```
