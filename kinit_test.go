package kinit

import "testing"

func TestDeclareWithNilFunction(t *testing.T) {
	err := Declare(nil)
	t.Logf("%+v", err)
	if err == nil {
		t.Fail()
		return
	}
}

func TestMustDeclareWithNilFunction(t *testing.T) {
	defer func() {
		v := recover()
		t.Logf("%+v", v)
		if v == nil {
			t.Fail()
			return
		}
	}()
	MustDeclare(nil)
}

func TestDeclareErrorProneWithNilFunction(t *testing.T) {
	err := DeclareErrorProne(nil)
	t.Logf("%+v", err)
	if err == nil {
		t.Fail()
		return
	}
}

func TestMustDeclareErrorProneWithNilFunction(t *testing.T) {
	defer func() {
		v := recover()
		t.Logf("%+v", v)
		if v == nil {
			t.Fail()
			return
		}
	}()
	MustDeclareErrorProne(nil)
}

func TestProvideWithNilConstructor(t *testing.T) {
	err := Provide(nil)
	t.Logf("%+v", err)
	if err == nil {
		t.Fail()
		return
	}
}

func TestMustProvideWithNilConstructor(t *testing.T) {
	defer func() {
		v := recover()
		t.Logf("%+v", v)
		if v == nil {
			t.Fail()
			return
		}
	}()
	MustProvide(nil)
}

func TestApplyWithNilConstructor(t *testing.T) {
	err := Apply(nil)
	t.Logf("%+v", err)
	if err == nil {
		t.Fail()
		return
	}
}

func TestMustApplyWithNilProcessor(t *testing.T) {
	defer func() {
		v := recover()
		t.Logf("%+v", v)
		if v == nil {
			t.Fail()
			return
		}
	}()
	MustApply(nil)
}

func TestInvokeWithNilExecutor(t *testing.T) {
	err := Invoke(nil)
	t.Logf("%+v", err)
	if err == nil {
		t.Fail()
		return
	}
}

func TestInvokeWithNilBootstrapper(t *testing.T) {
	err := Invoke(newTestExecutor(func() (Executor, error) { return nil, nil }), nil)
	t.Logf("%+v", err)
	if err == nil {
		t.Fail()
		return
	}
}

func TestMustInvokeWithNilExecutor(t *testing.T) {
	defer func() {
		v := recover()
		t.Logf("%+v", v)
		if v == nil {
			t.Fail()
			return
		}
	}()
	MustInvoke(nil)
}

func TestMustInvokeWithNilBootstrapper(t *testing.T) {
	defer func() {
		v := recover()
		t.Logf("%+v", v)
		if v == nil {
			t.Fail()
			return
		}
	}()
	MustInvoke(newTestExecutor(func() (Executor, error) { return nil, nil }), nil)
}
