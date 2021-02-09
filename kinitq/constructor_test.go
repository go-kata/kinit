package kinitq

import (
	"reflect"

	"github.com/go-kata/kdone"
)

type testConstructor struct {
	t  reflect.Type
	in []reflect.Type
}

func newTestConstructor(x interface{}) *testConstructor {
	ft := reflect.TypeOf(x)
	c := &testConstructor{
		t: ft.Out(0),
	}
	numIn := ft.NumIn()
	c.in = make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		c.in[i] = ft.In(i)
	}
	return c
}

func (c *testConstructor) Type() reflect.Type {
	return c.t
}

func (c *testConstructor) Parameters() []reflect.Type {
	return c.in
}

func (c *testConstructor) Create(a ...reflect.Value) (reflect.Value, kdone.Destructor, error) {
	return reflect.Value{}, kdone.Noop, nil
}
