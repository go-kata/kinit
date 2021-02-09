package kinit

import (
	"reflect"
	"testing"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

func TestArena(t *testing.T) {
	parent1 := NewArena()
	defer parent1.MustFinalize()
	x := int16(1)
	xt := reflect.TypeOf(x)
	parent1.MustPut(xt, reflect.ValueOf(x), kdone.Noop)

	parent2 := NewArena()
	defer parent2.MustFinalize()
	y := int32(1)
	yt := reflect.TypeOf(y)
	parent2.MustPut(yt, reflect.ValueOf(y), kdone.Noop)

	arena := NewArena(parent1, parent2)
	defer arena.MustFinalize()
	z := int64(1)
	zt := reflect.TypeOf(z)
	arena.MustPut(zt, reflect.ValueOf(z), kdone.Noop)

	if obj, ok := arena.Get(xt); !ok || obj.Interface() != x {
		t.Fail()
		return
	}
	if obj, ok := arena.Get(yt); !ok || obj.Interface() != y {
		t.Fail()
		return
	}
	if obj, ok := arena.Get(zt); !ok || obj.Interface() != z {
		t.Fail()
		return
	}
}

func TestArena__SameType(t *testing.T) {
	parent := NewArena()
	defer parent.MustFinalize()
	x := 1
	xt := reflect.TypeOf(x)
	parent.MustPut(xt, reflect.ValueOf(x), kdone.Noop)

	arena := NewArena(parent)
	defer arena.MustFinalize()
	y := 2
	yt := reflect.TypeOf(y)
	arena.MustPut(yt, reflect.ValueOf(y), kdone.Noop)

	if obj, ok := arena.Get(yt); !ok || obj.Interface() != y {
		t.Fail()
		return
	}
}

func TestArena__NilParent(t *testing.T) {
	parent := NewArena()
	defer parent.MustFinalize()
	x := 1
	xt := reflect.TypeOf(x)
	parent.MustPut(xt, reflect.ValueOf(x), kdone.Noop)

	arena := NewArena(nil, parent)
	defer arena.MustFinalize()

	if obj, ok := arena.Get(xt); !ok || obj.Interface() != x {
		t.Fail()
		return
	}
}

func TestArena__FinalizedParent(t *testing.T) {
	parent := NewArena()
	x := int16(1)
	xt := reflect.TypeOf(x)
	parent.MustPut(xt, reflect.ValueOf(x), kdone.Noop)

	arena := NewArena(parent)
	defer arena.MustFinalize()
	y := int64(1)
	yt := reflect.TypeOf(y)
	arena.MustPut(yt, reflect.ValueOf(y), kdone.Noop)

	parent.MustFinalize()

	if obj, ok := arena.Get(xt); ok {
		t.Logf("%+v", obj)
		t.Fail()
		return
	}
	if obj, ok := arena.Get(yt); !ok || obj.Interface() != y {
		t.Fail()
		return
	}
}

func TestArena_Put__NilObject(t *testing.T) {
	arena := NewArena()
	defer arena.MustFinalize()
	err := arena.Put(nil, reflect.Value{}, kdone.Noop)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestArena_Put__NilDestructor(t *testing.T) {
	arena := NewArena()
	defer arena.MustFinalize()
	x := 1
	err := arena.Put(reflect.TypeOf(x), reflect.ValueOf(x), nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestArena_Put__AmbiguousObject(t *testing.T) {
	arena := NewArena()
	defer arena.MustFinalize()
	x := 1
	xt := reflect.TypeOf(x)
	xv := reflect.ValueOf(x)
	arena.MustPut(xt, xv, kdone.Noop)
	err := arena.Put(xt, xv, kdone.Noop)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EAmbiguous {
		t.Fail()
		return
	}
}

func TestArena_Put__Finalized(t *testing.T) {
	arena := NewArena()
	arena.MustFinalize()
	x := 1
	err := arena.Put(reflect.TypeOf(x), reflect.ValueOf(x), nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EIllegal {
		t.Fail()
		return
	}
}

func TestArena_Get__NilType(t *testing.T) {
	arena := NewArena()
	defer arena.MustFinalize()
	x := 1
	arena.MustPut(reflect.TypeOf(x), reflect.ValueOf(x), kdone.Noop)
	if obj, ok := arena.Get(nil); ok {
		t.Logf("%+v", obj)
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
	arena.MustPut(reflect.TypeOf(x), reflect.ValueOf(x), kdone.DestructorFunc(func() error {
		c -= x
		return nil
	}))
}

func TestArena_Finalize__Finalized(t *testing.T) {
	arena := NewArena()
	arena.MustFinalize()
	err := arena.Finalize()
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EIllegal {
		t.Fail()
		return
	}
}

func TestArena_Finalized(t *testing.T) {
	arena := NewArena()
	arena.MustFinalize()
	if !arena.Finalized() {
		t.Fail()
		return
	}
}

func TestNilArena_Put(t *testing.T) {
	err := (*Arena)(nil).Put(nil, reflect.Value{}, nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENil {
		t.Fail()
		return
	}
}

func TestNilArena_Get(t *testing.T) {
	x := 1
	if obj, ok := (*Arena)(nil).Get(reflect.TypeOf(x)); ok {
		t.Logf("%+v", obj)
		t.Fail()
		return
	}
}

func TestNilArena_Finalize(t *testing.T) {
	if err := (*Arena)(nil).Finalize(); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
}

func TestNilArena_Finalized(t *testing.T) {
	if (*Arena)(nil).Finalized() {
		t.Fail()
		return
	}
}
