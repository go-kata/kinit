package kinit

import "reflect"

// Executor represents an activity executor.
//
// The usual identifier for variables of this type is exec.
//
// Deprecated: since 0.4.0, use Functor instead.
type Executor interface {
	// Parameters returns types of objects this executor depends on.
	Parameters() []reflect.Type
	// Execute executes an activity and returns a next executor
	// to continue invocation or nil to stop it.
	Execute(a ...reflect.Value) (Executor, error)
}
