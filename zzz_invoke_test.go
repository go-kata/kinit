package kinit

import (
	"testing"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

func TestContainer_Invoke(t *testing.T) {
	c := 0
	defer func() {
		if c != 1 {
			t.Fail()
			return
		}
	}()
	ctr := NewContainer()
	ctr.MustProvide(newTestConstructor(func() (*int, kdone.Destructor, error) { return &c, kdone.Noop, nil }))
	ctr.MustApply(newTestProcessor(processTestCounter))
	ctr.MustProvide(newTestConstructor(newTestObject1))
	ctr.MustProvide(newTestConstructor(newTestObject2))
	ctr.MustInvoke(
		newTestExecutor(func(*testObject1) (Executor, error) {
			if c != 2 {
				return nil, kerror.Newf(nil, "counter must be 2, %d found", c)
			}
			return newTestExecutor(func(*testObject2) (Executor, error) {
				if c != 0 {
					return nil, kerror.Newf(nil, "counter must be 0, %d found", c)
				}
				return nil, nil
			}), nil
		}),
		newTestBootstrapper(t),
	)
}

func TestContainer_Invoke__NilExecutor(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Invoke(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_Invoke__NilBootstrapper(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Invoke(newTestExecutor(func() (Executor, error) { return nil, nil }), nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestContainer_Invoke__BrokenGraph(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Invoke(
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

func TestContainer_Invoke__ErrorProneExecutor(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Invoke(
		newTestExecutor(func() (Executor, error) {
			return nil, kerror.New(kerror.ECustom, "test error")
		}),
	)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ECustom {
		t.Fail()
		return
	}
}

func TestContainer_Invoke__ErrorProneBootstrapper(t *testing.T) {
	ctr := NewContainer()
	err := ctr.Invoke(
		newTestExecutor(func() ([]Functor, error) {
			return nil, nil
		}),
		testErrorProneBootstrapper{},
	)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ECustom {
		t.Fail()
		return
	}
}

func TestNilContainer_Invoke(t *testing.T) {
	err := (*Container)(nil).Invoke(newTestExecutor(func() (Executor, error) { return nil, nil }))
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENil {
		t.Fail()
		return
	}
}
