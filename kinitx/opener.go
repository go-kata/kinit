package kinitx

import (
	"io"
	"reflect"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

// Opener represents a constructor based on a function that creates
// an implementation of the io.Closer interface.
type Opener struct {
	// t specifies the type of an object that is created by this opener.
	t reflect.Type
	// function specifies the reflection to a function value.
	function reflect.Value
	// inTypes specifies types of function input parameters.
	inTypes []reflect.Type
	// objectOutIndex specifies the index of a function output parameter that contains a created object.
	objectOutIndex int
	// errorOutIndex specifies the index of a function output parameter that contains an error.
	// The value -1 means that a function doesn't return an error.
	errorOutIndex int
}

// NewOpener returns a new opener.
//
// The argument x must be a function that is compatible with one of following signatures
// (C is an arbitrary implementation of the io.Closer interface):
//
//     func(...) C;
//
//     func(...) (C, error);
//
func NewOpener(x interface{}) (*Opener, error) {
	if x == nil {
		return nil, kerror.New(kerror.EViolation, "function expected, nil given")
	}
	ft := reflect.TypeOf(x)
	fv := reflect.ValueOf(x)
	if ft.Kind() != reflect.Func {
		return nil, kerror.Newf(kerror.EViolation, "function expected, %s given", ft)
	}
	if fv.IsNil() {
		return nil, kerror.New(kerror.EViolation, "function expected, nil given")
	}
	o := &Opener{
		function: fv,
	}
	numIn := ft.NumIn()
	if ft.IsVariadic() {
		numIn--
	}
	o.inTypes = make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		o.inTypes[i] = ft.In(i)
	}
	switch ft.NumOut() {
	default:
		return nil, kerror.Newf(kerror.EViolation, "function %s is not an opener", ft)
	case 1:
		if !ft.Out(0).Implements(closerType) {
			return nil, kerror.Newf(kerror.EViolation, "function %s is not an opener", ft)
		}
		o.t = ft.Out(0)
		o.objectOutIndex = 0
		o.errorOutIndex = -1
	case 2:
		if !ft.Out(0).Implements(closerType) || ft.Out(1) != errorType {
			return nil, kerror.Newf(kerror.EViolation, "function %s is not an opener", ft)
		}
		o.t = ft.Out(0)
		o.objectOutIndex = 0
		o.errorOutIndex = 1
	}
	return o, nil
}

// MustNewOpener is a variant of the NewOpener that panics on error.
func MustNewOpener(x interface{}) *Opener {
	o, err := NewOpener(x)
	if err != nil {
		panic(err)
	}
	return o
}

// Type implements the kinit.Constructor interface.
func (o *Opener) Type() reflect.Type {
	if o == nil {
		return nil
	}
	return o.t
}

// Parameters implements the kinit.Constructor interface.
func (o *Opener) Parameters() []reflect.Type {
	if o == nil {
		return nil
	}
	types := make([]reflect.Type, len(o.inTypes))
	copy(types, o.inTypes)
	return types
}

// Create implements the kinit.Constructor interface.
func (o *Opener) Create(a ...reflect.Value) (reflect.Value, kdone.Destructor, error) {
	if o == nil {
		return reflect.Value{}, kdone.Noop, nil
	}
	if len(a) != len(o.inTypes) {
		return reflect.Value{}, nil, kerror.Newf(kerror.EViolation,
			"%s opener expects %d argument(s), %d given",
			o.t, len(o.inTypes), len(a))
	}
	for i, v := range a {
		if v.Type() != o.inTypes[i] {
			return reflect.Value{}, nil, kerror.Newf(kerror.EViolation,
				"%s opener expects argument %d to be of %s type, %s given",
				o.t, i+1, o.inTypes[i], v.Type())
		}
	}
	out := o.function.Call(a)
	obj := out[o.objectOutIndex]
	dtor := kdone.DestructorFunc(obj.Interface().(io.Closer).Close)
	var err error
	if o.errorOutIndex >= 0 {
		if v := out[o.errorOutIndex].Interface(); v != nil {
			err = v.(error)
		}
	}
	return obj, dtor, err
}
