package kinit

import (
	"testing"

	"github.com/go-kata/kerror"
)

func TestDeclare__NilFunction(t *testing.T) {
	err := Declare(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestMustDeclare__NilFunction(t *testing.T) {
	err := kerror.Try(func() error {
		MustDeclare(nil)
		return nil
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestDeclareErrorProne__NilFunction(t *testing.T) {
	err := DeclareErrorProne(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestMustDeclareErrorProne__NilFunction(t *testing.T) {
	err := kerror.Try(func() error {
		MustDeclareErrorProne(nil)
		return nil
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestProvide__NilConstructor(t *testing.T) {
	err := Provide(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestMustProvide__NilConstructor(t *testing.T) {
	err := kerror.Try(func() error {
		MustProvide(nil)
		return nil
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestAttach__NilProcessor(t *testing.T) {
	err := Attach(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestMustAttach__NilProcessor(t *testing.T) {
	err := kerror.Try(func() error {
		MustAttach(nil)
		return nil
	})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}
