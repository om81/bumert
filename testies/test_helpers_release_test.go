//go:build !debug

package testies

import "testing"

// In release builds, the function f is expected to NOT panic (it's a no-op).
func MustPanicIfDebug(t *testing.T, f func()) {
	t.Helper()
	MustNotPanic(t, f) // Expect no panic because assertions are off
}

// In release builds, the function f is expected to NOT panic.
func MustNotPanicIfDebug(t *testing.T, f func()) {
	t.Helper()
	MustNotPanic(t, f)
}
