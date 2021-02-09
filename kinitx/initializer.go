package kinitx

import (
	"reflect"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

// Initializer represents a constructor based on a struct.
type Initializer struct {
	// t specifies the type of an object that is created by this initializer.
	t reflect.Type
	// assignableFieldTypes specifies types of assignable struct fields.
	assignableFieldTypes []reflect.Type
	// assignableFieldIndexes specifies indexes of assignable struct fields.
	assignableFieldIndexes []int
}

// NewInitializer returns a new initializer.
//
// The argument x must be a struct or a struct pointer.
func NewInitializer(x interface{}) (*Initializer, error) {
	if x == nil {
		return nil, kerror.New(kerror.EViolation, "struct or struct pointer expected, nil given")
	}
	t := reflect.TypeOf(x)
	var st reflect.Type
	switch t.Kind() {
	default:
		return nil, kerror.Newf(kerror.EViolation, "struct or struct pointer expected, %s given", t)
	case reflect.Struct:
		st = t
	case reflect.Ptr:
		st = t.Elem()
		if st.Kind() != reflect.Struct {
			return nil, kerror.Newf(kerror.EViolation, "struct or struct pointer expected, %s given", t)
		}
	}
	i := &Initializer{
		t: t,
	}
	for j, n := 0, st.NumField(); j < n; j++ {
		sf := st.Field(j)
		if sf.PkgPath != "" {
			continue
		}
		i.assignableFieldTypes = append(i.assignableFieldTypes, sf.Type)
		i.assignableFieldIndexes = append(i.assignableFieldIndexes, j)
	}
	return i, nil
}

// MustNewInitializer is a variant of the NewInitializer that panics on error.
func MustNewInitializer(x interface{}) *Initializer {
	i, err := NewInitializer(x)
	if err != nil {
		panic(err)
	}
	return i
}

// Type implements the kinit.Constructor interface.
func (i *Initializer) Type() reflect.Type {
	if i == nil {
		return nil
	}
	return i.t
}

// Parameters implements the kinit.Constructor interface.
func (i *Initializer) Parameters() []reflect.Type {
	if i == nil {
		return nil
	}
	types := make([]reflect.Type, len(i.assignableFieldTypes))
	copy(types, i.assignableFieldTypes)
	return types
}

// Create implements the kinit.Constructor interface.
func (i *Initializer) Create(a ...reflect.Value) (reflect.Value, kdone.Destructor, error) {
	if i == nil {
		return reflect.Value{}, kdone.Noop, nil
	}
	if len(a) != len(i.assignableFieldTypes) {
		return reflect.Value{}, nil, kerror.Newf(kerror.EViolation,
			"%s initializer expects %d argument(s), %d given",
			i.t, len(i.assignableFieldTypes), len(a))
	}
	for j, v := range a {
		if v.Type() != i.assignableFieldTypes[j] {
			return reflect.Value{}, nil, kerror.Newf(kerror.EViolation,
				"%s initializer expects argument %d to be of %s type, %s given",
				i.t, j+1, i.assignableFieldTypes[j], v.Type())
		}
	}
	var sp, obj reflect.Value
	if i.t.Kind() == reflect.Ptr {
		sp = reflect.New(i.t.Elem())
		obj = sp
	} else {
		sp = reflect.New(i.t)
		obj = sp.Elem()
	}
	sv := sp.Elem()
	for j, v := range a {
		sv.Field(i.assignableFieldIndexes[j]).Set(v)
	}
	return obj, kdone.Noop, nil
}
