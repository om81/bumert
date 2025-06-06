// Code generated by gen-release. DO NOT EDIT.
//go:build !debug && !bumert

package bumert

// Assertion is a no-op struct in release builds.
type Assertion struct{}

// Should is a no-op function in release builds.
// It returns a pointer to a singleton no-op Assertion,
// allowing method chains to execute silently without effect.
func Should(value any) *Assertion {
	return &noOpAssertion
}

// A singleton instance for the no-op assertion.
var noOpAssertion Assertion

// BeNil is a no-op method in release builds.
func (a *Assertion) BeNil() *Assertion {
	return a // Return receiver for chainability
}

// NotBeNil is a no-op method in release builds.
func (a *Assertion) NotBeNil() *Assertion {
	return a // Return receiver for chainability
}

// TrueFn is a no-op method in release builds.
func (a *Assertion) TrueFn(f func() bool) *Assertion {
	return a // Return receiver for chainability
}

// BeTrue is a no-op method in release builds.
func (a *Assertion) BeTrue() *Assertion {
	return a // Return receiver for chainability
}

// BeFalse is a no-op method in release builds.
func (a *Assertion) BeFalse() *Assertion {
	return a // Return receiver for chainability
}

// BeEqual is a no-op method in release builds.
func (a *Assertion) BeEqual(expected any) *Assertion {
	return a // Return receiver for chainability
}

// NotBeEqual is a no-op method in release builds.
func (a *Assertion) NotBeEqual(unexpected any) *Assertion {
	return a // Return receiver for chainability
}

// BeEmpty is a no-op method in release builds.
func (a *Assertion) BeEmpty() *Assertion {
	return a // Return receiver for chainability
}

// NotBeEmpty is a no-op method in release builds.
func (a *Assertion) NotBeEmpty() *Assertion {
	return a // Return receiver for chainability
}

// HaveLen is a no-op method in release builds.
func (a *Assertion) HaveLen(expectedLen int) *Assertion {
	return a // Return receiver for chainability
}

// Contain is a no-op method in release builds.
func (a *Assertion) Contain(expectedElement any) *Assertion {
	return a // Return receiver for chainability
}

// NotContain is a no-op method in release builds.
func (a *Assertion) NotContain(unexpectedElement any) *Assertion {
	return a // Return receiver for chainability
}

// ContainSubstring is a no-op method in release builds.
func (a *Assertion) ContainSubstring(substring string) *Assertion {
	return a // Return receiver for chainability
}

// HavePrefix is a no-op method in release builds.
func (a *Assertion) HavePrefix(prefix string) *Assertion {
	return a // Return receiver for chainability
}

// HaveSuffix is a no-op method in release builds.
func (a *Assertion) HaveSuffix(suffix string) *Assertion {
	return a // Return receiver for chainability
}

// BeZero is a no-op method in release builds.
func (a *Assertion) BeZero() *Assertion {
	return a // Return receiver for chainability
}

// NotBeZero is a no-op method in release builds.
func (a *Assertion) NotBeZero() *Assertion {
	return a // Return receiver for chainability
}

// BeGreaterThan is a no-op method in release builds.
func (a *Assertion) BeGreaterThan(expected any) *Assertion {
	return a // Return receiver for chainability
}

// BeLessThan is a no-op method in release builds.
func (a *Assertion) BeLessThan(expected any) *Assertion {
	return a // Return receiver for chainability
}

// BeGreaterThanOrEqualTo is a no-op method in release builds.
func (a *Assertion) BeGreaterThanOrEqualTo(expected any) *Assertion {
	return a // Return receiver for chainability
}

// BeLessThanOrEqualTo is a no-op method in release builds.
func (a *Assertion) BeLessThanOrEqualTo(expected any) *Assertion {
	return a // Return receiver for chainability
}

// BeError is a no-op method in release builds.
func (a *Assertion) BeError() *Assertion {
	return a // Return receiver for chainability
}

// NotBeError is a no-op method in release builds.
func (a *Assertion) NotBeError() *Assertion {
	return a // Return receiver for chainability
}

// BeErrorOfType is a no-op method in release builds.
func (a *Assertion) BeErrorOfType(target any) *Assertion {
	return a // Return receiver for chainability
}

// BeErrorWithMessage is a no-op method in release builds.
func (a *Assertion) BeErrorWithMessage(substring string) *Assertion {
	return a // Return receiver for chainability
}

// Assert is a no-op function in release builds.
func Assert(condition bool) {
	// No-op
}

// Assertf is a no-op function in release builds.
func Assertf(condition bool, format string, args ...any) {
	// No-op
}
