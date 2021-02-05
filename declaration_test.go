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
	if err := d.Fulfill(); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
	if c != 1 {
		t.Fail()
		return
	}
}

func TestDeclaration__Error(t *testing.T) {
	d := NewDeclaration()
	var c int
	d.MustDeclareErrorProne(func() error {
		c++
		return kerror.New(nil, "test error")
	})
	d.MustDeclare(func() { c++ })
	err := d.Fulfill()
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

func TestDeclaration_Declare__NilFunction(t *testing.T) {
	d := NewDeclaration()
	err := d.Declare(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestDeclaration_Declare__Fulfilled(t *testing.T) {
	d := NewDeclaration()
	d.MustFulfill()
	err := d.Declare(func() {})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EIllegal {
		t.Fail()
		return
	}
}

func TestDeclaration_DeclareErrorProne__NilFunction(t *testing.T) {
	d := NewDeclaration()
	err := d.DeclareErrorProne(nil)
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EInvalid {
		t.Fail()
		return
	}
}

func TestDeclaration_DeclareErrorProne__Fulfilled(t *testing.T) {
	d := NewDeclaration()
	d.MustFulfill()
	err := d.DeclareErrorProne(func() error { return nil })
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EIllegal {
		t.Fail()
		return
	}
}

func TestDeclaration_Fulfill__Fulfilled(t *testing.T) {
	d := NewDeclaration()
	d.MustFulfill()
	err := d.Fulfill()
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.EIllegal {
		t.Fail()
		return
	}
}

func TestDeclaration_Fulfilled(t *testing.T) {
	d := NewDeclaration()
	d.MustFulfill()
	if !d.Fulfilled() {
		t.Fail()
		return
	}
}

func TestNilDeclaration_Declare(t *testing.T) {
	err := (*Declaration)(nil).Declare(func() {})
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENil {
		t.Fail()
		return
	}
}

func TestNilDeclaration_DeclareErrorProne(t *testing.T) {
	err := (*Declaration)(nil).DeclareErrorProne(func() error { return nil })
	t.Logf("%+v", err)
	if kerror.ClassOf(err) != kerror.ENil {
		t.Fail()
		return
	}
}

func TestNilDeclaration_Fulfill(t *testing.T) {
	if err := (*Declaration)(nil).Fulfill(); err != nil {
		t.Logf("%+v", err)
		t.Fail()
		return
	}
}

func TestNilDeclaration_Fulfilled(t *testing.T) {
	if (*Declaration)(nil).Fulfilled() {
		t.Fail()
		return
	}
}
