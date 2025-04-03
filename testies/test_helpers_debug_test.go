//go:build debug

package testies

import "testing"

// In debug builds, the function f is expected to panic.
func MustPanicIfDebug(t *testing.T, f func()) {
	t.Helper()
	MustPanic(t, f)
}

// In debug builds, the function f is NOT expected to panic.
func MustNotPanicIfDebug(t *testing.T, f func()) {
	t.Helper()
	MustNotPanic(t, f)
}
