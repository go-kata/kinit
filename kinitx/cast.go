package kinitx

import (
	"reflect"

	"github.com/go-kata/kerror"
	"github.com/go-kata/kinit"
)

// castToConstructor returns a constructor based on the given entity.
//
// See the documentation for the Provide to find out possible values of the argument x.
func castToConstructor(x interface{}) (kinit.Constructor, error) {
	if x == nil {
		return nil, kerror.New(kerror.EViolation, "function, struct or struct pointer expected, nil given")
	}
	if ctor, ok := x.(kinit.Constructor); ok {
		return ctor, nil
	}
	var ctor kinit.Constructor
	var err error
	t := reflect.TypeOf(x)
	switch t.Kind() {
	default:
		return nil, kerror.Newf(kerror.EViolation, "function, struct or struct pointer expected, %s given", t)
	case reflect.Func:
		var isOpener bool
		switch t.NumOut() {
		case 2:
			if t.Out(1) != errorType {
				break
			}
			fallthrough
		case 1:
			isOpener = t.Out(0).Implements(closerType)
		}
		if isOpener {
			ctor, err = NewOpener(x)
		} else {
			ctor, err = NewConstructor(x)
		}
	case reflect.Struct, reflect.Ptr:
		ctor, err = NewInitializer(x)
	}
	return ctor, err
}

// castToProcessor returns a processor based on the given entity.
//
// See the documentation for the Attach to find out possible values of the argument x.
func castToProcessor(x interface{}) (kinit.Processor, error) {
	if x == nil {
		return nil, kerror.New(kerror.EViolation, "function expected, nil given")
	}
	if proc, ok := x.(kinit.Processor); ok {
		return proc, nil
	}
	return NewProcessor(x)
}

// castToFunctor returns a new functor based on the given entity.
//
// See the documentation for the Run to find out possible values of the argument x.
func castToFunctor(x interface{}) (kinit.Functor, error) {
	if x == nil {
		return nil, kerror.New(kerror.EViolation, "function or value expected, nil given")
	}
	if fun, ok := x.(kinit.Functor); ok {
		return fun, nil
	}
	var fun kinit.Functor
	var err error
	if reflect.TypeOf(x).Kind() == reflect.Func {
		fun, err = NewFunctor(x)
	} else {
		fun, err = NewInjector(x)
	}
	return fun, err
}
