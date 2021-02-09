package kinitx

import (
	"reflect"
	"testing"

	"github.com/go-kata/kerror"
	"github.com/go-kata/kinit"
)

func TestFunctor__FunctionReturningNothing(t *testing.T) {
	var c int
	fun := MustNewFunctor(func(v *int) { *v++ })
	t.Logf("%+v", fun.Parameters())
	if _, err := fun.Call(reflect.ValueOf(&c)); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if c != 1 {
		t.Fail()
		return
	}
}

func TestFunctor__FunctionReturningError(t *testing.T) {
	var c int
	fun := MustNewFunctor(func(v *int) error {
		*v++
		return nil
	})
	t.Logf("%+v", fun.Parameters())
	if _, err := fun.Call(reflect.ValueOf(&c)); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if c != 1 {
		t.Fail()
		return
	}
}

func TestFunctor__FunctionReturningSingleFurtherAndError(t *testing.T) {
	var c int
	fun := MustNewFunctor(func(v *int) (kinit.Functor, error) {
		*v++
		return MustNewFunctor(func() {}), nil
	})
	t.Logf("%+v", fun.Parameters())
	if _, err := fun.Call(reflect.ValueOf(&c)); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if c != 1 {
		t.Fail()
		return
	}
}

func TestFunctor__FunctionReturningMultipleFurtherAndError(t *testing.T) {
	var c int
	fun := MustNewFunctor(func(v *int) ([]kinit.Functor, error) {
		*v++
		return []kinit.Functor{MustNewFunctor(func() {})}, nil
	})
	t.Logf("%+v", fun.Parameters())
	if _, err := fun.Call(reflect.ValueOf(&c)); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if c != 1 {
		t.Fail()
		return
	}
}

func TestNewFunctor__Nil(t *testing.T) {
	_, err := NewFunctor(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewFunctor__NilFunction(t *testing.T) {
	_, err := NewFunctor((func())(nil))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewFunctor__WrongType(t *testing.T) {
	_, err := NewFunctor(0)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewFunctor__WrongSignature(t *testing.T) {
	_, err := NewFunctor(func() int { return 0 })
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestFunctor_Call__WrongNumberOfArguments(t *testing.T) {
	var c int
	fun := MustNewFunctor(func(v *int) { *v++ })
	t.Logf("%+v", fun.Parameters())
	_, err := fun.Call(reflect.ValueOf(&c), reflect.ValueOf(0))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestFunctor_Call__WrongArgumentType(t *testing.T) {
	fun := MustNewFunctor(func(v *int) { *v++ })
	t.Logf("%+v", fun.Parameters())
	_, err := fun.Call(reflect.ValueOf(""))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNilFunctor_Parameters(t *testing.T) {
	if (*Functor)(nil).Parameters() != nil {
		t.Fail()
		return
	}
}

func TestNilFunctor_Call(t *testing.T) {
	further, err := (*Functor)(nil).Call()
	if further != nil {
		t.Fail()
		return
	}
	if err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
}
