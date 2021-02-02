package kinit

import (
	"reflect"
	"testing"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

func TestArena(t *testing.T) {
	arena := NewArena()
	defer arena.MustFinalize()
	x := 1
	xt := reflect.TypeOf(x)
	arena.MustRegister(xt, reflect.ValueOf(x), kdone.Noop)
	if obj, ok := arena.Get(xt); !ok || obj.Interface() != x {
		t.Fail()
	}
}

func TestArena_RegisterWithNilObject(t *testing.T) {
	arena := NewArena()
	defer arena.MustFinalize()
	err := arena.Register(nil, reflect.Value{}, kdone.Noop)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
	}
}

func TestArena_RegisterWithNilDestructor(t *testing.T) {
	arena := NewArena()
	defer arena.MustFinalize()
	x := 1
	err := arena.Register(reflect.TypeOf(x), reflect.ValueOf(x), nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
	}
}

func TestArena_RegisterWithAmbiguousObject(t *testing.T) {
	arena := NewArena()
	defer arena.MustFinalize()
	x := 1
	xt := reflect.TypeOf(x)
	xv := reflect.ValueOf(x)
	arena.MustRegister(xt, xv, kdone.Noop)
	err := arena.Register(xt, xv, kdone.Noop)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EAmbiguous {
		t.Fail()
	}
}

func TestArena_GetWithNilType(t *testing.T) {
	arena := NewArena()
	defer arena.MustFinalize()
	x := 1
	arena.MustRegister(reflect.TypeOf(x), reflect.ValueOf(x), kdone.Noop)
	if _, ok := arena.Get(nil); ok {
		t.Fail()
		return
	}
}

func TestArena_Finalize(t *testing.T) {
	var c int
	defer func() {
		if c != -1 {
			t.Fail()
			return
		}
	}()
	arena := NewArena()
	defer arena.MustFinalize()
	x := 1
	arena.MustRegister(reflect.TypeOf(x), reflect.ValueOf(x), kdone.DestructorFunc(func() error {
		c -= x
		return nil
	}))
}

func TestNilArena_Register(t *testing.T) {
	defer func() {
		v := recover()
		t.Logf("%+v", v)
		if v == nil {
			t.Fail()
			return
		}
	}()
	_ = (*Arena)(nil).Register(nil, reflect.Value{}, nil)
}

func TestNilArena_Finalize(t *testing.T) {
	if err := (*Arena)(nil).Finalize(); err != nil {
		t.Logf("%+v", err)
		return
	}
}
