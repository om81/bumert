//go:build debug || bumert

package bumert

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// getCallerInfo retrieves the file and line number of the caller.
// It skips skipFrames stack frames.
func getCallerInfo(skipFrames int) string {
	_, file, line, ok := runtime.Caller(skipFrames + 1) // +1 to skip getCallerInfo itself
	if !ok {
		return "???:?" // Fallback if caller info is unavailable
	}
	// Find the last path separator to trim the prefix
	lastSlash := strings.LastIndex(file, "/")
	if lastSlash >= 0 {
		file = file[lastSlash+1:] // Keep only the filename
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// failAssertion formats the error message with caller info and panics.
func failAssertion(format string, args ...any) {
	callerInfo := getCallerInfo(2) // Skip failAssertion and the assertion method itself
	panic(fmt.Sprintf("%s: assertion failed: %s", callerInfo, fmt.Sprintf(format, args...)))
}

// Should starts a fluent assertion chain for debug builds.
// It performs checks only when the 'debug' build tag is enabled.
func Should(value any) *Assertion {
	return &Assertion{
		value: value,
	}
}

// Assertion holds the value being asserted and provides assertion methods.
// This struct is primarily used in debug builds.
type Assertion struct {
	value any
}

// isNil checks if the underlying value is nil using reflection.
// It correctly handles typed nil pointers, interfaces containing nil pointers, etc.
func isNil(value any) bool {
	if value == nil {
		return true // Catches untyped nil
	}
	// Use reflection to check if it's a typed nil (like *int(nil))
	// wrapped in an interface{}.
	v := reflect.ValueOf(value)
	switch v.Kind() {
	// Check kinds that can be nil
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	default:
		// Other kinds (structs, basic types) cannot be nil in the standard sense
		return false
	}
}

// BeNil checks if the asserted value is nil.
// Panics if the assertion fails in debug builds.
func (a *Assertion) BeNil() *Assertion {
	if !isNil(a.value) {
		failAssertion("should be nil: value (%#v) of type %T is not nil", a.value, a.value)
	}
	return a
}

// NotBeNil checks if the asserted value is not nil.
// Panics if the assertion fails in debug builds.
func (a *Assertion) NotBeNil() *Assertion {
	if isNil(a.value) {
		failAssertion("should not be nil: value is nil")
	}
	return a
}

// TrueFn checks if the provided function returns true.
// The original asserted value (from Should()) is ignored by this check.
// Panics if the function returns false in debug builds.
func (a *Assertion) TrueFn(f func() bool) *Assertion {
	if !f() {
		failAssertion("function returned false")
	}
	return a
}

// --- New Assertions ---

// BeTrue checks if the asserted value is true.
// Panics if the value is not boolean true.
func (a *Assertion) BeTrue() *Assertion {
	if b, ok := a.value.(bool); !ok || !b {
		failAssertion("should be true, but got: %#v (type %T)", a.value, a.value)
	}
	return a
}

// BeFalse checks if the asserted value is false.
// Panics if the value is not boolean false.
func (a *Assertion) BeFalse() *Assertion {
	if b, ok := a.value.(bool); !ok || b {
		failAssertion("should be false, but got: %#v (type %T)", a.value, a.value)
	}
	return a
}

// BeEqual checks if the asserted value is equal to the expected value.
// Uses reflect.DeepEqual for comparison.
// Panics if the values are not equal.
func (a *Assertion) BeEqual(expected any) *Assertion {
	if !reflect.DeepEqual(a.value, expected) {
		failAssertion(`should be equal:
  expected: %#v (type %T)
       got: %#v (type %T)`,
			expected, expected, a.value, a.value)
	}
	return a
}

// NotBeEqual checks if the asserted value is not equal to the expected value.
// Uses reflect.DeepEqual for comparison.
// Panics if the values are equal.
func (a *Assertion) NotBeEqual(unexpected any) *Assertion {
	if reflect.DeepEqual(a.value, unexpected) {
		failAssertion("should not be equal, but got: %#v (type %T)", a.value, a.value)
	}
	return a
}

// getLength returns the length of slice, array, map, chan, or string.
// Returns -1 and false if the type is not supported.
func getLength(value any) (int, bool) {
	if value == nil {
		return -1, false
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return v.Len(), true
	default:
		return -1, false
	}
}

// BeEmpty checks if the asserted value is empty (e.g., nil, zero length slice/map/string).
// Panics if the value is not empty.
func (a *Assertion) BeEmpty() *Assertion {
	if isNil(a.value) {
		return a // nil is considered empty
	}
	if length, ok := getLength(a.value); ok {
		if length == 0 {
			return a // Zero length is empty
		}
		failAssertion("should be empty, but got length %d: %#v", length, a.value)
	} else {
		// If we can't get length and it's not nil, consider it not empty
		failAssertion("should be empty, but got non-nil, non-length-testable value: %#v (type %T)", a.value, a.value)
	}
	return a
}

// NotBeEmpty checks if the asserted value is not empty.
// Panics if the value is empty.
func (a *Assertion) NotBeEmpty() *Assertion {
	if isNil(a.value) {
		failAssertion("should not be empty, but got nil")
	}
	if length, ok := getLength(a.value); ok && length == 0 {
		failAssertion("should not be empty, but got zero length: %#v", a.value)
	}
	// If it's not nil and either has length > 0 or isn't length-testable, it's not empty
	return a
}

// HaveLen checks if the asserted collection (slice, array, map, chan, string) has the expected length.
// Panics if the value is not a collection type or if the length does not match.
func (a *Assertion) HaveLen(expectedLen int) *Assertion {
	if length, ok := getLength(a.value); ok {
		if length != expectedLen {
			failAssertion("should have length %d, but got %d: %#v", expectedLen, length, a.value)
		}
	} else {
		failAssertion("should have length, but got non-length-testable value: %#v (type %T)", a.value, a.value)
	}
	return a
}

// Contain checks if the asserted slice, array, or string contains the expected element/substring.
// For slices/arrays, it iterates and uses reflect.DeepEqual.
// For strings, it uses strings.Contains.
// Panics if the element/substring is not found or the type is unsupported.
func (a *Assertion) Contain(expectedElement any) *Assertion {
	val := reflect.ValueOf(a.value)
	switch val.Kind() {
	case reflect.String:
		sub, ok := expectedElement.(string)
		if !ok {
			failAssertion("Contain expected a string element for string value, got %T", expectedElement)
		}
		if !strings.Contains(val.String(), sub) {
			failAssertion("string %#v should contain substring %#v", a.value, sub)
		}
	case reflect.Slice, reflect.Array:
		found := false
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(val.Index(i).Interface(), expectedElement) {
				found = true
				break
			}
		}
		if !found {
			failAssertion("collection %#v should contain element %#v", a.value, expectedElement)
		}
	// case reflect.Map: // Checking map containment is ambiguous (key or value?), omitting for now.
	default:
		failAssertion("Contain requires slice, array, or string, got %T", a.value)
	}
	return a
}

// NotContain checks if the asserted slice, array, or string does NOT contain the expected element/substring.
// For slices/arrays, it iterates and uses reflect.DeepEqual.
// For strings, it uses strings.Contains.
// Panics if the element/substring IS found or the type is unsupported.
func (a *Assertion) NotContain(unexpectedElement any) *Assertion {
	val := reflect.ValueOf(a.value)
	switch val.Kind() {
	case reflect.String:
		sub, ok := unexpectedElement.(string)
		if !ok {
			failAssertion("NotContain expected a string element for string value, got %T", unexpectedElement)
		}
		if strings.Contains(val.String(), sub) {
			failAssertion("string %#v should not contain substring %#v", a.value, sub)
		}
	case reflect.Slice, reflect.Array:
		found := false
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(val.Index(i).Interface(), unexpectedElement) {
				found = true
				break
			}
		}
		if found {
			failAssertion("collection %#v should not contain element %#v", a.value, unexpectedElement)
		}
	// case reflect.Map: // Checking map containment is ambiguous (key or value?), omitting for now.
	default:
		failAssertion("NotContain requires slice, array, or string, got %T", a.value)
	}
	return a
}

// ContainSubstring checks if the asserted string contains the expected substring.
// Panics if the value is not a string or does not contain the substring.
func (a *Assertion) ContainSubstring(substring string) *Assertion {
	str, ok := a.value.(string)
	if !ok {
		failAssertion("ContainSubstring requires a string value, got %T", a.value)
	}
	if !strings.Contains(str, substring) {
		failAssertion("string %#v should contain substring %#v", str, substring)
	}
	return a
}

// HavePrefix checks if the asserted string has the expected prefix.
// Panics if the value is not a string or does not have the prefix.
func (a *Assertion) HavePrefix(prefix string) *Assertion {
	str, ok := a.value.(string)
	if !ok {
		failAssertion("HavePrefix requires a string value, got %T", a.value)
	}
	if !strings.HasPrefix(str, prefix) {
		failAssertion("string %#v should have prefix %#v", str, prefix)
	}
	return a
}

// HaveSuffix checks if the asserted string has the expected suffix.
// Panics if the value is not a string or does not have the suffix.
func (a *Assertion) HaveSuffix(suffix string) *Assertion {
	str, ok := a.value.(string)
	if !ok {
		failAssertion("HaveSuffix requires a string value, got %T", a.value)
	}
	if !strings.HasSuffix(str, suffix) {
		failAssertion("string %#v should have suffix %#v", str, suffix)
	}
	return a
}

// isZero checks if the value is the zero value for its type.
func isZero(value any) bool {
	if value == nil {
		return true // Consider nil as zero for pointer types etc. handled by isNil before this
	}
	v := reflect.ValueOf(value)
	// Check if the value is the zero value for its type
	return v.IsZero()
}

// BeZero checks if the asserted value is the zero value for its type (e.g., 0, "", false, nil pointer).
// Panics if the value is not the zero value.
func (a *Assertion) BeZero() *Assertion {
	if !isZero(a.value) {
		failAssertion("should be the zero value, but got: %#v (type %T)", a.value, a.value)
	}
	return a
}

// NotBeZero checks if the asserted value is NOT the zero value for its type.
// Panics if the value is the zero value.
func (a *Assertion) NotBeZero() *Assertion {
	if isZero(a.value) {
		failAssertion("should not be the zero value, but got: %#v (type %T)", a.value, a.value)
	}
	return a
}

// compare performs numeric comparison (>, <, >=, <=) using reflection.
// Returns comparison result and true if comparable, false otherwise.
// op should be one of: ">", "<", ">=", "<=".
func compare(v1, v2 reflect.Value, op string) (bool, bool) {
	k1, k2 := v1.Kind(), v2.Kind()

	// Handle direct comparison for identical types first
	if v1.Type() == v2.Type() {
		switch k1 {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n1, n2 := v1.Int(), v2.Int()
			switch op {
			case ">":
				return n1 > n2, true
			case "<":
				return n1 < n2, true
			case ">=":
				return n1 >= n2, true
			case "<=":
				return n1 <= n2, true
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n1, n2 := v1.Uint(), v2.Uint()
			switch op {
			case ">":
				return n1 > n2, true
			case "<":
				return n1 < n2, true
			case ">=":
				return n1 >= n2, true
			case "<=":
				return n1 <= n2, true
			}
		case reflect.Float32, reflect.Float64:
			n1, n2 := v1.Float(), v2.Float()
			switch op {
			case ">":
				return n1 > n2, true
			case "<":
				return n1 < n2, true
			case ">=":
				return n1 >= n2, true
			case "<=":
				return n1 <= n2, true
			}
		}
	} // End identical type check

	// If types differ, attempt conversion to float64 for comparison
	isNumeric := func(k reflect.Kind) bool {
		return (k >= reflect.Int && k <= reflect.Int64) ||
			(k >= reflect.Uint && k <= reflect.Uintptr) ||
			(k >= reflect.Float32 && k <= reflect.Float64)
	}

	if isNumeric(k1) && isNumeric(k2) {
		var f1, f2 float64
		var ok1, ok2 bool

		f1, ok1 = convertToFloat64(v1)
		f2, ok2 = convertToFloat64(v2)

		if ok1 && ok2 {
			switch op {
			case ">":
				return f1 > f2, true
			case "<":
				return f1 < f2, true
			case ">=":
				return f1 >= f2, true
			case "<=":
				return f1 <= f2, true
			}
		}
	}

	return false, false // Types are not comparable or convertible
}

// convertToFloat64 converts numeric reflect.Value to float64.
// Returns the float64 value and true on success, 0 and false otherwise.
func convertToFloat64(v reflect.Value) (float64, bool) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(v.Uint()), true
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	default:
		return 0, false
	}
}

// BeGreaterThan checks if the asserted numeric value is greater than the expected value.
// Panics if types are not comparable or the assertion fails.
func (a *Assertion) BeGreaterThan(expected any) *Assertion {
	v1, v2 := reflect.ValueOf(a.value), reflect.ValueOf(expected)
	result, ok := compare(v1, v2, ">")
	if !ok {
		failAssertion("BeGreaterThan requires comparable numeric types, got %T and %T", a.value, expected)
	}
	if !result {
		failAssertion("should be greater than %#v, but got %#v", expected, a.value)
	}
	return a
}

// BeLessThan checks if the asserted numeric value is less than the expected value.
// Panics if types are not comparable or the assertion fails.
func (a *Assertion) BeLessThan(expected any) *Assertion {
	v1, v2 := reflect.ValueOf(a.value), reflect.ValueOf(expected)
	result, ok := compare(v1, v2, "<")
	if !ok {
		failAssertion("BeLessThan requires comparable numeric types, got %T and %T", a.value, expected)
	}
	if !result {
		failAssertion("should be less than %#v, but got %#v", expected, a.value)
	}
	return a
}

// BeGreaterThanOrEqualTo checks if the asserted numeric value is greater than or equal to the expected value.
// Panics if types are not comparable or the assertion fails.
func (a *Assertion) BeGreaterThanOrEqualTo(expected any) *Assertion {
	v1, v2 := reflect.ValueOf(a.value), reflect.ValueOf(expected)
	result, ok := compare(v1, v2, ">=")
	if !ok {
		failAssertion("BeGreaterThanOrEqualTo requires comparable numeric types, got %T and %T", a.value, expected)
	}
	if !result {
		failAssertion("should be greater than or equal to %#v, but got %#v", expected, a.value)
	}
	return a
}

// BeLessThanOrEqualTo checks if the asserted numeric value is less than or equal to the expected value.
// Panics if types are not comparable or the assertion fails.
func (a *Assertion) BeLessThanOrEqualTo(expected any) *Assertion {
	v1, v2 := reflect.ValueOf(a.value), reflect.ValueOf(expected)
	result, ok := compare(v1, v2, "<=")
	if !ok {
		failAssertion("BeLessThanOrEqualTo requires comparable numeric types, got %T and %T", a.value, expected)
	}
	if !result {
		failAssertion("should be less than or equal to %#v, but got %#v", expected, a.value)
	}
	return a
}

// BeError checks if the asserted value is an error (i.e., implements the error interface and is not nil).
// Panics if the value is not a non-nil error.
func (a *Assertion) BeError() *Assertion {
	if isNil(a.value) {
		failAssertion("should be an error, but got nil")
	}
	if _, ok := a.value.(error); !ok {
		failAssertion("should be an error, but got type %T with value %#v", a.value, a.value)
	}
	// It is an error and it's not nil
	return a
}

// NotBeError checks if the asserted value is nil or not an error.
// Typically used to assert that an error variable is nil.
// Panics if the value is a non-nil error.
func (a *Assertion) NotBeError() *Assertion {
	if !isNil(a.value) {
		if err, ok := a.value.(error); ok {
			failAssertion("should not be an error, but got: %v", err)
		}
		// It's not nil, but also not an error type, which is acceptable for NotBeError.
	}
	// Value is nil or not an error type
	return a
}

// BeErrorOfType checks if the asserted value is an error that matches the type of the target.
// target must be a pointer to a variable of the desired error type (e.g., var target *os.PathError).
// Panics if the value is not an error or does not match the target type.
func (a *Assertion) BeErrorOfType(target any) *Assertion {
	err, ok := a.value.(error)
	if !ok || isNil(err) {
		failAssertion("should be a non-nil error, but got: %#v (type %T)", a.value, a.value)
	}

	// Validate the target: must be a non-nil pointer for errors.As.
	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Ptr || targetVal.IsNil() {
		// errors.As itself will panic with a clearer message if the type is wrong,
		// so we just need to catch the nil/non-pointer case early.
		panic(fmt.Sprintf("internal bumert error: BeErrorOfType target must be a non-nil pointer, got %T", target))
	}

	if !errors.As(err, target) {
		// Attempt to get a meaningful type name for the error message.
		// This is tricky because target could be *MyError, **MyError, *error, etc.
		targetTypeName := "unknown (check target type)"
		if targetVal.Type().Elem().Kind() == reflect.Interface {
			targetTypeName = targetVal.Type().Elem().String() // e.g., "*fs.PathError"
		} else if targetVal.Type().Elem().Elem().Kind() != reflect.Invalid {
			targetTypeName = targetVal.Type().Elem().Elem().String() // e.g., "fs.PathError" for **fs.PathError
		}

		failAssertion("error type should be %s (or wrap it), but got type %T: %v", targetTypeName, err, err)
	}
	return a
}

// BeErrorWithMessage checks if the asserted value is a non-nil error whose message contains the expected substring.
// Panics if the value is not an error or the message does not contain the substring.
func (a *Assertion) BeErrorWithMessage(substring string) *Assertion {
	err, ok := a.value.(error)
	if !ok || isNil(err) {
		failAssertion("should be a non-nil error, but got: %#v (type %T)", a.value, a.value)
	}
	message := err.Error()
	if !strings.Contains(message, substring) {
		failAssertion("error message %#v should contain substring %#v", message, substring)
	}
	return a
}

// --- Direct Assertions ---

// Assert checks if the condition is true.
// Panics with a generic message if the condition is false in debug builds.
func Assert(condition bool) {
	if !condition {
		// Keep the message minimal as we don't have automatic context.
		// Use Assertf for specific messages.
		failAssertion("condition was false") // Skips failAssertion and Assert
	}
}

// Assertf checks if the condition is true.
// Panics with a formatted message if the condition is false in debug builds.
func Assertf(condition bool, format string, args ...any) {
	if !condition {
		failAssertion(format, args...) // Skips failAssertion and Assertf
	}
}
