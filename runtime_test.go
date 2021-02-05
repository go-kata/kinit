package kinit

import (
	"reflect"
	"testing"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

func TestRuntime(t *testing.T) {
	var c int
	ctr := NewContainer()
	ctr.MustProvide(newTestConstructor(func() (int32, kdone.Destructor, error) {
		c++
		return int32(c), kdone.Noop, nil
	}))
	ctr.MustProvide(newTestConstructor(func() (int64, kdone.Destructor, error) {
		c++
		return int64(c), kdone.Noop, nil
	}))
	arena := NewArena()
	defer arena.MustFinalize()
	runtime := MustNewRuntime(ctr, arena)
	runtime.MustRegister(reflect.TypeOf(t), reflect.ValueOf(t), kdone.Noop)
	runtime.MustRun(newTestFunctor(func(innerRuntime1 *Runtime, i32 int32) ([]Functor, error) {
		if i32 != 1 {
			return nil, kerror.Newf(kerror.EInvalid, "int32: %d expected, %d given", 1, i32)
		}
		innerRuntime1.MustRegister(reflect.TypeOf(t), reflect.ValueOf(t), kdone.Noop)
		innerRuntime1.MustRun(newTestFunctor(func(innerRuntime2 *Runtime, i32 int32, i64 int64) ([]Functor, error) {
			if i32 != 1 {
				return nil, kerror.Newf(kerror.EInvalid, "int32: %d expected, %d given", 1, i32)
			}
			if i64 != 2 {
				return nil, kerror.Newf(kerror.EInvalid, "int64: %d expected, %d given", 2, i64)
			}
			innerRuntime2.MustRegister(reflect.TypeOf(t), reflect.ValueOf(t), kdone.Noop)
			return nil, nil
		}))
		return nil, nil
	}))
}

func TestNewRuntime__NilContainer(t *testing.T) {
	_, err := NewRuntime(nil, NewArena())
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestNewRuntime__NilArena(t *testing.T) {
	_, err := NewRuntime(NewContainer(), nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestNilRuntime_Register(t *testing.T) {
	x := 1
	err := (*Runtime)(nil).Register(reflect.TypeOf(x), reflect.ValueOf(x), kdone.Noop)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENil {
		t.Fail()
		return
	}
}

func TestNilRuntime_Run(t *testing.T) {
	err := (*Runtime)(nil).Run()
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENil {
		t.Fail()
		return
	}
}
