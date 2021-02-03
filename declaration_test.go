package kinit

import (
	"testing"

	"github.com/go-kata/kerror"
)

func TestDeclaration(t *testing.T) {
	d := NewDeclaration()
	var c int
	if err := d.Declare(func() { c++ }); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if err := d.Perform(); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if c != 1 {
		t.Fail()
		return
	}
}

func TestDeclarationWithError(t *testing.T) {
	d := NewDeclaration()
	var c int
	d.MustDeclareErrorProne(func() error {
		c++
		return kerror.New(nil, "test error")
	})
	d.MustDeclare(func() { c++ })
	err := d.Perform()
	t.Logf("%+v", err)
	if err == nil {
		t.Fail()
		return
	}
	if c != 1 {
		t.Fail()
		return
	}
}

func TestDeclaration_DeclareWithNilFunction(t *testing.T) {
	d := NewDeclaration()
	err := d.Declare(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestDeclaration_DeclareErrorProneWithNilFunction(t *testing.T) {
	d := NewDeclaration()
	err := d.DeclareErrorProne(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestDeclaration_DeclareWhenPerformed(t *testing.T) {
	d := NewDeclaration()
	d.MustPerform()
	err := d.Declare(func() {})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EIllegal {
		t.Fail()
		return
	}
}

func TestDeclaration_PerformWhenPerformed(t *testing.T) {
	d := NewDeclaration()
	d.MustPerform()
	err := d.Perform()
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EIllegal {
		t.Fail()
		return
	}
}

func TestNilDeclaration_Declare(t *testing.T) {
	defer func() {
		v := recover()
		t.Logf("%+v", v)
		if v == nil {
			t.Fail()
			return
		}
	}()
	_ = (*Declaration)(nil).Declare(func() {})
}

func TestNilDeclaration_DeclareErrorProne(t *testing.T) {
	defer func() {
		v := recover()
		t.Logf("%+v", v)
		if v == nil {
			t.Fail()
			return
		}
	}()
	_ = (*Declaration)(nil).DeclareErrorProne(func() error { return nil })
}

func TestNilDeclaration_Perform(t *testing.T) {
	if err := (*Declaration)(nil).Perform(); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
}
