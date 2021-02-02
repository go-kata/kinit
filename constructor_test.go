package kinit

import (
	"reflect"

	"github.com/go-kata/kdone"
)

// func(...) (T, kdone.Destructor, error)

type testConstructor struct {
	t  reflect.Type
	f  reflect.Value
	in []reflect.Type
}

func newTestConstructor(x interface{}) *testConstructor {
	ft := reflect.TypeOf(x)
	c := &testConstructor{
		t: ft.Out(0),
		f: reflect.ValueOf(x),
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
	out := c.f.Call(a)
	obj := out[0]
	var dtor kdone.Destructor = kdone.Noop
	if v := out[1].Interface(); v != nil {
		dtor = v.(kdone.Destructor)
	}
	var err error
	if v := out[2].Interface(); v != nil {
		err = v.(error)
	}
	return obj, dtor, err
}

type testBrokenConstructor struct{}

func (testBrokenConstructor) Type() reflect.Type {
	return nil
}

func (testBrokenConstructor) Parameters() []reflect.Type {
	return nil
}

func (testBrokenConstructor) Create(a ...reflect.Value) (reflect.Value, kdone.Destructor, error) {
	return reflect.Value{}, nil, nil
}
