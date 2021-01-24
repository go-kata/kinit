package kinit

// hooks specifies functions globally registered to be called on the default container invocation.
var hooks []func() error

// Hook registers the given function as a hook and returns it's index. Nil function will be discarded.
//
// Hooks may be used to protect from slow reflection calls without a subsequent default container invocation.
// This is useful for libraries that registers defaults for this package.
func Hook(f func() error) int {
	if f == nil {
		return -1
	}
	hooks = append(hooks, f)
	return len(hooks) - 1
}

// MustHook is a variant of the Hook registering a function that panics on error.
func MustHook(f func()) int {
	if f == nil {
		return -1
	}
	return Hook(func() error {
		f()
		return nil
	})
}

// callHooks calls registered hooks.
func callHooks() error {
	for _, hook := range hooks {
		if err := hook(); err != nil {
			return err
		}
	}
	return nil
}
