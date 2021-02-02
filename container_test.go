package kinit

import (
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
	c := NewContainer()
	c.MustProvide(newTestConstructor(func() (*int, kdone.Destructor, error) { return &counter, kdone.Noop, nil }))
	c.MustApply(newTestProcessor(processTestCounter))
	c.MustProvide(newTestConstructor(newTestObject1))
	c.MustProvide(newTestConstructor(newTestObject2))
	c.MustInvoke(
		newTestExecutor(func(*testObject1) (Executor, error) {
			if counter != 2 {
				return nil, kerror.Newf(nil, "counter must be 2, %d found", counter)
			}
			return newTestExecutor(func(*testObject2) (Executor, error) {
				if counter != 0 {
					return nil, kerror.Newf(nil, "counter must be 0, %d found", counter)
				}
				return nil, nil
			}), nil
		}),
		newTestBootstrapper(t),
	)
}

func TestContainer_ProvideWithNilConstructor(t *testing.T) {
	c := NewContainer()
	err := c.Provide(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_ProvideWithBrokenConstructor(t *testing.T) {
	c := NewContainer()
	err := c.Provide(testBrokenConstructor{})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_ProvideWithAmbiguousConstructor(t *testing.T) {
	c := NewContainer()
	c.MustProvide(newTestConstructor(newTestObject1))
	err := c.Provide(newTestConstructor(newTestObject1))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EAmbiguous {
		t.Fail()
		return
	}
}

func TestContainer_ApplyWithNilProcessor(t *testing.T) {
	c := NewContainer()
	err := c.Apply(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_ApplyWithBrokenProcessor(t *testing.T) {
	c := NewContainer()
	err := c.Apply(testBrokenProcessor{})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_InvokeWithNilExecutor(t *testing.T) {
	c := NewContainer()
	err := c.Invoke(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_InvokeWithNilBootstrapper(t *testing.T) {
	c := NewContainer()
	err := c.Invoke(newTestExecutor(func() (Executor, error) { return nil, nil }), nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_InvokeWithBrokenGraph(t *testing.T) {
	c := NewContainer()
	c.MustProvide(newTestConstructor(newTestObject1))
	err := c.Invoke(
		newTestExecutor(func(*testObject1) (Executor, error) {
			return nil, nil
		}),
		newTestBootstrapper(t),
	)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENotFound {
		t.Fail()
		return
	}
}

func TestNilContainer_Provide(t *testing.T) {
	defer func() {
		v := recover()
		t.Logf("%+v", v)
		if v == nil {
			t.Fail()
			return
		}
	}()
	_ = (*Container)(nil).Provide(newTestConstructor(newTestObject1))
}

func TestNilContainer_Apply(t *testing.T) {
	defer func() {
		v := recover()
		t.Logf("%+v", v)
		if v == nil {
			t.Fail()
			return
		}
	}()
	_ = (*Container)(nil).Apply(newTestProcessor(processTestCounter))
}

func TestNilContainer_Invoke(t *testing.T) {
	defer func() {
		v := recover()
		t.Logf("%+v", v)
		if v == nil {
			t.Fail()
			return
		}
	}()
	_ = (*Container)(nil).Invoke(newTestExecutor(func() (Executor, error) { return nil, nil }))
}
