package testies

import "testing"

// Helper function to capture panics
func MustPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

// Helper function to ensure no panic occurs
func MustNotPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked unexpectedly: %v", r)
		}
	}()
	f()
}
