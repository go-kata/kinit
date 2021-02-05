package kinit

// Bootstrapper represents an invocation bootstrapper.
//
// The usual identifier for variables of this type is boot.
//
// Deprecated: since 0.4.0, use Functor and Runtime instead.
type Bootstrapper interface {
	// Bootstrap bootstraps an invocation.
	Bootstrap(arena *Arena) error
}
