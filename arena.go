package kinit

import (
	"reflect"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

// Arena represents an objects holder.
type Arena struct {
	// parents specifies parent arenas.
	parents []*Arena
	// objects specifies registered objects.
	objects map[reflect.Type]reflect.Value
	// reaper specifies the reaper for registered objects.
	reaper *kdone.Reaper
	// finalized specifies whether were registered objects destroyed.
	finalized bool
}

// NewArena returns a new arena with given parent arenas.
func NewArena(parents ...*Arena) *Arena {
	a := &Arena{
		objects: make(map[reflect.Type]reflect.Value),
		reaper:  kdone.NewReaper(),
	}
	if len(parents) > 0 {
		a.parents = make([]*Arena, len(parents))
		copy(a.parents, parents)
	}
	return a
}

// Put registers the given object on this arena.
func (a *Arena) Put(t reflect.Type, obj reflect.Value, dtor kdone.Destructor) error {
	if a == nil {
		return kerror.New(kerror.ENil, "nil arena cannot register object")
	}
	if a.finalized {
		return kerror.New(kerror.EIllegal, "arena has already destroyed objects")
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

// MustPut is a variant of the Put that panics on error.
func (a *Arena) MustPut(t reflect.Type, obj reflect.Value, dtor kdone.Destructor) {
	if err := a.Put(t, obj, dtor); err != nil {
		panic(err)
	}
}

// Get returns an object of the given type if registered on this arena
// or on the one of non-finalized parent arenas which will be bypassed
// in the order that they were passed to the NewArena.
func (a *Arena) Get(t reflect.Type) (obj reflect.Value, ok bool) {
	if a == nil || t == nil {
		return reflect.Value{}, false
	}
	if obj, ok = a.objects[t]; ok {
		return
	}
	for _, parent := range a.parents {
		if parent.Finalized() {
			continue
		}
		if obj, ok = parent.Get(t); ok {
			return
		}
	}
	return
}

// Finalize destroys objects registered on this arena.
func (a *Arena) Finalize() error {
	if a == nil {
		return nil
	}
	if a.finalized {
		return kerror.New(kerror.EIllegal, "arena has already destroyed objects")
	}
	defer func() {
		a.finalized = true
	}()
	return a.reaper.Finalize()
}

// MustFinalize is a variant of the Finalize that panics on error.
func (a *Arena) MustFinalize() {
	if err := a.Finalize(); err != nil {
		panic(err)
	}
}

// Finalized returns boolean specifies were objects registered on this arena destroyed.
func (a *Arena) Finalized() bool {
	if a == nil {
		return false
	}
	return a.finalized
}
