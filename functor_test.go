package kinit

import "reflect"

// func(...) ([]Functor, error)

type testFunctor struct {
	f  reflect.Value
	in []reflect.Type
}

func newTestFunctor(x interface{}) *testFunctor {
	ft := reflect.TypeOf(x)
	r := &testFunctor{
		f: reflect.ValueOf(x),
	}
	numIn := ft.NumIn()
	r.in = make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		r.in[i] = ft.In(i)
	}
	return r
}

func (f *testFunctor) Parameters() []reflect.Type {
	return f.in
}

func (f *testFunctor) Call(a ...reflect.Value) ([]Functor, error) {
	out := f.f.Call(a)
	var further []Functor
	if v := out[0].Interface(); v != nil {
		further = v.([]Functor)
	}
	var err error
	if v := out[1].Interface(); v != nil {
		err = v.(error)
	}
	return further, err
}

type testFunctorWithBrokenParameters struct{}

func (testFunctorWithBrokenParameters) Parameters() []reflect.Type {
	return []reflect.Type{nil}
}

func (testFunctorWithBrokenParameters) Call(a ...reflect.Value) ([]Functor, error) {
	return nil, nil
}
