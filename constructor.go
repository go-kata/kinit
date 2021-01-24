package kinit

import (
	"reflect"

	"github.com/go-kata/kdone"
)

// Constructor represents an object constructor.
//
// The usual identifier for variables of this type is ctor.
type Constructor interface {
	// Type returns a type of an object that is created by this constructor.
	Type() reflect.Type
	// Parameters returns types of objects this constructor depends on.
	Parameters() []reflect.Type
	// Create creates and returns a new object.
	Create(a ...reflect.Value) (reflect.Value, kdone.Destructor, error)
}
