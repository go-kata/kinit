package kinit

import (
	"reflect"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

// Arena represents an invocation context.
type Arena struct {
	// objects specifies the list of registered objects.
	objects map[reflect.Type]reflect.Value
	// reaper specifies the reaper for registered objects.
	reaper *kdone.Reaper
}

// NewArena returns a new arena.
func NewArena() *Arena {
	return &Arena{
		objects: make(map[reflect.Type]reflect.Value),
		reaper:  kdone.NewReaper(),
	}
}

// Register registers the given object on this arena.
func (a *Arena) Register(t reflect.Type, obj reflect.Value, dtor kdone.Destructor) error {
	if a == nil {
		kerror.NPE()
		return nil
	}
	if t == nil {
		return kerror.New(kerror.EInvalid, "arena cannot register object of nil type")
	}
	if _, ok := a.objects[t]; ok {
		return kerror.Newf(kerror.EAmbiguous, "%s object already registered", t)
	}
	if err := a.reaper.Assume(dtor); err != nil {
		return err
	}
	a.objects[t] = obj
	return nil
}

// MustRegister is a variant of the Register that panics on error.
func (a *Arena) MustRegister(t reflect.Type, obj reflect.Value, dtor kdone.Destructor) {
	if err := a.Register(t, obj, dtor); err != nil {
		panic(err)
	}
}

// Get returns an object of the given type if registered on this arena.
func (a *Arena) Get(t reflect.Type) (obj reflect.Value, ok bool) {
	if a == nil || t == nil {
		return reflect.Value{}, false
	}
	obj, ok = a.objects[t]
	return
}

// Finalize destroys objects registered on this arena.
func (a *Arena) Finalize() error {
	if a == nil {
		return nil
	}
	return a.reaper.Finalize()
}

// MustFinalize is a variant of the Finalize that panics on error.
func (a *Arena) MustFinalize() {
	if err := a.Finalize(); err != nil {
		panic(err)
	}
}
