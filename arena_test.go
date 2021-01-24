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

func TestArenaWithNilObject(t *testing.T) {
	arena := NewArena()
	defer arena.MustFinalize()
	err := arena.Register(nil, reflect.Value{}, kdone.Noop)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
	}
}

func TestArenaWithAmbiguousObject(t *testing.T) {
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

func TestArenaFinalization(t *testing.T) {
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
