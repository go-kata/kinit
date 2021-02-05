// Package kinit provides tools for creating objects.
package kinit

// globalDeclaration specifies the global declaration.
var globalDeclaration = NewDeclaration()

// Declare declares the given function to call at the first global invocation.
func Declare(f func()) error {
	return globalDeclaration.Declare(f)
}

// MustDeclare is a variant of the Declare that panics on error.
func MustDeclare(f func()) struct{} {
	if err := Declare(f); err != nil {
		panic(err)
	}
	return struct{}{}
}

// DeclareErrorProne declares the given error-prone function to call at the first global invocation.
func DeclareErrorProne(f func() error) error {
	return globalDeclaration.DeclareErrorProne(f)
}

// MustDeclareErrorProne is a variant of the DeclareErrorProne that panics on error.
func MustDeclareErrorProne(f func() error) struct{} {
	if err := DeclareErrorProne(f); err != nil {
		panic(err)
	}
	return struct{}{}
}

// globalContainer specifies the global container.
var globalContainer = NewContainer()

// Provide calls the Provide method of the global container.
func Provide(ctor Constructor) error {
	return globalContainer.Provide(ctor)
}

// MustProvide is a variant of the Provide that panics on error.
func MustProvide(ctor Constructor) {
	if err := Provide(ctor); err != nil {
		panic(err)
	}
}

// Attach calls the Attach method of the global container.
func Attach(proc Processor) error {
	return globalContainer.Attach(proc)
}

// MustAttach is a variant of the Attach that panics on error.
func MustAttach(proc Processor) {
	if err := Attach(proc); err != nil {
		panic(err)
	}
}

// Run calls declared functions if not called yet and then
// calls the Run method of the global container.
func Run(functors ...Functor) error {
	if !globalDeclaration.Fulfilled() {
		if err := globalDeclaration.Fulfill(); err != nil {
			return err
		}
	}
	return globalContainer.Run(functors...)
}

// MustRun is a variant of the Run that panics on error.
func MustRun(functors ...Functor) {
	if err := Run(functors...); err != nil {
		panic(err)
	}
}
