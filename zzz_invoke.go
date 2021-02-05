package kinit

import "github.com/go-kata/kerror"

// Invoke calls declared functions if not called yet and then
// calls the Invoke method of the global container.
//
// Deprecated: since 0.4.0, use Run instead.
func Invoke(exec Executor, bootstrappers ...Bootstrapper) error {
	if !globalDeclaration.Fulfilled() {
		if err := globalDeclaration.Fulfill(); err != nil {
			return err
		}
	}
	return globalContainer.Invoke(exec, bootstrappers...)
}

// MustInvoke is a variant of the Invoke that panics on error.
//
// Deprecated: since 0.4.0, use MustRun instead.
func MustInvoke(exec Executor, bootstrappers ...Bootstrapper) {
	if err := Invoke(exec, bootstrappers...); err != nil {
		panic(err)
	}
}

// Invoke applies given bootstrappers, resolves the dependency graph based on parameters
// of the given executor using this container and then executes an activity. Dependencies
// of each subsequent executor will be resolved dynamically before it's activity execution.
//
// Deprecated: since 0.4.0, use Run instead.
func (c *Container) Invoke(exec Executor, bootstrappers ...Bootstrapper) (err error) {
	if c == nil {
		return kerror.New(kerror.ENil, "nil container cannot invoke executor")
	}
	if exec == nil {
		return kerror.New(kerror.EInvalid, "container cannot invoke nil executor")
	}
	arena := NewArena()
	defer func() {
		err = kerror.Join(err, arena.Finalize())
	}()
	for _, boot := range bootstrappers {
		if boot == nil {
			return kerror.New(kerror.EInvalid, "container cannot apply nil bootstrapper")
		}
		if err := boot.Bootstrap(arena); err != nil {
			return err
		}
	}
	return c.invoke(arena, exec)
}

// MustInvoke is a variant of Invoke that panics on error.
//
// Deprecated: since 0.4.0, use MustRun instead.
func (c *Container) MustInvoke(exec Executor, bootstrappers ...Bootstrapper) {
	if err := c.Invoke(exec, bootstrappers...); err != nil {
		panic(err)
	}
}

// invoke executes an activity using the given arena.
func (c *Container) invoke(arena *Arena, exec Executor) error {
	a, err := c.resolve(arena, exec.Parameters())
	if err != nil {
		return err
	}
	next, err := exec.Execute(a...)
	if err != nil {
		return err
	}
	if next != nil {
		return c.invoke(arena, next)
	}
	return nil
}
