# KInit

[![Go Reference](https://pkg.go.dev/badge/github.com/go-kata/kinit.svg)](https://pkg.go.dev/github.com/go-kata/kinit)
[![codecov](https://codecov.io/gh/go-kata/kinit/branch/master/graph/badge.svg?token=NBFR4LKON8)](https://codecov.io/gh/go-kata/kinit)

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

A local container must be filled up with *constructors* and *processors* manually whereas the global one can be
filled up when initializing packages. You may use init functions fot this, but it's recommended to use *hooks*.

### Hooks

Registered hooks will be called only when the global container starts an *invocation*. It's useful for libraries
which may not perform container filling up operations if their entities just used manually.

To register a hook use `kinit.Hook` or `kinit.MustHook` method:

```go
kinit.Hook(func() error { /* fill up the global container here with returning error if occurred */ })

kinit.MustHook(func() { /* fill up the global container here */ })
```

Both methods return an index of registered hook starting from zero. It is not very useful information, but
you can use this fact to simplify syntax:

```go
var _ = kinit.MustHook(func() { ... })
```

instead of

```go
func init() {
	kinit.MustHook(func() { ... })
}
```

### Constructors

Constructors are entities which creates objects (dependencies for injection in context of the DI). This library
considers constructors to have the following interface:

```go
type Constructor interface {
	
	Type() reflect.Type
	
	Parameters() []reflect.Type
	
	Create(a ...reflect.Value) (reflect.Value, kdone.Destructor, error)
	
}
```

Here `Type` returns a type of object to create, `Parameters` returns types of other objects required to this
object creation (dependencies) and finally `Create` creates and returns a new object and its destructor (use
`kdone.Noop`, not nil, if object has no destructor).

To register a new constructor in container use `kinit.Provide` method (or `ctr.Provide` for local containers).
There are also `kinit.MustProvide` method (or `ctr.MustProvide` for local containers) that panics on error.

Container allows to have only one constructor for each type. However, the `reflect.Type` is the interface, and
you can implement it as you want, e.g. as following:

```go
type NamedType stuct {
	reflect.Type
	Name string
}
```

Just keep in mind that container uses a simple comparison of `reflect.Type` instances to find necessary constructors
(as well as processors and already created objects).

The **[KInitX](https://github.com/go-kata/kinitx)** library provides following constructor implementations:

* **Constructor** is the function-based constructor.
* **Opener** is the function-based constructor which creates an object that implements the `io.Closer` interface.
  The `Close` method is treated as object destructor.
* **Initializer** is the memberwise initializer for structs that doesn't provide a destructor.

You can find more details in the documentation for the library.

### Processors

Processors are entities which processes already created objects. This library applies processors immediately after
object creation and before it will be injected as a dependency at the first time. It considers processors to have
the following interface:

```go
type Processor interface {
	
	Type() reflect.Type
	
	Parameters() []reflect.Type

	Process(obj reflect.Value, a ...reflect.Value) error
	
}
```

Here `Type` returns a type of object to process, `Parameters` returns types of other objects required to this
object processing (dependencies) and finally `Process` processes an object.

To register a new processor in container use `kinit.Apply` method (or `ctr.Apply` for local containers).
There are also `kinit.MustApply` method (or `ctr.MustApply` for local containers) that panics on error.

Container allows to have an unlimited number of processors for each type but doesn't guarantee the order of their
calling.

The **[KInitX](https://github.com/go-kata/kinitx)** library provides the function-based processor implementation.

### Invocation

Invocation is the process of some activity execution resolving a dependency tree as roots of which are treated
this activity dependencies.

To perform an invocation use `kinit.Invoke` method (or `ctr.Invoke` for local containers). There are also
`kinit.MustInvoke` method (or `ctr.MustInvoke` for local containers) that panics on error. Both methods
require an *executor* and allow to pass one or more *bootstrappers*.

At the start of invocation container creates so-called *arena* that contains all created objects (only one object
for each type). If some object that is required as a dependency is already on the arena it will be used, otherwise
it will be previously created and processed. All objects that are on the arena at the end of invocation will be
automatically destroyed using their destructors.

When all dependencies of an activity are resolved, container executes it. Executed activity may provide other
activity to continue invocation - its dependencies will be also resolved using the same arena. This process is
called the *cascade injection*. When a currently executed activity doesn't provide a next activity to execute,
invocation ends.

### Executors

Executors are representations of activities. This library considers executors to have the following interface:

```go
type Executor interface {
	
	Parameters() []reflect.Type
	
	Execute(a ...reflect.Value) (Executor, error)
	
}
```

Here `Parameters` returns type of objects required to activity execution (root dependencies) and `Execute` executes
activity and may return a next executor to continue invocation.

The **[KInitX](https://github.com/go-kata/kinitx)** library provides the function-based executor implementation.

### Bootstrappers

Bootstrappers are special entities which allows the arena bootstrapping at the start of invocation. This library
considers bootstrappers to have the following interface:

```go
type Bootstrapper interface {
	
	Bootstrap(arena *kinit.Arena) error
	
}
```

The **[KInitX](https://github.com/go-kata/kinitx)** library provides the bootstrapper implementation called *literal*.
Literals are objects that are registered on the arena for direct use instead of being created during dependency tree
resolution. Those objects have no destructors to call at the end of invocation - they must be destroyed manually.

### Putting all together

In the [github.com/go-kata/examples](https://github.com/go-kata/examples) repository you can find examples of how may
the code uses this library looks like.

## References

**[KInitX](https://github.com/go-kata/kinitx)** is the library that provides default extensions for the **KInit**.

**[KDone](https://github.com/go-kata/kdone)** is the library that provides tools for destroying objects.

**[KError](https://github.com/go-kata/kerror)** is the library that provides tools for handling errors.
