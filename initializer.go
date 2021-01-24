package kinit

// Initializer represents an invocation initializer.
//
// The usual identifier for variables of this type is init.
type Initializer interface {
	// Initialize initializes an invocation.
	Initialize(arena *Arena) error
}
