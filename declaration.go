package kinit

import "github.com/go-kata/kerror"

// Declaration represents a registry of declared functions.
type Declaration struct {
	// functions specifies declared functions.
	functions []func() error
	// performed specifies whether were the declared functions called.
	performed bool
}

// NewDeclaration returns a new declaration.
func NewDeclaration() *Declaration {
	return &Declaration{}
}

// Declare appends the given function to this declaration.
func (d *Declaration) Declare(f func()) error {
	var fe func() error
	if f != nil {
		fe = func() error {
			f()
			return nil
		}
	}
	return d.DeclareErrorProne(fe)
}

// MustDeclare is a variant of the Declare that panics on error.
func (d *Declaration) MustDeclare(f func()) {
	if err := d.Declare(f); err != nil {
		panic(err)
	}
}

// Declare appends the given function to this declaration.
func (d *Declaration) DeclareErrorProne(f func() error) error {
	if d == nil {
		kerror.NPE()
		return nil
	}
	if d.performed {
		return kerror.New(kerror.EIllegal, "declared functions already called")
	}
	if f == nil {
		return kerror.New(kerror.EInvalid, "nil function cannot be declared")
	}
	d.functions = append(d.functions, f)
	return nil
}

// MustDeclareErrorProne is a variant of the MustDeclareErrorProne that panics on error.
func (d *Declaration) MustDeclareErrorProne(f func() error) {
	if err := d.DeclareErrorProne(f); err != nil {
		panic(err)
	}
}

// Perform calls declared functions.
func (d *Declaration) Perform() error {
	if d == nil {
		return nil
	}
	if d.performed {
		return kerror.New(kerror.EIllegal, "declared functions already called")
	}
	d.performed = true
	for _, f := range d.functions {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

// MustPerform is a variant of the Perform that panics on error.
func (d *Declaration) MustPerform() {
	if err := d.Perform(); err != nil {
		panic(err)
	}
}
