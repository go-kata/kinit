# KInit

[![Go Reference](https://pkg.go.dev/badge/github.com/go-kata/kinit.svg)](https://pkg.go.dev/github.com/go-kata/kinit)
[![codecov](https://codecov.io/gh/go-kata/kinit/branch/master/graph/badge.svg?token=NBFR4LKON8)](https://codecov.io/gh/go-kata/kinit)
[![Report Card](https://goreportcard.com/badge/github.com/go-kata/kinit)](https://goreportcard.com/report/github.com/go-kata/kinit)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

**[Usage examples](https://github.com/go-kata/examples)**

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
the automatic [Dependency Injection](https://en.wikipedia.org/wiki/Dependency_injection):

```go
kinit.Global()
```

If you need a local one (e.g. for tests), you can create it as following:

```go
ctr := kinit.NewContainer()
```

A local container must be filled up with *constructors* and *processors* manually whereas the global one
can be filled up when [initializing packages](https://golang.org/doc/effective_go.html#init).

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

To register a constructor in the container use the `Provide` method.

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

To register a processor in the container use the `Attach` method.

The container allows to have an unlimited number of processors for each type but doesn't guarantee the order of
their calling.

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

To run functors in the container use the `Run` method.

At the start of run the container creates so-called *arena* that holds all created objects (only one object
for each type). If some object required as a dependency is already on the arena it will be used, otherwise
it will be firstly created and processed. All objects that are on the arena at the end of run will be
automatically destroyed.

The container runs given functors sequentially. Their dependencies are resolved recursively using registered
constructors and processors. If functor (let's call it *branched*) returns further functors, the container runs
all of them before continue running functors following the branched one. This is called the *Depth-First Run*.
  
## KInitX

[![Go Reference](https://pkg.go.dev/badge/github.com/go-kata/kinit/kinitx.svg)](https://pkg.go.dev/github.com/go-kata/kinit/kinitx)

This subpackage provides the expansion set includes default handy implementations of main library interfaces
along with other handy tools. In most cases the **KInitX** is all you need to use the entire **KInit** functionality.

There are following implementations:

**Constructor** represents a constructor based on a function. It accepts `func(...) T`, `func(...) (T, error)` and
`func(...) (T, kdone.Destructor, error)` signatures where `T` is an arbitrary Go type.

```go
kinitx.MustProvide(func(config *Config) (*Object, kdone.Destructor, error) { ... })
```

**Opener** represents a constructor based on a function that creates an implementation of the `io.Closer` interface.
It accepts `func(...) C` and `func(...) (C, error)` signatures where `C` is an arbitrary implementation of the
`io.Closer` interface.

```go
kinitx.MustProvide(func(logger *log.Logger) (*sql.DB, error) { ... })
```

**Initializer** represents a memberwise initializer of a struct. It accepts a template struct like a `YourType{}`
and a template struct pointer like a `(*YourType)(nil)` or `new(YourType)`.

```go
kinitx.MustProvide((*Config)(nil))
```

**Binder** represents a pseudo-constructor that casts an object to an interface. It accepts an interface pointer
like a `(*YourInterface)(nil)`.

```go
kinitx.MustBind((*StorageInterface)(nil), (*PostgresStrorage)(nil))
```

**Processor** represents a processor based on a function. It accepts `func(T, ...)` and `func(T, ...) error`
signatures where `T` is an arbitrary Go type.

```go
kinitx.MustAttach((*Object).SetOptionalProperty)
```

**Functor** represents a functor based on a function. It accepts `func(...)`, `func(...) error`,
`func(...) (kinit.Functor, error)` and `func(...) ([]kinit.Functor, error)` signatures.

```go
kinitx.MustRun(func(app *Application) error { ... })
```

## KInitQ

[![Go Reference](https://pkg.go.dev/badge/github.com/go-kata/kinit/kinitq.svg)](https://pkg.go.dev/github.com/go-kata/kinit/kinitq)

The DI mechanism provided by the main library is reflection-based and works in the runtime. However, this subpackage
makes it possible to validate the dependency graph semi-statically thanks to build tags.

Just add two main functions as following:

`main.go`

```go
// +build !inspect

package main

import "github.com/go-kata/kinit/kinitx"

func main() { kinitx.MustRun(EntryPoint) }
```

`main_inspect.go`

```go
// +build inspect

package main

import "github.com/go-kata/kinit/kinitx"

func main() { kinitx.MustInspect(nil) }
```

Now to validate the dependency graph of your program just run:

`go run -tags inspect`

Example output:

```
2 errors occurred:
    #1 🠖 cyclic dependency: *config.Config 🠖 config.Loader 🠖 *config.FileLoader 🠖 *config.Config
    #2 🠖 unsatisfied dependency: *sql.DB 🠖 *log.Logger
```

For more details learn the documentation and explore examples.

## Putting all together

In the [github.com/go-kata/examples](https://github.com/go-kata/examples) repository you can find examples of
how may the code uses this library looks like.

## References

**[KDone](https://github.com/go-kata/kdone)** is the library that provides tools for destroying objects.

**[KError](https://github.com/go-kata/kerror)** is the library that provides tools for handling errors.
