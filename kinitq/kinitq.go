// Package kinitq provides the KInit quality inspection kit.
package kinitq

// globalInspector specifies the global inspector.
var globalInspector = NewInspector()

// Global returns the global inspector.
func Global() *Inspector {
	return globalInspector
}
