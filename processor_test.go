package kinit

import "reflect"

// func(T, ...) error

type testProcessor struct {
	t  reflect.Type
	f  reflect.Value
	in []reflect.Type
}

func newTestProcessor(x interface{}) *testProcessor {
	ft := reflect.TypeOf(x)
	p := &testProcessor{
		t: ft.In(0),
		f: reflect.ValueOf(x),
	}
	numIn := ft.NumIn()
	p.in = make([]reflect.Type, numIn-1)
	for i := 1; i < numIn; i++ {
		p.in[i-1] = ft.In(i)
	}
	return p
}

func (p *testProcessor) Type() reflect.Type {
	return p.t
}

func (p *testProcessor) Parameters() []reflect.Type {
	return p.in
}

func (p *testProcessor) Process(obj reflect.Value, a ...reflect.Value) error {
	out := p.f.Call(append([]reflect.Value{obj}, a...))
	var err error
	if v := out[0].Interface(); v != nil {
		err = v.(error)
	}
	return err
}
