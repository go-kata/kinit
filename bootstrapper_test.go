package kinit

import (
	"reflect"

	"github.com/go-kata/kdone"
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
