package kinit

import "github.com/go-kata/kerror"

// Declaration represents a registry of declared functions.
type Declaration struct {
	// functions specifies declared functions.
	functions []func() error
	// fulfilled specifies whether were the declared functions called.
	fulfilled bool
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

// DeclareErrorProne appends the given function to this declaration.
func (d *Declaration) DeclareErrorProne(f func() error) error {
	if d == nil {
		return kerror.New(kerror.ENil, "nil declaration cannot declare function")
	}
	if d.fulfilled {
		return kerror.New(kerror.EIllegal, "declaration has already called declared functions")
	}
	if f == nil {
		return kerror.New(kerror.EInvalid, "declaration cannot declare nil function")
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

// Fulfill calls declared functions.
func (d *Declaration) Fulfill() error {
	if d == nil {
		return nil
	}
	if d.fulfilled {
		return kerror.New(kerror.EIllegal, "declaration has already called declared functions")
	}
	defer func() {
		d.fulfilled = true
	}()
	for _, f := range d.functions {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

// MustFulfill is a variant of the Fulfill that panics on error.
func (d *Declaration) MustFulfill() {
	if err := d.Fulfill(); err != nil {
		panic(err)
	}
}

// Fulfilled returns boolean whether were the declared functions called.
func (d *Declaration) Fulfilled() bool {
	if d == nil {
		return false
	}
	return d.fulfilled
}
