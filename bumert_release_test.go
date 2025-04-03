//go:build !debug && !bumert

package bumert_test

import (
	"testing"

	"github.com/deblasis/bumert"
	"github.com/deblasis/bumert/testies"
)

// TestAssertionBehaviorRelease checks assertions in release mode.
// Assertions should be no-ops and never panic.
func TestAssertionBehaviorRelease(t *testing.T) {
	var nilPtr *int = nil
	var nonNilPtr = new(int)

	// --- BeNil --- (Release Mode: Should never panic)
	testies.MustNotPanic(t, func() { bumert.Should(nonNilPtr).BeNil() }) // Would panic in debug
	testies.MustNotPanic(t, func() { bumert.Should(nilPtr).BeNil() })

	// --- NotBeNil --- (Release Mode: Should never panic)
	testies.MustNotPanic(t, func() { bumert.Should(nilPtr).NotBeNil() }) // Would panic in debug
	testies.MustNotPanic(t, func() { bumert.Should(nonNilPtr).NotBeNil() })
}

// TestTrueFnAssertionRelease checks TrueFn assertion in release mode.
// Assertions should be no-ops and never panic.
func TestTrueFnAssertionRelease(t *testing.T) {
	// --- TrueFn --- (Release Mode: Should never panic)
	testies.MustNotPanic(t, func() { bumert.Should(nil).TrueFn(func() bool { return true }) })
	// testies.MustNotPanic(t, func() { should.ShouldTrueFn(func() bool { return true }) })

	// Function returns false, but should not panic in release
	testies.MustNotPanic(t, func() { bumert.Should(nil).TrueFn(func() bool { return false }) }) // Would panic in debug
	// testies.MustNotPanic(t, func() { should.ShouldTrueFn(func() bool { return false }) })       // Would panic in debug
}
