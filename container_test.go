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
			return nil, runtime.Register(reflect.TypeOf(t), reflect.ValueOf(t), kdone.Noop)
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

func TestNilContainer_Run(t *testing.T) {
	err := (*Container)(nil).Run()
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENil {
		t.Fail()
		return
	}
}
