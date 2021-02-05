package kinit

import (
	"reflect"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

// Runtime represents a runtime.
type Runtime struct {
	// container specifies the container associated with this runtime.
	container *Container
	// arena specifies the arena associated with this runtime.
	arena *Arena
}

// NewRuntime returns a new runtime associated with given container and arena.
func NewRuntime(ctr *Container, arena *Arena) (*Runtime, error) {
	if ctr == nil {
		return nil, kerror.New(kerror.EInvalid, "runtime cannot be associated with nil container")
	}
	if arena == nil {
		return nil, kerror.New(kerror.EInvalid, "runtime cannot be associated with nil arena")
	}
	return &Runtime{
		container: ctr,
		arena:     arena,
	}, nil
}

// MustNewRuntime is a variant of the NewRuntime that panics on error.
func MustNewRuntime(ctr *Container, arena *Arena) *Runtime {
	r, err := NewRuntime(ctr, arena)
	if err != nil {
		panic(err)
	}
	return r
}

// Register registers the given object on the associated arena.
func (r *Runtime) Register(t reflect.Type, obj reflect.Value, dtor kdone.Destructor) error {
	if r == nil {
		return kerror.New(kerror.ENil, "nil runtime cannot register object")
	}
	return r.arena.Register(t, obj, dtor)
}

// MustRegister is a variant of the Register that panics on error.
func (r *Runtime) MustRegister(t reflect.Type, obj reflect.Value, dtor kdone.Destructor) {
	if err := r.Register(t, obj, dtor); err != nil {
		panic(err)
	}
}

// Run runs given functors using the associated container.
// The created separate arena will use the associated arena as a parent.
func (r *Runtime) Run(functors ...Functor) (err error) {
	if r == nil {
		return kerror.New(kerror.ENil, "nil runtime cannot run functors")
	}
	arena := NewArena(r.arena)
	defer func() {
		err = kerror.Join(err, arena.Finalize())
	}()
	runtime, err := NewRuntime(r.container, arena)
	if err != nil {
		return err
	}
	if err := arena.Register(reflect.TypeOf(runtime), reflect.ValueOf(runtime), kdone.Noop); err != nil {
		return err
	}
	return r.container.run(arena, functors...)
}

// MustRun is a variant of Run that panics on error.
func (r *Runtime) MustRun(functors ...Functor) {
	if err := r.Run(functors...); err != nil {
		panic(err)
	}
}
