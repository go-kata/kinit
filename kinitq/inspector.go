package kinitq

import (
	"reflect"

	"github.com/go-kata/kerror"
	"github.com/go-kata/kinit"
)

// Inspector represents a container inspector.
type Inspector struct {
	// types specifies registered types.
	//
	// Here true means all dependencies of this type are permanently satisfied
	// (e.g. via the kinit.Runtime) and this type must be ignored on inspection
	// and false means it must be inspected that dependencies of this type can
	// be successfully satisfied.
	types map[reflect.Type]bool
}

// NewInspector returns a new inspector.
func NewInspector() *Inspector {
	return &Inspector{
		types: make(map[reflect.Type]bool),
	}
}

// Require registers the given required type in this inspector.
func (i *Inspector) Require(t reflect.Type) error {
	if i == nil {
		return kerror.New(kerror.ENil, "nil inspector cannot register type")
	}
	if t == nil {
		return kerror.New(kerror.EInvalid, "inspector cannot register nil type")
	}
	if _, ok := i.types[t]; !ok {
		i.types[t] = false
	}
	return nil
}

// MustRequire us a variant of the Require that panics on error.
func (i *Inspector) MustRequire(t reflect.Type) {
	if err := i.Require(t); err != nil {
		panic(err)
	}
}

// Ignore registers the given ignored type in this inspector.
func (i *Inspector) Ignore(t reflect.Type) error {
	if i == nil {
		return kerror.New(kerror.ENil, "nil inspector cannot register type")
	}
	if t == nil {
		return kerror.New(kerror.EInvalid, "inspector cannot register nil type")
	}
	i.types[t] = true
	return nil
}

// MustIgnore us a variant of the Ignore that panics on error.
func (i *Inspector) MustIgnore(t reflect.Type) {
	if err := i.Ignore(t); err != nil {
		panic(err)
	}
}

// Options represents inspection options.
type Options struct {
	// InspectOnlyRequired specifies whether to inspect only required types
	// skipping all constructors and processors unnecessary for their satisfaction.
	InspectOnlyRequired bool
	// AllowIrrelevantProcessors specifies whether to suppress errors indicating
	// the presence of processors for types without constructors.
	//
	// This option only applies if the InspectOnlyRequired is false.
	AllowIrrelevantProcessors bool
}

// background represents an inspection background.
type background struct {
	// history specifies the inspection history.
	//
	// If a type is absent in the history it means that it wasn't inspected yet.
	// Otherwise true means that the type inspection was already done
	// and false means that a type is currently being inspected.
	history map[reflect.Type]bool
	// stack specifies the inspection stack.
	stack []reflect.Type
}

// Inspect inspects the given container for the absence of cyclic and unsatisfied dependencies.
func (i *Inspector) Inspect(ctr *kinit.Container, opt *Options) error {
	if i == nil {
		return nil
	}
	if opt == nil {
		opt = &Options{}
	}
	coerr := kerror.NewCollector()
	bg := &background{
		history: make(map[reflect.Type]bool),
	}
	for t, ignore := range i.types {
		if ignore {
			continue
		}
		coerr.Collect(i.inspectType(ctr, t, bg))
	}
	if !opt.InspectOnlyRequired {
		ctr.Explore(func(t reflect.Type, ctor kinit.Constructor, processors []kinit.Processor) (next bool) {
			if ctor != nil {
				coerr.Collect(i.inspectType(ctr, t, bg))
				return true
			}
			if !opt.AllowIrrelevantProcessors {
				coerr.Collect(kerror.Newf(kerror.EInvalid, "%s processor(s) found in absence of constructor", t))
			}
			for _, proc := range processors {
				coerr.Collect(i.inspectTypes(ctr, proc.Parameters(), bg))
			}
			return true
		})
	}
	return coerr.Error()
}

// MustInspect is a variant of the Inspect that panics on error.
func (i *Inspector) MustInspect(ctr *kinit.Container, opt *Options) {
	if err := i.Inspect(ctr, opt); err != nil {
		panic(err)
	}
}

// inspectType inspects that the dependency of the given type
// can be successfully satisfied by the given container.
func (i *Inspector) inspectType(ctr *kinit.Container, t reflect.Type, bg *background) error {
	if i.types[t] {
		return nil
	}
	if ended, begun := bg.history[t]; begun {
		if ended {
			return nil
		}
		s := t.String()
		n := len(bg.stack)
		for j := n - 1; j >= 0; j-- {
			if bg.stack[j] == t {
				for k := j + 1; k < n; k++ {
					s += " 🠖 " + bg.stack[k].String()
				}
				break
			}
		}
		s += " 🠖 " + t.String()
		return kerror.Newf(kerror.EAmbiguous, "cyclic dependency: %s", s)
	}
	bg.history[t] = false
	defer func() {
		bg.history[t] = true
	}()
	ctor, processors := ctr.Lookup(t)
	if ctor == nil {
		s := t.String()
		if n := len(bg.stack); n > 0 {
			s = bg.stack[n-1].String() + " 🠖 " + s
		}
		return kerror.Newf(kerror.ENotFound, "unsatisfied dependency: %s", s)
	}
	bg.stack = append(bg.stack, t)
	defer func() {
		bg.stack = bg.stack[:len(bg.stack)-1]
	}()
	coerr := kerror.NewCollector()
	coerr.Collect(i.inspectTypes(ctr, ctor.Parameters(), bg))
	for _, proc := range processors {
		coerr.Collect(i.inspectTypes(ctr, proc.Parameters(), bg))
	}
	return coerr.Error()
}

// inspectTypes inspects given types together.
func (i *Inspector) inspectTypes(ctr *kinit.Container, types []reflect.Type, bg *background) error {
	coerr := kerror.NewCollector()
	for _, t := range types {
		coerr.Collect(i.inspectType(ctr, t, bg))
	}
	return coerr.Error()
}
