// Package kinit provides tools for creating objects.
package kinit

// globalContainer specifies the global container.
var globalContainer = NewContainer()

// Global returns the global container.
func Global() *Container {
	return globalContainer
}
