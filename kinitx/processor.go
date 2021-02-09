package kinitx

import (
	"reflect"

	"github.com/go-kata/kerror"
)

// Processor represents a processor based on a function.
type Processor struct {
	// t specifies the type of an object that is processed by this processor.
	t reflect.Type
	// function specifies the reflection to a function value.
	function reflect.Value
	// inTypes specifies types of function input parameters.
	inTypes []reflect.Type
	// errorOutIndex specifies the index of a function output parameter that contains an error.
	// The value -1 means that a function doesn't return an error.
	errorOutIndex int
}

// NewProcessor returns a new processor.
//
// The argument x must be a function that is compatible with one of following signatures
// (T is an arbitrary Go type):
//
//     func(T, ...);
//
//     func(T, ...) error.
//
func NewProcessor(x interface{}) (*Processor, error) {
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
	p := &Processor{
		function: fv,
	}
	numIn := ft.NumIn()
	if ft.IsVariadic() {
		numIn--
	}
	if numIn < 1 {
		return nil, kerror.Newf(kerror.EViolation, "function %s is not a processor", ft)
	}
	p.t = ft.In(0)
	p.inTypes = make([]reflect.Type, numIn-1)
	for i := 1; i < numIn; i++ {
		p.inTypes[i-1] = ft.In(i)
	}
	switch ft.NumOut() {
	default:
		return nil, kerror.Newf(kerror.EViolation, "function %s is not a processor", ft)
	case 0:
		p.errorOutIndex = -1
	case 1:
		if ft.Out(0) != errorType {
			return nil, kerror.Newf(kerror.EViolation, "function %s is not a processor", ft)
		}
		p.errorOutIndex = 0
	}
	return p, nil
}

// MustNewProcessor is a variant of the NewProcessor that panics on error.
func MustNewProcessor(x interface{}) *Processor {
	p, err := NewProcessor(x)
	if err != nil {
		panic(err)
	}
	return p
}

// Type implements the kinit.Processor interface.
func (p *Processor) Type() reflect.Type {
	if p == nil {
		return nil
	}
	return p.t
}

// Parameters implements the kinit.Processor interface.
func (p *Processor) Parameters() []reflect.Type {
	if p == nil {
		return nil
	}
	types := make([]reflect.Type, len(p.inTypes))
	copy(types, p.inTypes)
	return types
}

// Process implements the kinit.Processor interface.
func (p *Processor) Process(obj reflect.Value, a ...reflect.Value) error {
	if p == nil {
		return nil
	}
	if obj.Type() != p.t {
		return kerror.Newf(kerror.EViolation,
			"%s processor doesn't accept objects of %s type", p.t, obj.Type())
	}
	if len(a) != len(p.inTypes) {
		return kerror.Newf(kerror.EViolation,
			"%s processor expects %d argument(s), %d given", p.t, len(p.inTypes), len(a))
	}
	in := make([]reflect.Value, len(p.inTypes)+1)
	in[0] = obj
	for i, v := range a {
		if v.Type() != p.inTypes[i] {
			return kerror.Newf(kerror.EViolation,
				"%s processor expects argument %d to be of %s type, %s given",
				p.t, i+1, p.inTypes[i], v.Type())
		}
		in[i+1] = a[i]
	}
	out := p.function.Call(in)
	var err error
	if p.errorOutIndex >= 0 {
		if v := out[p.errorOutIndex].Interface(); v != nil {
			err = v.(error)
		}
	}
	return err
}
