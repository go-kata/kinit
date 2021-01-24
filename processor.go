package kinit

import "reflect"

// Processor represents an object processor.
//
// The usual identifier for variables of this type is proc.
type Processor interface {
	// Type returns a type of an object that is processed by this processor.
	Type() reflect.Type
	// Parameters returns types of objects this processor depends on.
	Parameters() []reflect.Type
	// Process processes the given object.
	Process(obj reflect.Value, a ...reflect.Value) error
}
