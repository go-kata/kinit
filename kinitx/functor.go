package kinitx

import (
	"reflect"

	"github.com/go-kata/kerror"
	"github.com/go-kata/kinit"
)

// Functor represents a functor based on a function.
type Functor struct {
	// function specifies the reflection to a function value.
	function reflect.Value
	// inTypes specifies types of function input parameters.
	inTypes []reflect.Type
	// furtherOutIndex specifies the index of a function output parameter that contains further functor(s).
	// The value -1 means that a function doesn't return subsequent functor(s).
	furtherOutIndex int
	// furtherIsSlice specifies whether the function returns multiple further functors.
	furtherIsSlice bool
	// errorOutIndex specifies the index of a function output parameter that contains an error.
	// The value -1 means that a function doesn't return an error.
	errorOutIndex int
}

// NewFunctor returns a new functor.
//
// The argument x must be a function that is compatible with one of following signatures:
//
//     func(...)
//
//     func(...) error
//
//     func(...) (kinit.Functor, error)
//
//     func(...) ([]kinit.Functor, error)
//
func NewFunctor(x interface{}) (*Functor, error) {
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
	f := &Functor{
		function: fv,
	}
	numIn := ft.NumIn()
	if ft.IsVariadic() {
		numIn--
	}
	f.inTypes = make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		f.inTypes[i] = ft.In(i)
	}
	switch ft.NumOut() {
	default:
		return nil, kerror.Newf(kerror.EViolation, "function %s is not a functor", ft)
	case 0:
		f.furtherOutIndex = -1
		f.furtherIsSlice = false
		f.errorOutIndex = -1
	case 1:
		if ft.Out(0) != errorType {
			return nil, kerror.Newf(kerror.EViolation, "function %s is not a functor", ft)
		}
		f.furtherOutIndex = -1
		f.furtherIsSlice = false
		f.errorOutIndex = 0
	case 2:
		out0 := ft.Out(0)
		out0IsFunctorSlice := out0 == functorSliceType
		if !(out0IsFunctorSlice || out0 == functorType) || ft.Out(1) != errorType {
			return nil, kerror.Newf(kerror.EViolation, "function %s is not a functor", ft)
		}
		f.furtherOutIndex = 0
		f.furtherIsSlice = out0IsFunctorSlice
		f.errorOutIndex = 1
	}
	return f, nil
}

// MustNewFunctor is a variant of the NewFunctor that panics on error.
func MustNewFunctor(x interface{}) *Functor {
	f, err := NewFunctor(x)
	if err != nil {
		panic(err)
	}
	return f
}

// Parameters implements the kinit.Functor interface.
func (f *Functor) Parameters() []reflect.Type {
	if f == nil {
		return nil
	}
	types := make([]reflect.Type, len(f.inTypes))
	copy(types, f.inTypes)
	return types
}

// Call implements the kinit.Functor interface.
func (f *Functor) Call(a ...reflect.Value) ([]kinit.Functor, error) {
	if f == nil {
		return nil, nil
	}
	if len(a) != len(f.inTypes) {
		return nil, kerror.Newf(kerror.EViolation,
			"functor expects %d argument(s), %d given", len(f.inTypes), len(a))
	}
	for i, v := range a {
		if v.Type() != f.inTypes[i] {
			return nil, kerror.Newf(kerror.EViolation,
				"functor expects argument %d to be of %s type, %s given",
				i+1, f.inTypes[i], v.Type())
		}
	}
	out := f.function.Call(a)
	var further []kinit.Functor
	if f.furtherOutIndex >= 0 {
		if v := out[f.furtherOutIndex].Interface(); v != nil {
			if f.furtherIsSlice {
				further = v.([]kinit.Functor)
			} else {
				further = []kinit.Functor{v.(kinit.Functor)}
			}
		}
	}
	var err error
	if f.errorOutIndex >= 0 {
		if v := out[f.errorOutIndex].Interface(); v != nil {
			err = v.(error)
		}
	}
	return further, err
}
