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
		kerror.NPE()
		return nil
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

// Apply registers the given processor in this container.
//
// Multiple processors may be registered for one type, but there are no guaranty of order of their call.
func (c *Container) Apply(proc Processor) error {
	if c == nil {
		kerror.NPE()
		return nil
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
func (c *Container) MustApply(proc Processor) {
	if err := c.Apply(proc); err != nil {
		panic(err)
	}
}

// Invoke applies given bootstrappers, resolves the dependency graph based on parameters
// of the given executor using this container and then executes an activity. Dependencies
// of each subsequent executor will be resolved dynamically before it's activity execution.
func (c *Container) Invoke(exec Executor, bootstrappers ...Bootstrapper) error {
	if c == nil {
		kerror.NPE()
		return nil
	}
	if exec == nil {
		return kerror.New(kerror.EInvalid, "container cannot invoke nil executor")
	}
	arena := NewArena()
	defer arena.MustFinalize()
	for _, boot := range bootstrappers {
		if boot == nil {
			return kerror.New(kerror.EInvalid, "container cannot apply nil bootstrapper")
		}
		if err := boot.Bootstrap(arena); err != nil {
			return err
		}
	}
	return c.execute(arena, exec)
}

// MustInvoke is a variant of Invoke that panics on error.
func (c *Container) MustInvoke(exec Executor, bootstrappers ...Bootstrapper) {
	if err := c.Invoke(exec, bootstrappers...); err != nil {
		panic(err)
	}
}

// get creates if needed and returns an object of the given type.
func (c *Container) get(arena *Arena, t reflect.Type) (reflect.Value, error) {
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
	obj, dtor, err := c.create(arena, ctor)
	if err != nil {
		return reflect.Value{}, err
	}
	for _, proc := range c.processors[t] {
		if err := c.process(arena, obj, proc); err != nil {
			return reflect.Value{}, err
		}
	}
	arena.MustRegister(t, obj, dtor)
	return obj, nil
}

// create creates and returns a new object.
func (c *Container) create(arena *Arena, ctor Constructor) (reflect.Value, kdone.Destructor, error) {
	p := ctor.Parameters()
	a := make([]reflect.Value, len(p))
	for i := range p {
		aobj, err := c.get(arena, p[i])
		if err != nil {
			return reflect.Value{}, nil, err
		}
		a[i] = aobj
	}
	return ctor.Create(a...)
}

// process processes the given object.
func (c *Container) process(arena *Arena, obj reflect.Value, proc Processor) error {
	p := proc.Parameters()
	a := make([]reflect.Value, len(p))
	for i := range p {
		aobj, err := c.get(arena, p[i])
		if err != nil {
			return err
		}
		a[i] = aobj
	}
	return proc.Process(obj, a...)
}

// execute executes an activity.
func (c *Container) execute(arena *Arena, exec Executor) error {
	p := exec.Parameters()
	a := make([]reflect.Value, len(p))
	for i := range p {
		aobj, err := c.get(arena, p[i])
		if err != nil {
			return err
		}
		a[i] = aobj
	}
	next, err := exec.Execute(a...)
	if err != nil {
		return err
	}
	if next != nil {
		return c.execute(arena, next)
	}
	return nil
}
