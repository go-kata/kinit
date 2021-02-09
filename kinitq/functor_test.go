package kinitq

import (
	"reflect"

	"github.com/go-kata/kinit"
)

type testFunctor struct {
	in []reflect.Type
}

func newTestFunctor(x interface{}) *testFunctor {
	ft := reflect.TypeOf(x)
	f := &testFunctor{}
	numIn := ft.NumIn()
	f.in = make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		f.in[i] = ft.In(i)
	}
	return f
}

func (f *testFunctor) Parameters() []reflect.Type {
	return f.in
}

func (f *testFunctor) Call(a ...reflect.Value) ([]kinit.Functor, error) {
	return nil, nil
}
