//go:build debug || bumert

package bumert_test

import (
	"testing"

	"github.com/deblasis/bumert"
	"github.com/deblasis/bumert/testies"
)

// mustPanic is a local helper for debug tests.
// It checks if the function f panics, failing the test if it does not.
func mustPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

// TestAssertionBehavior checks assertions based on build tags.
// This test itself needs to be run with and without the 'debug' tag
// to verify both scenarios.
func TestAssertionBehavior(t *testing.T) {
	var nilPtr *int = nil
	var nonNilPtr = new(int)

	// --- BeNil ---
	// Expect panic in debug, no-op in release
	mustPanic(t, func() { bumert.Should(nonNilPtr).BeNil() })
	// Expect no panic in either build
	testies.MustNotPanic(t, func() { bumert.Should(nilPtr).BeNil() })

	// --- NotBeNil ---
	// Expect panic in debug, no-op in release
	mustPanic(t, func() { bumert.Should(nilPtr).NotBeNil() })
	// Expect no panic in either build
	testies.MustNotPanic(t, func() { bumert.Should(nonNilPtr).NotBeNil() })
}

// TestTrueFnAssertion checks the TrueFn assertion and ShouldTrueFn alias.
func TestTrueFnAssertion(t *testing.T) {
	// --- TrueFn ---
	// Expect no panic in either build (function returns true)
	testies.MustNotPanic(t, func() { bumert.Should(nil).TrueFn(func() bool { return true }) })
	// testies.MustNotPanic(t, func() { should.ShouldTrueFn(func() bool { return true }) })

	// Expect panic in debug, no-op in release (function returns false)
	mustPanic(t, func() { bumert.Should(nil).TrueFn(func() bool { return false }) })
	// mustPanic(t, func() { should.ShouldTrueFn(func() bool { return false }) })
}
