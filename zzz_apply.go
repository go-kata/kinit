package kinit

import "github.com/go-kata/kerror"

// Apply calls the Apply method of the global container.
//
// Deprecated: since 0.4.0, use Attach instead.
func Apply(proc Processor) error {
	return globalContainer.Apply(proc)
}

// MustApply is a variant of the Apply that panics on error.
//
// Deprecated: since 0.4.0, use MustAttach instead.
func MustApply(proc Processor) {
	if err := Apply(proc); err != nil {
		panic(err)
	}
}

// Apply registers the given processor in this container.
//
// Multiple processors may be registered for one type, but there are no guaranty of order of their call.
//
// Deprecated: since 0.4.0, use Attach instead.
func (c *Container) Apply(proc Processor) error {
	if c == nil {
		return kerror.New(kerror.ENil, "nil container cannot register processor")
	}
	if proc == nil {
		return kerror.New(kerror.EInvalid, "container cannot register nil processor")
	}
	t := proc.Type()
	if t == nil {
		return kerror.New(kerror.EInvalid, "container cannot register processor for nil type")
	}
	c.processors[t] = append(c.processors[t], proc)
	return nil
}

// MustApply is a variant of the Apply that panics on error.
//
// Deprecated: since 0.4.0, use MustAttach instead.
func (c *Container) MustApply(proc Processor) {
	if err := c.Apply(proc); err != nil {
		panic(err)
	}
}
