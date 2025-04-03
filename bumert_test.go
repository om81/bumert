//go:build debug || bumert

package bumert_test

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/deblasis/bumert"
)

// assertPanics asserts that the given function f panics with a message containing expectedSubstring.
func assertPanics(t *testing.T, f func(), expectedSubstring string) {
	t.Helper()
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Errorf("Expected function to panic, but it did not")
			return
		}
		panicMsg := fmt.Sprintf("%v", recovered)
		if expectedSubstring != "" && !strings.Contains(panicMsg, expectedSubstring) {
			t.Errorf("Panic message did not contain expected substring.\nExpected: %q\nActual:   %q", expectedSubstring, panicMsg)
		}
	}()
	f()
}

// assertNotPanics asserts that the given function f does not panic.
func assertNotPanics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Errorf("Expected function not to panic, but it did: %v", recovered)
		}
	}()
	f()
}

// --- Test Cases ---

func TestAssertion_BeNil(t *testing.T) {
	var typedNil *int = nil
	var nonNilInt = 5

	tests := []struct {
		name        string
		value       any
		shouldPanic bool
		panicSubstr string
	}{
		{"Untyped nil", nil, false, ""},
		{"Typed nil pointer", typedNil, false, ""},
		{"Nil interface", interface{}(nil), false, ""},
		{"Non-nil pointer", &nonNilInt, true, "should be nil"},
		{"Non-nil value", 42, true, "should be nil"},
		{"Empty string", "", true, "should be nil"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).BeNil() }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_NotBeNil(t *testing.T) {
	var typedNil *int = nil
	var nonNilInt = 5

	tests := []struct {
		name        string
		value       any
		shouldPanic bool
		panicSubstr string
	}{
		{"Untyped nil", nil, true, "should not be nil"},
		{"Typed nil pointer", typedNil, true, "should not be nil"},
		{"Nil interface", interface{}(nil), true, "should not be nil"},
		{"Non-nil pointer", &nonNilInt, false, ""},
		{"Non-nil value", 42, false, ""},
		{"Empty string", "", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).NotBeNil() }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_TrueFn(t *testing.T) {
	tests := []struct {
		name        string
		f           func() bool
		shouldPanic bool
		panicSubstr string
	}{
		{"Function returns true", func() bool { return true }, false, ""},
		{"Function returns false", func() bool { return false }, true, "function returned false"},
		{"Function checks variable (true)", func() bool { x := 5; return x > 0 }, false, ""},
		{"Function checks variable (false)", func() bool { x := -1; return x > 0 }, true, "function returned false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The initial Should(nil) value doesn't matter for TrueFn
			f := func() { bumert.Should(nil).TrueFn(tt.f) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_BeTrue(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		shouldPanic bool
		panicSubstr string
	}{
		{"Boolean true", true, false, ""},
		{"Boolean false", false, true, "should be true"},
		{"Non-boolean (int)", 1, true, "should be true"},
		{"Non-boolean (string)", "true", true, "should be true"},
		{"Nil", nil, true, "should be true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).BeTrue() }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_BeFalse(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		shouldPanic bool
		panicSubstr string
	}{
		{"Boolean false", false, false, ""},
		{"Boolean true", true, true, "should be false"},
		{"Non-boolean (int)", 0, true, "should be false"},
		{"Non-boolean (string)", "false", true, "should be false"},
		{"Nil", nil, true, "should be false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).BeFalse() }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

type comparableStruct struct {
	A int
	B string
}

func TestAssertion_BeEqual(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		expected    any
		shouldPanic bool
		panicSubstr string
	}{
		{"Equal ints", 10, 10, false, ""},
		{"Unequal ints", 10, 11, true, "should be equal"},
		{"Equal floats", 3.14, 3.14, false, ""},
		{"Unequal floats", 3.14, 3.15, true, "should be equal"},
		{"Equal strings", "hello", "hello", false, ""},
		{"Unequal strings", "hello", "world", true, "should be equal"},
		{"Equal structs", comparableStruct{1, "a"}, comparableStruct{1, "a"}, false, ""},
		{"Unequal structs", comparableStruct{1, "a"}, comparableStruct{2, "a"}, true, "should be equal"},
		{"Equal slices", []int{1, 2}, []int{1, 2}, false, ""},
		{"Unequal slices (value)", []int{1, 2}, []int{1, 3}, true, "should be equal"},
		{"Unequal slices (length)", []int{1, 2}, []int{1, 2, 3}, true, "should be equal"},
		{"Equal maps", map[string]int{"a": 1}, map[string]int{"a": 1}, false, ""},
		{"Unequal maps (value)", map[string]int{"a": 1}, map[string]int{"a": 2}, true, "should be equal"},
		{"Unequal maps (key)", map[string]int{"a": 1}, map[string]int{"b": 1}, true, "should be equal"},
		{"Nil vs Empty Slice", nil, []int{}, true, "should be equal"}, // DeepEqual distinguishes nil from empty
		{"Nil vs Nil", nil, nil, false, ""},
		{"Typed Nil vs Typed Nil", (*int)(nil), (*int)(nil), false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).BeEqual(tt.expected) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_NotBeEqual(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		unexpected  any
		shouldPanic bool
		panicSubstr string
	}{
		{"Unequal ints", 10, 11, false, ""},
		{"Equal ints", 10, 10, true, "should not be equal"},
		{"Unequal floats", 3.14, 3.15, false, ""},
		{"Equal floats", 3.14, 3.14, true, "should not be equal"},
		{"Unequal strings", "hello", "world", false, ""},
		{"Equal strings", "hello", "hello", true, "should not be equal"},
		{"Unequal structs", comparableStruct{1, "a"}, comparableStruct{2, "a"}, false, ""},
		{"Equal structs", comparableStruct{1, "a"}, comparableStruct{1, "a"}, true, "should not be equal"},
		{"Unequal slices (value)", []int{1, 2}, []int{1, 3}, false, ""},
		{"Equal slices", []int{1, 2}, []int{1, 2}, true, "should not be equal"},
		{"Nil vs Empty Slice", nil, []int{}, false, ""},
		{"Nil vs Nil", nil, nil, true, "should not be equal"},
		{"Typed Nil vs Typed Nil", (*int)(nil), (*int)(nil), true, "should not be equal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).NotBeEqual(tt.unexpected) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_BeEmpty(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		shouldPanic bool
		panicSubstr string
	}{
		{"Untyped nil", nil, false, ""},
		{"Typed nil pointer", (*int)(nil), false, ""},
		{"Empty string", "", false, ""},
		{"Empty slice", []int{}, false, ""},
		{"Empty map", map[string]int{}, false, ""},
		{"Empty channel", make(chan int), false, ""},
		{"Nil slice", []int(nil), false, ""},
		{"Nil map", map[string]int(nil), false, ""},
		{"Non-empty string", "a", true, "should be empty"},
		{"Non-empty slice", []int{1}, true, "should be empty"},
		{"Non-empty map", map[string]int{"a": 1}, true, "should be empty"},
		{"Non-length-testable (int)", 0, true, "non-length-testable"},
		{"Non-length-testable (struct)", struct{}{}, true, "non-length-testable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).BeEmpty() }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_NotBeEmpty(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		shouldPanic bool
		panicSubstr string
	}{
		{"Non-empty string", "a", false, ""},
		{"Non-empty slice", []int{1}, false, ""},
		{"Non-empty map", map[string]int{"a": 1}, false, ""},
		{"Non-length-testable (int)", 0, false, ""},             // 0 is considered not empty here
		{"Non-length-testable (struct)", struct{}{}, false, ""}, // non-nil struct is not empty
		{"Untyped nil", nil, true, "should not be empty"},
		{"Typed nil pointer", (*int)(nil), true, "should not be empty"},
		{"Empty string", "", true, "should not be empty"},
		{"Empty slice", []int{}, true, "should not be empty"},
		{"Empty map", map[string]int{}, true, "should not be empty"},
		{"Nil slice", []int(nil), true, "should not be empty"},
		{"Nil map", map[string]int(nil), true, "should not be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).NotBeEmpty() }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_HaveLen(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		expectedLen int
		shouldPanic bool
		panicSubstr string
	}{
		{"Correct length string", "abc", 3, false, ""},
		{"Incorrect length string", "abc", 2, true, "should have length 2, but got 3"},
		{"Correct length slice", []int{1, 2}, 2, false, ""},
		{"Incorrect length slice", []int{1, 2}, 3, true, "should have length 3, but got 2"},
		{"Correct length map", map[int]bool{1: true}, 1, false, ""},
		{"Incorrect length map", map[int]bool{1: true}, 0, true, "should have length 0, but got 1"},
		{"Empty string, len 0", "", 0, false, ""},
		{"Empty slice, len 0", []string{}, 0, false, ""},
		{"Nil slice, incorrect len", []int(nil), 1, true, "should have length 1, but got 0"},
		{"Nil map, incorrect len", map[string]int(nil), 1, true, "should have length 1, but got 0"},
		{"Not length-testable (int)", 123, 3, true, "non-length-testable"},
		{"Not length-testable (struct)", struct{}{}, 1, true, "non-length-testable"},
		{"Nil interface", nil, 0, true, "non-length-testable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).HaveLen(tt.expectedLen) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_Contain(t *testing.T) {
	tests := []struct {
		name            string
		collection      any
		expectedElement any
		shouldPanic     bool
		panicSubstr     string
	}{
		{"String contains substring", "hello world", "lo wor", false, ""},
		{"String does not contain substring", "hello world", "goodbye", true, "should contain substring"},
		{"Slice contains element", []int{1, 2, 3}, 2, false, ""},
		{"Slice does not contain element", []int{1, 2, 3}, 4, true, "should contain element"},
		{"Slice contains struct", []comparableStruct{{1, "a"}, {2, "b"}}, comparableStruct{1, "a"}, false, ""},
		{"Slice does not contain struct", []comparableStruct{{1, "a"}}, comparableStruct{2, "b"}, true, "should contain element"},
		{"Array contains element", [3]string{"a", "b", "c"}, "b", false, ""},
		{"Array does not contain element", [3]string{"a", "b", "c"}, "d", true, "should contain element"},
		{"Empty slice does not contain", []int{}, 1, true, "should contain element"},
		{"Nil slice does not contain", []string(nil), "a", true, "should contain element"},
		{"String contains empty string", "abc", "", false, ""}, // strings.Contains behavior
		{"Unsupported type (int)", 123, 2, true, "requires slice, array, or string"},
		{"Unsupported type (map)", map[int]int{1: 1}, 1, true, "requires slice, array, or string"},
		{"String requires string element", "abc", 1, true, "expected a string element"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.collection).Contain(tt.expectedElement) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_NotContain(t *testing.T) {
	tests := []struct {
		name              string
		collection        any
		unexpectedElement any
		shouldPanic       bool
		panicSubstr       string
	}{
		{"String does not contain substring", "hello world", "goodbye", false, ""},
		{"String contains substring", "hello world", "lo wor", true, "should not contain substring"},
		{"Slice does not contain element", []int{1, 2, 3}, 4, false, ""},
		{"Slice contains element", []int{1, 2, 3}, 2, true, "should not contain element"},
		{"Slice does not contain struct", []comparableStruct{{1, "a"}}, comparableStruct{2, "b"}, false, ""},
		{"Slice contains struct", []comparableStruct{{1, "a"}, {2, "b"}}, comparableStruct{1, "a"}, true, "should not contain element"},
		{"Array does not contain element", [3]string{"a", "b", "c"}, "d", false, ""},
		{"Array contains element", [3]string{"a", "b", "c"}, "b", true, "should not contain element"},
		{"Empty slice does not contain", []int{}, 1, false, ""},
		{"Nil slice does not contain", []string(nil), "a", false, ""},
		{"String contains empty string (panics)", "abc", "", true, "should not contain substring"}, // strings.Contains behavior
		{"Unsupported type (int)", 123, 2, true, "requires slice, array, or string"},
		{"String requires string element", "abc", 1, true, "expected a string element"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.collection).NotContain(tt.unexpectedElement) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_ContainSubstring(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		substring   string
		shouldPanic bool
		panicSubstr string
	}{
		{"Contains", "hello world", "lo wor", false, ""},
		{"Does not contain", "hello world", "goodbye", true, "should contain substring"},
		{"Contains empty", "abc", "", false, ""},
		{"Empty contains non-empty", "", "a", true, "should contain substring"},
		{"Empty contains empty", "", "", false, ""},
		{"Not a string", 123, "1", true, "requires a string value"},
		{"Nil", nil, "a", true, "requires a string value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).ContainSubstring(tt.substring) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_HavePrefix(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		prefix      string
		shouldPanic bool
		panicSubstr string
	}{
		{"Has prefix", "hello world", "hello", false, ""},
		{"Does not have prefix", "hello world", "world", true, "should have prefix"},
		{"Has empty prefix", "abc", "", false, ""},
		{"Empty has non-empty prefix", "", "a", true, "should have prefix"},
		{"Empty has empty prefix", "", "", false, ""},
		{"Not a string", 123, "1", true, "requires a string value"},
		{"Nil", nil, "a", true, "requires a string value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).HavePrefix(tt.prefix) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_HaveSuffix(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		suffix      string
		shouldPanic bool
		panicSubstr string
	}{
		{"Has suffix", "hello world", "world", false, ""},
		{"Does not have suffix", "hello world", "hello", true, "should have suffix"},
		{"Has empty suffix", "abc", "", false, ""},
		{"Empty has non-empty suffix", "", "a", true, "should have suffix"},
		{"Empty has empty suffix", "", "", false, ""},
		{"Not a string", 123, "3", true, "requires a string value"},
		{"Nil", nil, "a", true, "requires a string value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).HaveSuffix(tt.suffix) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_BeZero(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		shouldPanic bool
		panicSubstr string
	}{
		{"Zero int", 0, false, ""},
		{"Non-zero int", 1, true, "should be the zero value"},
		{"Zero float", 0.0, false, ""},
		{"Non-zero float", 0.1, true, "should be the zero value"},
		{"Empty string", "", false, ""},
		{"Non-empty string", "a", true, "should be the zero value"},
		{"False boolean", false, false, ""},
		{"True boolean", true, true, "should be the zero value"},
		{"Nil pointer", (*int)(nil), false, ""},
		{"Non-nil pointer", new(int), true, "should be the zero value"},
		{"Nil slice", []int(nil), false, ""},
		{"Empty slice", []int{}, true, "should be the zero value"}, // Empty is not zero for slice
		{"Nil map", map[string]int(nil), false, ""},
		{"Empty map", map[string]int{}, true, "should be the zero value"}, // Empty is not zero for map
		{"Zero struct", comparableStruct{}, false, ""},
		{"Non-zero struct", comparableStruct{A: 1}, true, "should be the zero value"},
		{"Nil interface", nil, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).BeZero() }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_NotBeZero(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		shouldPanic bool
		panicSubstr string
	}{
		{"Non-zero int", 1, false, ""},
		{"Zero int", 0, true, "should not be the zero value"},
		{"Non-zero float", 0.1, false, ""},
		{"Zero float", 0.0, true, "should not be the zero value"},
		{"Non-empty string", "a", false, ""},
		{"Empty string", "", true, "should not be the zero value"},
		{"True boolean", true, false, ""},
		{"False boolean", false, true, "should not be the zero value"},
		{"Non-nil pointer", new(int), false, ""},
		{"Nil pointer", (*int)(nil), true, "should not be the zero value"},
		{"Empty slice", []int{}, false, ""}, // Empty is not zero
		{"Nil slice", []int(nil), true, "should not be the zero value"},
		{"Empty map", map[string]int{}, false, ""}, // Empty is not zero
		{"Nil map", map[string]int(nil), true, "should not be the zero value"},
		{"Non-zero struct", comparableStruct{A: 1}, false, ""},
		{"Zero struct", comparableStruct{}, true, "should not be the zero value"},
		{"Nil interface", nil, true, "should not be the zero value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).NotBeZero() }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_BeGreaterThan(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		expected    any
		shouldPanic bool
		panicSubstr string
	}{
		{"Int > Int", 10, 9, false, ""},
		{"Int == Int", 10, 10, true, "should be greater than"},
		{"Int < Int", 9, 10, true, "should be greater than"},
		{"Float > Float", 3.15, 3.14, false, ""},
		{"Float == Float", 3.14, 3.14, true, "should be greater than"},
		{"Float < Float", 3.14, 3.15, true, "should be greater than"},
		{"Uint > Uint", uint(10), uint(9), false, ""},
		{"Int > Float (convertible)", 10, 9.5, false, ""},
		{"Float > Int (convertible)", 10.5, 10, false, ""},
		{"Incompatible types (string vs int)", "10", 9, true, "requires comparable numeric types"},
		{"Incompatible types (struct vs int)", struct{}{}, 0, true, "requires comparable numeric types"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).BeGreaterThan(tt.expected) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_BeLessThan(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		expected    any
		shouldPanic bool
		panicSubstr string
	}{
		{"Int < Int", 9, 10, false, ""},
		{"Int == Int", 10, 10, true, "should be less than"},
		{"Int > Int", 10, 9, true, "should be less than"},
		{"Float < Float", 3.14, 3.15, false, ""},
		{"Float == Float", 3.14, 3.14, true, "should be less than"},
		{"Float > Float", 3.15, 3.14, true, "should be less than"},
		{"Uint < Uint", uint(9), uint(10), false, ""},
		{"Int < Float (convertible)", 9, 9.5, false, ""},
		{"Float < Int (convertible)", 9.5, 10, false, ""},
		{"Incompatible types (string vs int)", "9", 10, true, "requires comparable numeric types"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).BeLessThan(tt.expected) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_BeGreaterThanOrEqualTo(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		expected    any
		shouldPanic bool
		panicSubstr string
	}{
		{"Int > Int", 10, 9, false, ""},
		{"Int == Int", 10, 10, false, ""},
		{"Int < Int", 9, 10, true, "should be greater than or equal to"},
		{"Float > Float", 3.15, 3.14, false, ""},
		{"Float == Float", 3.14, 3.14, false, ""},
		{"Float < Float", 3.14, 3.15, true, "should be greater than or equal to"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).BeGreaterThanOrEqualTo(tt.expected) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_BeLessThanOrEqualTo(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		expected    any
		shouldPanic bool
		panicSubstr string
	}{
		{"Int < Int", 9, 10, false, ""},
		{"Int == Int", 10, 10, false, ""},
		{"Int > Int", 10, 9, true, "should be less than or equal to"},
		{"Float < Float", 3.14, 3.15, false, ""},
		{"Float == Float", 3.14, 3.14, false, ""},
		{"Float > Float", 3.15, 3.14, true, "should be less than or equal to"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).BeLessThanOrEqualTo(tt.expected) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_BeError(t *testing.T) {
	var nilErr error = nil
	tests := []struct {
		name        string
		value       any
		shouldPanic bool
		panicSubstr string
	}{
		{"Non-nil error", errors.New("test error"), false, ""},
		{"Nil error variable", nilErr, true, "should be an error, but got nil"},
		{"Untyped nil", nil, true, "should be an error, but got nil"},
		{"Non-error type (int)", 1, true, "should be an error, but got type int"},
		{"Non-error type (string)", "error", true, "should be an error, but got type string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).BeError() }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_NotBeError(t *testing.T) {
	var nilErr error = nil
	tests := []struct {
		name        string
		value       any
		shouldPanic bool
		panicSubstr string
	}{
		{"Nil error variable", nilErr, false, ""},
		{"Untyped nil", nil, false, ""},
		{"Non-error type (int)", 1, false, ""},          // Not an error, so passes
		{"Non-error type (string)", "error", false, ""}, // Not an error, so passes
		{"Non-nil error", errors.New("test error"), true, "should not be an error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).NotBeError() }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertion_BeErrorOfType(t *testing.T) {
	var pathErr *os.PathError // Target type for errors.As
	var customErr *CustomError

	tests := []struct {
		name        string
		value       any // The error value being asserted
		target      any // The target for errors.As (e.g., &pathErr)
		shouldPanic bool
		panicSubstr string
	}{
		{"Correct type (os.PathError)", &os.PathError{Op: "read"}, &pathErr, false, ""},
		{"Wrapped correct type", fmt.Errorf("wrapped: %w", &os.PathError{Op: "write"}), &pathErr, false, ""},
		{"Correct custom type", &CustomError{Code: 1}, &customErr, false, ""},
		{"Incorrect type", errors.New("generic error"), &pathErr, true, "error type should be"},
		{"Nil error value", nil, &pathErr, true, "should be a non-nil error"},
		{"Non-error value", 123, &pathErr, true, "should be a non-nil error"},
		{"Nil target (internal panic)", &os.PathError{}, nil, true, "target must be a non-nil pointer"}, // This tests the internal check
		// {"Non-pointer target (internal panic)", &os.PathError{}, pathErr, true, "target must be a pointer"}, // Cannot create non-pointer variable of interface type easily
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).BeErrorOfType(tt.target) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

type CustomError struct {
	Code int
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("custom error code %d", e.Code)
}

func TestAssertion_BeErrorWithMessage(t *testing.T) {
	tests := []struct {
		name        string
		value       any
		substring   string
		shouldPanic bool
		panicSubstr string
	}{
		{"Message contains substring", errors.New("file not found"), "not found", false, ""},
		{"Message does not contain substring", errors.New("file exists"), "not found", true, "should contain substring"},
		{"Wrapped error contains substring", fmt.Errorf("wrapped: %w", errors.New("inner fail")), "inner fail", false, ""},
		{"Custom error contains substring", &CustomError{Code: 404}, "code 404", false, ""},
		{"Custom error does not contain substring", &CustomError{Code: 500}, "code 404", true, "should contain substring"},
		{"Nil error", nil, "any", true, "should be a non-nil error"},
		{"Non-error value", 123, "any", true, "should be a non-nil error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Should(tt.value).BeErrorWithMessage(tt.substring) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

// --- Direct Assertion Tests ---

func TestAssert(t *testing.T) {
	tests := []struct {
		name        string
		condition   bool
		shouldPanic bool
		panicSubstr string
	}{
		{"Condition true", true, false, ""},
		{"Condition false", false, true, "condition was false"},
		{"Complex condition true", (5 > 3 && "a" == "a"), false, ""},
		{"Complex condition false", (1 == 2 || false), true, "condition was false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Assert(tt.condition) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}

func TestAssertf(t *testing.T) {
	value := 42
	tests := []struct {
		name        string
		condition   bool
		format      string
		args        []any
		shouldPanic bool
		panicSubstr string
	}{
		{"Condition true, no message", true, "this should not panic", nil, false, ""},
		{"Condition false, simple message", false, "value should be positive", nil, true, "value should be positive"},
		{"Condition false, formatted message", false, "value %d is not > 100", []any{value}, true, "value 42 is not > 100"},
		{"Condition true, complex message", true, "value %d > %d", []any{value, 10}, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := func() { bumert.Assertf(tt.condition, tt.format, tt.args...) }
			if tt.shouldPanic {
				assertPanics(t, f, tt.panicSubstr)
			} else {
				assertNotPanics(t, f)
			}
		})
	}
}
