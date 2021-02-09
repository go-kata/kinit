package kinit

import (
	"reflect"
	"testing"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

type testObject1 struct{}

func newTestObject1(c *int, t *testing.T) (*testObject1, kdone.Destructor, error) {
	t.Logf("counter before object #1 initialization: %d", *c)
	*c += 1
	t.Logf("counter after object #1 initialization: %d", *c)
	return &testObject1{}, kdone.DestructorFunc(func() error {
		t.Logf("counter before object #1 finalization: %d", *c)
		*c -= 1
		t.Logf("counter after object #1 finalization: %d", *c)
		return nil
	}), nil
}

type testObject2 struct{}

func newTestObject2(c *int, t *testing.T) (*testObject2, kdone.Destructor, error) {
	t.Logf("counter before object #2 initialization: %d", *c)
	*c -= 2
	t.Logf("counter after object #2 initialization: %d", *c)
	return &testObject2{}, kdone.DestructorFunc(func() error {
		t.Logf("counter before object #2 finalization: %d", *c)
		*c += 2
		t.Logf("counter after object #2 finalization: %d", *c)
		return nil
	}), nil
}

func processTestCounter(c *int, t *testing.T) error {
	t.Logf("counter before processing: %d", *c)
	*c++
	t.Logf("counter after processing: %d", *c)
	return nil
}

func TestContainer(t *testing.T) {
	counter := 0
	defer func() {
		if counter != 1 {
			t.Fail()
			return
		}
	}()
	ctr := NewContainer()
	ctr.MustProvide(newTestConstructor(func() (*int, kdone.Destructor, error) { return &counter, kdone.Noop, nil }))
	ctr.MustAttach(newTestProcessor(processTestCounter))
	ctr.MustProvide(newTestConstructor(newTestObject1))
	ctr.MustProvide(newTestConstructor(newTestObject2))
	ctr.MustRun(
		newTestFunctor(func(runtime *Runtime) ([]Functor, error) {
			return nil, runtime.Put(reflect.TypeOf(t), reflect.ValueOf(t), kdone.Noop)
		}),
		newTestFunctor(func(*testObject1) ([]Functor, error) {
			if counter != 2 {
				return nil, kerror.Newf(nil, "counter must be 2, %d found", counter)
			}
			return []Functor{newTestFunctor(func(*testObject2) ([]Functor, error) {
				if counter != 0 {
					return nil, kerror.Newf(nil, "counter must be 0, %d found", counter)
				}
				return nil, nil
			})}, nil
		}),
	)
}

func TestContainer_Provide__NilConstructor(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Provide(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_Provide__ConstructorWithBrokenType(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Provide(testConstructorWithBrokenType{})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_Provide__AmbiguousConstructor(t *testing.T) {
	ctr := NewContainer()
	ctr.MustProvide(newTestConstructor(newTestObject1))
	err := ctr.Provide(newTestConstructor(newTestObject1))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EAmbiguous {
		t.Fail()
		return
	}
}

func TestContainer_Attach__NilProcessor(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Attach(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_Attach__ProcessorWithBrokenType(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Attach(testProcessorWithBrokenType{})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_Lookup(t *testing.T) {
	ctor := newTestConstructor(func() (int, kdone.Destructor, error) { return 0, kdone.Noop, nil })
	proc1 := newTestProcessor(func(int) error { return nil })
	proc2 := newTestProcessor(func(int) error { return nil })
	ctr := NewContainer()
	ctr.MustProvide(ctor)
	ctr.MustAttach(proc1)
	ctr.MustAttach(proc2)
	if c, pp := ctr.Lookup(reflect.TypeOf(0)); c != ctor || len(pp) != 2 || pp[0] != proc1 || pp[1] != proc2 {
		t.Fail()
		return
	}
}

func TestContainer_Lookup__Nil(t *testing.T) {
	if ctor, processors := NewContainer().Lookup(nil); ctor != nil || len(processors) > 0 {
		t.Fail()
		return
	}
}

func TestContainer_Lookup__Unregistered(t *testing.T) {
	if ctor, processors := NewContainer().Lookup(reflect.TypeOf(0)); ctor != nil || len(processors) > 0 {
		t.Fail()
		return
	}
}

func TestContainer_Lookup__OnlyConstructor(t *testing.T) {
	ctor := newTestConstructor(func() (int, kdone.Destructor, error) { return 0, kdone.Noop, nil })
	ctr := NewContainer()
	ctr.MustProvide(ctor)
	if c, pp := ctr.Lookup(reflect.TypeOf(0)); c != ctor || len(pp) > 0 {
		t.Fail()
		return
	}
}

func TestContainer_Lookup__OnlyProcessors(t *testing.T) {
	proc1 := newTestProcessor(func(int) error { return nil })
	proc2 := newTestProcessor(func(int) error { return nil })
	ctr := NewContainer()
	ctr.MustAttach(proc1)
	ctr.MustAttach(proc2)
	if c, pp := ctr.Lookup(reflect.TypeOf(0)); c != nil || len(pp) != 2 || pp[0] != proc1 || pp[1] != proc2 {
		t.Fail()
		return
	}
}

func TestContainer_Explore(t *testing.T) {
	types := map[reflect.Type]struct {
		constructor Constructor
		processors  []Processor
	}{
		reflect.TypeOf(int16(0)): {
			constructor: newTestConstructor(func() (int16, kdone.Destructor, error) { return 0, kdone.Noop, nil }),
			processors: []Processor{
				newTestProcessor(func(int16) error { return nil }),
				newTestProcessor(func(int16) error { return nil }),
			},
		},
		reflect.TypeOf(int32(0)): {
			constructor: newTestConstructor(func() (int32, kdone.Destructor, error) { return 0, kdone.Noop, nil }),
		},
		reflect.TypeOf(int64(0)): {
			processors: []Processor{
				newTestProcessor(func(int64) error { return nil }),
				newTestProcessor(func(int64) error { return nil }),
			},
		},
	}
	ctr := NewContainer()
	for _, typ := range types {
		if typ.constructor != nil {
			ctr.MustProvide(typ.constructor)
		}
		for _, proc := range typ.processors {
			ctr.MustAttach(proc)
		}
	}
	ctr.Explore(func(rt reflect.Type, ctor Constructor, processors []Processor) (next bool) {
		typ, ok := types[rt]
		if !ok {
			t.Logf("%s", rt)
			t.Fail()
			return false
		}
		if ctor != typ.constructor {
			t.Fail()
			return false
		}
		for i := range typ.processors {
			if i > len(processors) || processors[i] != typ.processors[i] {
				t.Fail()
				return false
			}
		}
		return true
	})
}

func TestContainer_Explore__Nil(t *testing.T) {
	(*Container)(nil).Explore(nil)
}

func TestContainer_Explore__BreakOnlyConstructor(t *testing.T) {
	c := 0
	ctr := NewContainer()
	ctr.MustProvide(newTestConstructor(func() (int32, kdone.Destructor, error) { return 0, kdone.Noop, nil }))
	ctr.MustProvide(newTestConstructor(func() (int64, kdone.Destructor, error) { return 0, kdone.Noop, nil }))
	ctr.Explore(func(reflect.Type, Constructor, []Processor) (next bool) {
		c++
		return false
	})
	if c != 1 {
		t.Fail()
		return
	}
}

func TestContainer_Explore__BreakOnlyProcessors(t *testing.T) {
	c := 0
	ctr := NewContainer()
	ctr.MustAttach(newTestProcessor(func(int32) error { return nil }))
	ctr.MustAttach(newTestProcessor(func(int64) error { return nil }))
	ctr.Explore(func(reflect.Type, Constructor, []Processor) (next bool) {
		c++
		return false
	})
	if c != 1 {
		t.Fail()
		return
	}
}

func TestContainer_Run__NilFunctor(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Run(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_Run__ErrorProneFunctor(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Run(
		newTestFunctor(func() ([]Functor, error) {
			return nil, kerror.New(kerror.ECustom, "test error")
		}),
	)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ECustom {
		t.Fail()
		return
	}
}

func TestContainer_Run__NilSubsequentFunctor(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Run(
		newTestFunctor(func() ([]Functor, error) {
			return []Functor{nil}, nil
		}),
	)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_Run__ErrorProneSubsequentFunctor(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Run(
		newTestFunctor(func() ([]Functor, error) {
			return []Functor{newTestFunctor(func() ([]Functor, error) {
				return nil, kerror.New(kerror.ECustom, "test error")
			})}, nil
		}),
	)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ECustom {
		t.Fail()
		return
	}
}

func TestContainer_Run__ErrorProneConstructor(t *testing.T) {
	ctr := NewContainer()
	ctr.MustProvide(newTestConstructor(func() (int, kdone.Destructor, error) {
		return 0, kdone.Noop, kerror.New(kerror.ECustom, "test error")
	}))
	err := ctr.Run(
		newTestFunctor(func(int) ([]Functor, error) {
			return nil, nil
		}),
	)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ECustom {
		t.Fail()
		return
	}
}

func TestContainer_Run__ErrorProneProcessor(t *testing.T) {
	ctr := NewContainer()
	ctr.MustProvide(newTestConstructor(func() (int, kdone.Destructor, error) {
		return 1, kdone.Noop, nil
	}))
	ctr.MustAttach(newTestProcessor(func(int) error {
		return kerror.New(kerror.ECustom, "test error")
	}))
	err := ctr.Run(
		newTestFunctor(func(int) ([]Functor, error) {
			return nil, nil
		}),
	)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ECustom {
		t.Fail()
		return
	}
}

func TestContainer_Run__BrokenGraph(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Run(
		newTestFunctor(func(int) ([]Functor, error) {
			return nil, nil
		}),
	)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENotFound {
		t.Fail()
		return
	}
}

func TestContainer_Run__ConstructorWithBrokenParameters(t *testing.T) {
	ctr := NewContainer()
	ctr.MustProvide(testConstructorWithBrokenParameters{})
	err := ctr.Run(
		newTestFunctor(func(int) ([]Functor, error) {
			return nil, nil
		}),
	)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_Run__ConstructorWithBrokenDestructor(t *testing.T) {
	ctr := NewContainer()
	ctr.MustProvide(testConstructorWithBrokenDestructor{})
	err := ctr.Run(
		newTestFunctor(func(int) ([]Functor, error) {
			return nil, nil
		}),
	)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_Run__ProcessorWithBrokenParameters(t *testing.T) {
	ctr := NewContainer()
	ctr.MustProvide(newTestConstructor(func() (int, kdone.Destructor, error) {
		return 1, kdone.Noop, nil
	}))
	ctr.MustAttach(testProcessorWithBrokenParameters{})
	err := ctr.Run(
		newTestFunctor(func(int) ([]Functor, error) {
			return nil, nil
		}),
	)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_Run__FunctorWithBrokenParameters(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Run(testFunctorWithBrokenParameters{})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestNilContainer_Provide(t *testing.T) {
	err := (*Container)(nil).Provide(newTestConstructor(newTestObject1))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENil {
		t.Fail()
		return
	}
}

func TestNilContainer_Attach(t *testing.T) {
	err := (*Container)(nil).Attach(newTestProcessor(processTestCounter))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENil {
		t.Fail()
		return
	}
}

func TestNilContainer_Lookup(t *testing.T) {
	if ctor, processors := (*Container)(nil).Lookup(reflect.TypeOf(0)); ctor != nil || len(processors) > 0 {
		t.Fail()
		return
	}
}

func TestNilContainer_Explore(t *testing.T) {
	c := 0
	(*Container)(nil).Explore(func(reflect.Type, Constructor, []Processor) (next bool) {
		c++
		return true
	})
	if c != 0 {
		t.Fail()
		return
	}
}

func TestNilContainer_Run(t *testing.T) {
	err := (*Container)(nil).Run()
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENil {
		t.Fail()
		return
	}
}
