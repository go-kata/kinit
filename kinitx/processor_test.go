package kinitx

import (
	"reflect"
	"testing"

	"github.com/go-kata/kerror"
)

func TestProcessor__FunctionReturningNothing(t *testing.T) {
	var c int
	proc := MustNewProcessor(func(v *int, i int8) { *v += int(i) })
	t.Logf("%+v %+v", proc.Type(), proc.Parameters())
	if err := proc.Process(reflect.ValueOf(&c), reflect.ValueOf(int8(1))); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if c != 1 {
		t.Fail()
		return
	}
}

func TestProcessor__FunctionReturningError(t *testing.T) {
	var c int
	proc := MustNewProcessor(func(v *int, i int8) error {
		*v += int(i)
		return nil
	})
	t.Logf("%+v %+v", proc.Type(), proc.Parameters())
	if err := proc.Process(reflect.ValueOf(&c), reflect.ValueOf(int8(1))); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if c != 1 {
		t.Fail()
		return
	}
}

func TestNewProcessor__Nil(t *testing.T) {
	_, err := NewProcessor(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewProcessor__NilFunction(t *testing.T) {
	_, err := NewProcessor((func())(nil))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewProcessor__WrongType(t *testing.T) {
	_, err := NewProcessor(0)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNewProcessor__WrongSignature(t *testing.T) {
	_, err := NewProcessor(func() {})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestProcessor_Process__WrongObjectType(t *testing.T) {
	proc := MustNewProcessor(func(v *int, i int8) { *v += int(i) })
	t.Logf("%+v %+v", proc.Type(), proc.Parameters())
	err := proc.Process(reflect.ValueOf(""), reflect.ValueOf(int8(1)))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestProcessor_Process__WrongNumberOfArguments(t *testing.T) {
	var c int
	proc := MustNewProcessor(func(v *int, i int8) { *v += int(i) })
	t.Logf("%+v %+v", proc.Type(), proc.Parameters())
	err := proc.Process(reflect.ValueOf(&c))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestProcessor_Process__WrongArgumentType(t *testing.T) {
	var c int
	proc := MustNewProcessor(func(v *int, i int8) { *v += int(i) })
	t.Logf("%+v %+v", proc.Type(), proc.Parameters())
	err := proc.Process(reflect.ValueOf(&c), reflect.ValueOf(""))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNilProcessor_Type(t *testing.T) {
	if (*Processor)(nil).Type() != nil {
		t.Fail()
		return
	}
}

func TestNilProcessor_Parameters(t *testing.T) {
	if (*Processor)(nil).Parameters() != nil {
		t.Fail()
		return
	}
}

func TestNilProcessor_Process(t *testing.T) {
	if err := (*Processor)(nil).Process(reflect.Value{}); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
}
