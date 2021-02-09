// Package kinitx provides the KInit expansion set.
package kinitx

import (
	"io"
	"reflect"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
	"github.com/go-kata/kinit"
	"github.com/go-kata/kinit/kinitq"
)

// errorType specifies the reflection to the error interface.
var errorType = reflect.TypeOf((*error)(nil)).Elem()

// closerType specifies the reflection to the io.Closer interface.
var closerType = reflect.TypeOf((*io.Closer)(nil)).Elem()

// destructorType specifies the reflection to the kdone.Destructor interface.
var destructorType = reflect.TypeOf((*kdone.Destructor)(nil)).Elem()

// functorType specifies the reflection to the kinit.Functor interface.
var functorType = reflect.TypeOf((*kinit.Functor)(nil)).Elem()

// functorSliceType specifies the reflection to the slice of functors.
var functorSliceType = reflect.SliceOf(functorType)

// runtimeType specifies the reflection to the kinit.Runtime interface.
var runtimeType = reflect.TypeOf((*kinit.Runtime)(nil))

// Provide calls the Provide method of the global container by passing a constructor based on the given entity.
//
// The x argument will be parsed corresponding to following rules:
//
// - x must not be nil;
//
// - if x implements the kinit.Constructor interface it will used by itself;
//
// - if x is a function it will be parsed using the NewOpener only when returns
//   an implementation of the io.Closer interface at the first position and, optionally,
//   error at the second; all other functions will be parsed using the NewConstructor;
//
// - if x is a struct or pointer it will be parsed using the NewInitializer;
//
// - all other variants of x are unacceptable.
func Provide(x interface{}) error {
	ctor, err := castToConstructor(x)
	if err != nil {
		return err
	}
	return kinit.Global().Provide(ctor)
}

// MustProvide is a variant of the Provide that panics on error.
func MustProvide(x interface{}) {
	if err := Provide(x); err != nil {
		panic(err)
	}
}

// Bind calls the Provide method of the global container bu passing a binder based on given interface and object.
//
// See the documentation for the NewBinder to find out possible values of the argument x.
func Bind(i, x interface{}) error {
	ctor, err := NewBinder(i, x)
	if err != nil {
		return err
	}
	return kinit.Global().Provide(ctor)
}

// MustBind is a variant of the Bind that panics on error.
func MustBind(i, x interface{}) {
	if err := Bind(i, x); err != nil {
		panic(err)
	}
}

// Attach calls the Attach method of the global container by passing a processor based on the given entity.
//
// The x argument will be parsed corresponding to following rules:
//
// - x must not be nil;
//
// - if x implements the kinit.Processor interface it will be used by itself;
//
// - otherwise x will be parsed using the NewProcessor.
func Attach(x interface{}) error {
	proc, err := castToProcessor(x)
	if err != nil {
		return err
	}
	return kinit.Global().Attach(proc)
}

// MustAttach is a variant of the Attach that panics on error.
func MustAttach(x interface{}) {
	if err := Attach(x); err != nil {
		panic(err)
	}
}

// Run calls the Run method of the global container by passing functors based on given entities.
//
// Items of the xx argument (let's name each item as x) will be parsed corresponding to following rules:
//
// - x must not be nil;
//
// - if x is a function it will be parsed using the NewFunctor;
//
// - otherwise x will be parsed using the NewInjector.
func Run(xx ...interface{}) error {
	functors := make([]kinit.Functor, len(xx))
	for i, x := range xx {
		fun, err := castToFunctor(x)
		if err != nil {
			return err
		}
		functors[i] = fun
	}
	return kinit.Global().Run(functors...)
}

// MustRun is a variant of the Run that panics on error.
func MustRun(xx ...interface{}) {
	if err := Run(xx...); err != nil {
		panic(err)
	}
}

// Require calls the Require method of the global inspector by passing the type of the given entity.
//
// The argument x must not be nil.
func Require(x interface{}) error {
	if x == nil {
		return kerror.New(kerror.EViolation, "value expected, nil given")
	}
	return kinitq.Global().Require(reflect.TypeOf(x))
}

// MustRequire is a variant of the Require that panics on error.
func MustRequire(x interface{}) {
	if err := Require(x); err != nil {
		panic(err)
	}
}

// Ignore calls the Ignore method of the global inspector by passing the type of the given entity.
//
// The argument x must not be nil.
func Ignore(x interface{}) error {
	if x == nil {
		return kerror.New(kerror.EViolation, "value expected, nil given")
	}
	return kinitq.Global().Ignore(reflect.TypeOf(x))
}

// MustIgnore is a variant of the Ignore that panics on error.
func MustIgnore(x interface{}) {
	if err := Ignore(x); err != nil {
		panic(err)
	}
}

// Consider calls the Require method of the global inspector
// for each dependency of the functor based on the given entity.
//
// See the documentation for the Run to find out possible values of the argument x.
func Consider(x interface{}) error {
	fun, err := castToFunctor(x)
	if err != nil {
		return err
	}
	for _, p := range fun.Parameters() {
		if err := kinitq.Global().Require(p); err != nil {
			return err
		}
	}
	return nil
}

// MustConsider is a variant of the Consider that panics on error.
func MustConsider(x interface{}) {
	if err := Consider(x); err != nil {
		panic(err)
	}
}

// Inspect calls the Inspect method of the global inspector on the global container.
func Inspect(opt *kinitq.Options) error {
	return kinitq.Global().Inspect(kinit.Global(), opt)
}

// MustInspect is a variant of the Inspect that panics on error.
func MustInspect(opt *kinitq.Options) {
	if err := Inspect(opt); err != nil {
		panic(err)
	}
}
