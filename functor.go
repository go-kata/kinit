package kinit

import "reflect"

// Functor represents a function.
//
// The usual identifier for variables of this type is fun.
type Functor interface {
	// Parameters returns types of objects this functor depends on.
	Parameters() []reflect.Type
	// Call calls a function and may return further functors.
	Call(a ...reflect.Value) ([]Functor, error)
}
