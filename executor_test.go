package kinit

import "reflect"

// func(...) (Executor, error)

type testExecutor struct {
	f  reflect.Value
	in []reflect.Type
}

func newTestExecutor(x interface{}) *testExecutor {
	ft := reflect.TypeOf(x)
	e := &testExecutor{
		f: reflect.ValueOf(x),
	}
	numIn := ft.NumIn()
	e.in = make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		e.in[i] = ft.In(i)
	}
	return e
}

func (e *testExecutor) Parameters() []reflect.Type {
	return e.in
}

func (e *testExecutor) Execute(a ...reflect.Value) (Executor, error) {
	out := e.f.Call(a)
	var exec Executor
	if v := out[0].Interface(); v != nil {
		exec = v.(Executor)
	}
	var err error
	if v := out[1].Interface(); v != nil {
		err = v.(error)
	}
	return exec, err
}
