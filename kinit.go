// Package kinit provides tools for creating objects.
//
// In contrast to the Wire (github.com/google/wire) the dependency injection mechanism provided by this package
// doesn't require any codegeneration and relies on the reflection. The dependency injection is a process that
// takes place only once on a program startup in the vast majority of cases and the low performance of reflection
// in this case is not crucial. On other hand this solution gives more control over a program on debug (because
// it includes injection process itself) and doesn't divide program to the real code and configuration. Also it
// makes the dependency injection process more customizable thanks to interfaces.
//
// However, taking in account that the reflection is slow, it is better to use hooks provided by this package
// instead of raw init functions to fill up the default container in libraries. It avoids slow reflection calls
// when library entities are used manually.
//
// You may divide resolving of the dependency graph into steps using the executors chaining. It makes possible
// to chose the next step depending on some conditions, e.g. depends on what command of console program is
// executed for now or what module is mounted to the extension point of framework.
//
// Objects are uniquely identified by container using reflection of their types. Default implementations of
// interfaces from this package based on the reflect.Type interface directly are provided by the
// github.com/go-kata/kinitx package. To extend the objects identification (for example, with custom names
// like in the github.com/uber-go/dig) you may write your own implementation of the type reflection like a
//
//     type NamedType struct {
//         reflect.Type
//         Name string
//     }
//
// and then use it as a base of your own implementations of package interfaces.
//
// Objects are created when container is invoked and destroyed (finalized) at the end of invocation. Object is
// destroyed using it's destructor returning by it's constructor on creation. You may be sure that all destructors
// of correctly created objects will be guaranteed called in a correct order even in case of panic at any step of
// the container invocation.
package kinit

// defaultContainer specifies the default global container.
var defaultContainer = NewContainer()

// Provide calls the Provide method of the default container.
func Provide(ctor Constructor) error {
	return defaultContainer.Provide(ctor)
}

// MustProvide calls the MustProvide method of the default container.
func MustProvide(ctor Constructor) {
	defaultContainer.MustProvide(ctor)
}

// Apply calls the Apply method of the default container.
func Apply(proc Processor) error {
	return defaultContainer.Apply(proc)
}

// MustApply calls the MustApply method of the default container.
func MustApply(proc Processor) {
	defaultContainer.MustApply(proc)
}

// Invoke calls the Invoke method of the default container.
func Invoke(exec Executor, initializers ...Initializer) error {
	if err := callHooks(); err != nil {
		return err
	}
	return defaultContainer.Invoke(exec, initializers...)
}

// MustInvoke calls the MustInvoke method of the default container.
func MustInvoke(exec Executor, initializers ...Initializer) {
	if err := callHooks(); err != nil {
		panic(err)
	}
	defaultContainer.MustInvoke(exec, initializers...)
}
