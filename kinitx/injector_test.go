package kinitx

import (
	"reflect"
	"testing"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
	"github.com/go-kata/kinit"
)

func TestInjector(t *testing.T) {
	x := 1
	fun := MustNewInjector(x)
	t.Logf("%+v", fun.Parameters())
	ctr := kinit.NewContainer()
	arena := kinit.NewArena()
	defer arena.MustFinalize()
	runtime := kinit.MustNewRuntime(ctr, arena)
	arena.MustPut(reflect.TypeOf(runtime), reflect.ValueOf(runtime), kdone.Noop)
	if _, err := fun.Call(reflect.ValueOf(runtime)); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if obj, ok := arena.Get(reflect.TypeOf(x)); !ok || obj.Interface() != x {
		t.Fail()
		return
	}
}

func TestNewInjector__Nil(t *testing.T) {
	_, err := NewInjector(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestInjector_Call__NilRuntime(t *testing.T) {
	x := 1
	fun := MustNewInjector(x)
	t.Logf("%+v", fun.Parameters())
	_, err := fun.Call(reflect.ValueOf((*kinit.Runtime)(nil)))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENil {
		t.Fail()
		return
	}
}

func TestInjector_Call__WrongNumberOfArguments(t *testing.T) {
	fun := MustNewInjector(1)
	t.Logf("%+v", fun.Parameters())
	_, err := fun.Call()
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestInjector_Call__WrongArgumentType(t *testing.T) {
	fun := MustNewInjector(1)
	t.Logf("%+v", fun.Parameters())
	_, err := fun.Call(reflect.ValueOf(""))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestNilInjector_Parameters(t *testing.T) {
	if (*Injector)(nil).Parameters() != nil {
		t.Fail()
		return
	}
}

func TestNilInjector_Call(t *testing.T) {
	further, err := (*Injector)(nil).Call()
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
