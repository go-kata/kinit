# KInit

[![Go Reference](https://pkg.go.dev/badge/github.com/go-kata/kinit.svg)](https://pkg.go.dev/github.com/go-kata/kinit)
[![codecov](https://codecov.io/gh/go-kata/kinit/branch/master/graph/badge.svg?token=NBFR4LKON8)](https://codecov.io/gh/go-kata/kinit)
[![Report Card](https://goreportcard.com/badge/github.com/go-kata/kinit)](https://goreportcard.com/report/github.com/go-kata/kinit)

## Installation

`go get github.com/go-kata/kinit`

## Status

**This is a beta version.** API is not stabilized for now.

## Versioning

*Till the first major release* minor version (`v0.x.0`) must be treated as a major version
and patch version (`v0.0.x`) must be treated as a minor version.

For example, changing of version from `v0.1.0` to `v0.1.1` indicates compatible changes,
but when version changes `v0.1.1` to `v0.2.0` this means that the last version breaks the API.

## How to use

This library provides the global [IoC](https://en.wikipedia.org/wiki/Inversion_of_control) container which does
the automatic [Dependency Injection](https://en.wikipedia.org/wiki/Dependency_injection). If you need a local one
(e.g. for tests), you can create it as following:

```go
ctr := kinit.NewContainer()
```

A local container must be filled up with *constructors* and *processors* manually whereas the global one
can be filled up when initializing packages. You may use init functions for this, but it's recommended
to use *declared functions*.

> Further, when referring to the *container*, the global container is meant. To use some method for a
> local container just replace `kinit.FunctionName` to `ctr.FunctionName`.

### Declared functions

Declared functions will be called only when the container runs *functors* at the first time. It's useful
for libraries which may not perform container filling up operations if their entities just used manually.

To declare a function use `kinit.Declare` and `kinit.DeclareErrorProne` methods:

```go
kinit.Declare(func() { /* fill up the container here */ })

kinit.DeclareErrorProne(func() error { /* fill up the container here with returning error if occurred */ })
```

There are also `kinit.MustDeclare` and `kinit.MustDeclareErrorProne` methods that panic on error. Both methods
return `struct{}` - it's not very informative result, but you can use this fact to simplify syntax:

```go
var _ = kinit.MustDeclare(func() { ... })
```

instead of

```go
func init() {
	kinit.MustDeclare(func() { ... })
}
```

> Declared functions are not applicable for local containers.

> All the methods mentioned below also have `Must` versions that panic on error but don't return `struct{}`.

### Constructors

Constructors are entities that create *objects* (dependencies for injection in context of the DI). They have the
following interface:

```go
type Constructor interface {
	
	Type() reflect.Type
	
	Parameters() []reflect.Type
	
	Create(a ...reflect.Value) (reflect.Value, kdone.Destructor, error)
	
}
```

Here `Type` returns a type of object to create, `Parameters` returns types of other objects required to create
this object and finally `Create` creates and returns a new object and its destructor (use`kdone.Noop`, not nil,
if object has no destructor).

To register a constructor in the container use `kinit.Provide` method.

The container allows to have only one constructor for each type. However, the `reflect.Type` is the interface, and
you can implement it as you want, e.g. as following:

```go
type NamedType stuct {
	reflect.Type
	Name string
}
```

Just keep in mind that the container uses a simple comparison of `reflect.Type` instances when looks up for
necessary constructors (as well as processors and already created objects).

In addition to the function-based constructor implementation the **[KInitX](https://github.com/go-kata/kinitx)**
library provides following implementations:

* **Opener** is the function-based constructor creating an object that implements the `io.Closer` interface.
  The object's `Close` method is treated as a destructor.
* **Initializer** is the memberwise initializer of a struct. It doesn't provide a destructor.

You can find more details in the [documentation](https://pkg.go.dev/github.com/go-kata/kinitx) for the library.

### Processors

Processors are entities that process already created objects. The container applies processors immediately after
the object creation and before an object will be injected as a dependency at the first time. Processors have the
following interface:

```go
type Processor interface {
	
	Type() reflect.Type
	
	Parameters() []reflect.Type

	Process(obj reflect.Value, a ...reflect.Value) error
	
}
```

Here `Type` returns a type of object to process, `Parameters` returns types of other objects required to process
this object and finally `Process` processes an object.

To register a processor in the container use `kinit.Attach` method.

The container allows to have an unlimited number of processors for each type but doesn't guarantee the order of
their calling.

The **[KInitX](https://github.com/go-kata/kinitx)** library provides the function-based processor implementation only.

### Functors

Functors represent functions to be run in the container and have the following interface:

```go
type Functor interface {
	
	Parameters() []reflect.Type
	
	Call(a ...reflect.Value) ([]Functor, error)
	
}
```

Here `Parameters` returns types of objects required to call a function and `Call` calls a function and may return
*further functors*.

To run functors in the container use `kinit.Run` method.

At the start of run the container creates so-called *arena* that holds all created objects (only one object
for each type). If some object required as a dependency is already on the arena it will be used, otherwise
it will be firstly created and processed. All objects that are on the arena at the end of run will be
automatically destroyed.

The container runs given functors sequentially. Their dependencies are resolved recursively using registered
constructors and processors. If functor (let's call it *branched*) returns further functors, the container runs
all of them before continue running functors following the branched one. This is called the *Depth-First Run*.

In addition to the function-based functor implementation the **[KInitX](https://github.com/go-kata/kinitx)**
library provides following implementations:

* **Injector** is the provider of object that is directly registered on the arena instead of being created
  during the dependency tree resolution. The provided object must be destroyed manually after run ends.

### Putting all together

In the [github.com/go-kata/examples](https://github.com/go-kata/examples) repository you can find examples of
how may the code uses this library looks like.

## References

**[KInitX](https://github.com/go-kata/kinitx)** is the library that provides default extensions for the **KInit**.

**[KDone](https://github.com/go-kata/kdone)** is the library that provides tools for destroying objects.

**[KError](https://github.com/go-kata/kerror)** is the library that provides tools for handling errors.
