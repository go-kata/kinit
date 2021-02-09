package kinitx

import (
	"testing"

	"github.com/go-kata/kerror"
)

func TestProvide__Nil(t *testing.T) {
	err := Provide(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestMustProvide__Nil(t *testing.T) {
	err := kerror.Try(func() error {
		MustProvide(nil)
		return nil
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestBind__NilInterfacePointer(t *testing.T) {
	err := Bind(nil, 0)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestMustBind__NilInterfacePointer(t *testing.T) {
	err := kerror.Try(func() error {
		MustBind(nil, 0)
		return nil
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestAttach__Nil(t *testing.T) {
	err := Attach(nil)
	t.Logf("%+v", err)
	if err == nil {
		t.Fail()
		return
	}
}

func TestMustAttach__Nil(t *testing.T) {
	err := kerror.Try(func() error {
		MustAttach(nil)
		return nil
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestRun__Nil(t *testing.T) {
	err := Run(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestMustRun__Nil(t *testing.T) {
	err := kerror.Try(func() error {
		MustRun(nil)
		return nil
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestRequire__Nil(t *testing.T) {
	err := Require(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestMustRequire__Nil(t *testing.T) {
	err := kerror.Try(func() error {
		MustRequire(nil)
		return nil
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestIgnore__Nil(t *testing.T) {
	err := Ignore(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestMustIgnore__Nil(t *testing.T) {
	err := kerror.Try(func() error {
		MustIgnore(nil)
		return nil
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestConsider__Nil(t *testing.T) {
	err := Consider(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestMustConsider__Nil(t *testing.T) {
	err := kerror.Try(func() error {
		MustConsider(nil)
		return nil
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EViolation {
		t.Fail()
		return
	}
}

func TestInspect(t *testing.T) {
	if err := Inspect(nil); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
}

func TestMustInspect(t *testing.T) {
	if err := kerror.Try(func() error {
		MustInspect(nil)
		return nil
	}); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
}
