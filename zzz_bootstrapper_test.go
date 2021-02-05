package kinit

import (
	"reflect"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

type testBootstrapper struct {
	t reflect.Type
	v reflect.Value
}

func newTestBootstrapper(x interface{}) *testBootstrapper {
	return &testBootstrapper{
		t: reflect.TypeOf(x),
		v: reflect.ValueOf(x),
	}
}

func (b *testBootstrapper) Bootstrap(arena *Arena) error {
	return arena.Register(b.t, b.v, kdone.Noop)
}

type testErrorProneBootstrapper struct{}

func (testErrorProneBootstrapper) Bootstrap(arena *Arena) error {
	return kerror.New(kerror.ECustom, "test error")
}
