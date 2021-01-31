package kinit

// Bootstrapper represents an invocation bootstrapper.
//
// The usual identifier for variables of this type is boot.
type Bootstrapper interface {
	// Bootstrap bootstraps an invocation.
	Bootstrap(arena *Arena) error
}
