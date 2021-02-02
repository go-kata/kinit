package kinit

import "testing"

func TestHook(t *testing.T) {
	if i := Hook(func() error { return nil }); i == -1 {
		t.Fail()
		return
	}
}

func TestHookWithNil(t *testing.T) {
	if i := Hook(nil); i != -1 {
		t.Fail()
		return
	}
}

func TestMustHook(t *testing.T) {
	if i := MustHook(func() {}); i == -1 {
		t.Fail()
		return
	}
}

func TestMustHookWithNil(t *testing.T) {
	if i := MustHook(nil); i != -1 {
		t.Fail()
		return
	}
}

func TestHooksCalled(t *testing.T) {
	if err := callHooks(); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if !hooksCalled {
		t.Fail()
		return
	}
	if i := Hook(func() error { return nil }); i != -1 {
		t.Fail()
		return
	}
	if i := MustHook(func() {}); i != -1 {
		t.Fail()
		return
	}
}
