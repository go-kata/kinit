package kinitx

import (
	"reflect"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

// Binder represents a pseudo-constructor that casts an object to an interface.
type Binder struct {
	// t specifies the type of an object that is created by this binder.
	t reflect.Type
	// inType specifies the type of the input object to cast.
	inType reflect.Type
}

// NewBinder returns a new binder.
//
// The argument i must be an interface pointer and the argument x must not be nil.
func NewBinder(i, x interface{}) (*Binder, error) {
	if i == nil {
		return nil, kerror.New(kerror.EViolation, "interface pointer expected, nil given")
	}
	if x == nil {
		return nil, kerror.New(kerror.EViolation, "value expected, nil given")
	}
	pt := reflect.TypeOf(i)
	if pt.Kind() != reflect.Ptr {
		return nil, kerror.Newf(kerror.EViolation, "interface pointer expected, %s given", pt)
	}
	it := pt.Elem()
	if it.Kind() != reflect.Interface {
		return nil, kerror.Newf(kerror.EViolation, "interface pointer expected, %s given", pt)
	}
	ot := reflect.TypeOf(x)
	if !ot.Implements(it) {
		return nil, kerror.Newf(kerror.EViolation, "%s doesn't implement %s", ot, pt)
	}
	return &Binder{
		t:      it,
		inType: ot,
	}, nil
}

// MustNewBinder is a variant of the NewBinder that panics on error.
func MustNewBinder(i, x interface{}) *Binder {
	b, err := NewBinder(i, x)
	if err != nil {
		panic(err)
	}
	return b
}

// Type implements the kinit.Constructor interface.
func (b *Binder) Type() reflect.Type {
	if b == nil {
		return nil
	}
	return b.t
}

// Parameters implements the kinit.Constructor interface.
func (b *Binder) Parameters() []reflect.Type {
	if b == nil {
		return nil
	}
	return []reflect.Type{b.inType}
}

// Create implements the kinit.Constructor interface.
func (b *Binder) Create(a ...reflect.Value) (reflect.Value, kdone.Destructor, error) {
	if b == nil {
		return reflect.Value{}, kdone.Noop, nil
	}
	if len(a) != 1 {
		return reflect.Value{}, nil, kerror.Newf(kerror.EViolation,
			"%s binder expects %d argument(s), %d given", b.t, 1, len(a))
	}
	if a[0].Type() != b.inType {
		return reflect.Value{}, nil, kerror.Newf(kerror.EViolation,
			"%s binder expects argument %d to be of %s type, %s given",
			b.t, 1, b.inType, a[0].Type())
	}
	return a[0].Convert(b.t), kdone.Noop, nil
}
