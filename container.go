package kinit

import (
	"reflect"

	"github.com/go-kata/kdone"
	"github.com/go-kata/kerror"
)

// Container represents a dependency injection container.
//
// The usual identifier for variables of this type is ctr.
type Container struct {
	// constructors specifies registered constructors associated with types of objects they are create.
	constructors map[reflect.Type]Constructor
	// processors specifies registered processors associated with types of objects they are process.
	processors map[reflect.Type][]Processor
}

// NewContainer returns a new dependency injection container.
func NewContainer() *Container {
	return &Container{
		constructors: make(map[reflect.Type]Constructor),
		processors:   make(map[reflect.Type][]Processor),
	}
}

// Provide registers the given constructor in this container.
//
// Only one constructor for a type may be registered.
func (c *Container) Provide(ctor Constructor) error {
	if c == nil {
		return kerror.New(kerror.ENil, "nil container cannot register constructor")
	}
	if ctor == nil {
		return kerror.New(kerror.EInvalid, "container cannot register nil constructor")
	}
	t := ctor.Type()
	if t == nil {
		return kerror.New(kerror.EInvalid, "container cannot register constructor for nil type")
	}
	if _, ok := c.constructors[t]; ok {
		return kerror.Newf(kerror.EAmbiguous, "%s constructor already registered", t)
	}
	c.constructors[t] = ctor
	return nil
}

// MustProvide is a variant of the Provide that panics on error.
func (c *Container) MustProvide(ctor Constructor) {
	if err := c.Provide(ctor); err != nil {
		panic(err)
	}
}

// Attach registers the given processor in this container.
//
// Multiple processors may be registered for one type, but there are no guaranty of order of their call.
func (c *Container) Attach(proc Processor) error {
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

// MustAttach is a variant of the Attach that panics on error.
func (c *Container) MustAttach(proc Processor) {
	if err := c.Attach(proc); err != nil {
		panic(err)
	}
}

// Lookup returns constructor and processors that are registered for the given type in this container.
//
// Nil constructor indicates that there are no registered constructor for the type.
func (c *Container) Lookup(t reflect.Type) (Constructor, []Processor) {
	if c == nil || t == nil {
		return nil, nil
	}
	ctor := c.constructors[t]
	var processors []Processor
	if pp, ok := c.processors[t]; ok {
		processors = make([]Processor, len(pp))
		copy(processors, pp)
	}
	return ctor, processors
}

// Explore calls f for each type presented in this container.
//
// Nil constructor indicates that there are no registered constructor for the type
// but registered processors are there.
//
// The traversal will be broken if f will return false.
func (c *Container) Explore(f func(reflect.Type, Constructor, []Processor) (next bool)) {
	if c == nil || f == nil {
		return
	}
	for t, ctor := range c.constructors {
		var processors []Processor
		if pp, ok := c.processors[t]; ok {
			processors = make([]Processor, len(pp))
			copy(processors, pp)
		}
		if !f(t, ctor, processors) {
			return
		}
	}
	for t, pp := range c.processors {
		if _, ok := c.constructors[t]; !ok {
			processors := make([]Processor, len(pp))
			copy(processors, pp)
			if !f(t, nil, processors) {
				return
			}
		}
	}
}

// Run runs given functors sequentially resolving their dependencies recursively using this container.
// If some functor returns further functors all of them will be run before the running of functors that follows it.
//
// All objects created during run will be automatically destroyed when it ends.
func (c *Container) Run(functors ...Functor) (err error) {
	if c == nil {
		return kerror.New(kerror.ENil, "nil container cannot run functors")
	}
	arena := NewArena()
	defer func() {
		err = kerror.Join(err, arena.Finalize())
	}()
	runtime, err := NewRuntime(c, arena)
	if err != nil {
		return err
	}
	if err := arena.Put(reflect.TypeOf(runtime), reflect.ValueOf(runtime), kdone.Noop); err != nil {
		return err
	}
	return c.run(arena, functors)
}

// MustRun is a variant of the Run that panics on error.
func (c *Container) MustRun(functors ...Functor) {
	if err := c.Run(functors...); err != nil {
		panic(err)
	}
}

// run runs given functors using the given arena.
func (c *Container) run(arena *Arena, functors []Functor) error {
	for _, fun := range functors {
		if fun == nil {
			return kerror.New(kerror.EInvalid, "container cannot run nil functor")
		}
		a, err := c.resolveTypes(arena, fun.Parameters())
		if err != nil {
			return err
		}
		further, err := fun.Call(a...)
		if err != nil {
			return err
		}
		if err := c.run(arena, further); err != nil {
			return err
		}
	}
	return nil
}

// resolveType returns the object of the given type. If the object is already on the given arena, it will be used.
// Otherwise it will be firstly created and processed using this container and registered on the arena.
func (c *Container) resolveType(arena *Arena, t reflect.Type) (reflect.Value, error) {
	if t == nil {
		return reflect.Value{}, kerror.New(kerror.EInvalid, "container cannot resolve dependency of nil type")
	}
	if obj, ok := arena.Get(t); ok {
		return obj, nil
	}
	ctor, ok := c.constructors[t]
	if !ok {
		return reflect.Value{}, kerror.Newf(kerror.ENotFound, "%s constructor is not registered", t)
	}
	a, err := c.resolveTypes(arena, ctor.Parameters())
	if err != nil {
		return reflect.Value{}, err
	}
	obj, dtor, err := ctor.Create(a...)
	if err != nil {
		return reflect.Value{}, err
	}
	for _, proc := range c.processors[t] {
		a, err := c.resolveTypes(arena, proc.Parameters())
		if err != nil {
			return reflect.Value{}, err
		}
		if err := proc.Process(obj, a...); err != nil {
			return reflect.Value{}, err
		}
	}
	if err := arena.Put(t, obj, dtor); err != nil {
		return reflect.Value{}, err
	}
	return obj, nil
}

// resolveTypes resolves given types together.
func (c *Container) resolveTypes(arena *Arena, types []reflect.Type) ([]reflect.Value, error) {
	objects := make([]reflect.Value, len(types))
	for i, t := range types {
		obj, err := c.resolveType(arena, t)
		if err != nil {
			return nil, err
		}
		objects[i] = obj
	}
	return objects, nil
}
