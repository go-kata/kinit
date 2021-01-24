package kinit

import (
	"reflect"

	"github.com/go-kata/kdone"
)

type testInitializer struct {
	t reflect.Type
	v reflect.Value
}

func newTestInitializer(x interface{}) *testInitializer {
	return &testInitializer{
		t: reflect.TypeOf(x),
		v: reflect.ValueOf(x),
	}
}

func (i *testInitializer) Initialize(arena *Arena) error {
	return arena.Register(i.t, i.v, kdone.Noop)
}
