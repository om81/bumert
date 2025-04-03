# bumert

`bumert` is a simple, fluent assertion library for Go.

**Core Feature:** Assertions only execute and are only compiled into your binary when the `debug` build tag (or the `bumert` build tag) is specified. In standard/production builds (when the tag is _not_ provided), the assertion calls become no-ops, incurring zero runtime overhead.

## Installation

```bash
go get github.com/deblasis/bumert
```

## Usage

### Fluent API (Recommended)

Use the `Should` function for a fluent interface:

```go
package main

import (
	"fmt"
	"github.com/deblasis/bumert"
)

func process(data *string, count int) {
	// These assertions only run if compiled with `-tags debug`
	bumert.Should(data).NotBeNil() // Checks if pointer is non-nil
	bumert.Should(*data).NotBeEmpty() // Checks if string value is non-empty
	bumert.Should(count).BeGreaterThan(0)

	fmt.Println("Processing:", *data, "Count:", count)
}

func main() {
	val := "some data"
	process(&val, 5)

	var nilData *string
	// This would panic if compiled with `-tags debug`:
	// bumert.Should(nilData).NotBeNil()

	// This would panic if compiled with `-tags debug`:
	// bumert.Should(0).BeGreaterThan(0)
}
```

### Direct Assertion Functions (Alias Package)

You can optionally use the `should` sub-package for direct assertion functions:

```go
package main

import (
	"fmt"
	"github.com/deblasis/bumert/should"
)

func process(data *string, count int) {
	// These assertions only run if compiled with `-tags debug`
	should.ShouldNotBeNil(data)
	should.ShouldNotBeEmpty(*data)
	should.ShouldBeGreaterThan(count, 0)

	fmt.Println("Processing:", *data, "Count:", count)
}

// ... main function similar to above ...
```

## Conditional Compilation: Enabling Assertions (`debug` Tag)

`bumert` relies on Go's build tag system. By default (without any special tags), all assertion calls are effectively removed during compilation. They have **zero performance impact** on your production builds.

To **enable** assertions, you must include the `debug` build tag:

- **Testing:**
  - Use the Makefile target: `make test-debug`
  - Or run directly: `go test -tags debug ./...`
- **Running:**
  - `go run -tags debug main.go`
- **Building:**
  - `go build -tags debug -o myapp_debug main.go`

**Running/Testing without the tag** (e.g., `go test ./...` or `make test` or `go build main.go`) produces a binary where all `bumert` checks are **completely absent**.

## Panic Behavior

When the `debug` tag is enabled, a **failed assertion will cause a panic**. This is intentional and follows the common pattern for assertion libraries in Go (and other languages). It immediately halts execution at the point of the violated assumption during development or debugging, making issues easy to spot. The panic message includes the file and line number of the failed assertion.

In production builds (without the `debug` tag), since the assertions are compiled out, they **cannot** panic.

## Makefile Targets

```bash
# Run tests with assertions enabled
make test-debug

# Run tests without assertions (production mode)
make test

# Format code
make fmt

# Tidy dependencies
make tidy

# Regenerate release stubs (needed if you modify debug assertions)
make generate
```

## Available Assertions

### Direct Assertions (Top Level)

These function directly at the package level:

- `Assert(condition bool)`: Panics with a generic message (`condition was false`) if `condition` is false. For more informative messages, prefer `Assertf`.
- `Assertf(condition bool, format string, args ...any)`: Panics with a formatted message if `condition` is false. **Recommended** over `Assert` for providing context.

### Fluent Assertions (`Should`)

These are chained off the `Should(value)` function:

#### General

- `Should(value).BeNil()` / `should.ShouldBeNil(value)`
- `Should(value).NotBeNil()` / `should.ShouldNotBeNil(value)`
- `Should(value).BeEqual(expected)` / `should.ShouldBeEqual(value, expected)` (uses `reflect.DeepEqual`)
- `Should(value).NotBeEqual(unexpected)` / `should.ShouldNotBeEqual(value, unexpected)` (uses `reflect.DeepEqual`)
- `Should(value).BeZero()` / `should.ShouldBeZero(value)` (checks for the type's zero value)
- `Should(value).NotBeZero()` / `should.ShouldNotBeZero(value)`

### Booleans

- `Should(value).BeTrue()` / `should.ShouldBeTrue(value)`
- `Should(value).BeFalse()` / `should.ShouldBeFalse(value)`

### Collections (Slices, Arrays, Maps, Strings, Channels)

- `Should(value).BeEmpty()` / `should.ShouldBeEmpty(value)` (nil or zero length)
- `Should(value).NotBeEmpty()` / `should.ShouldNotBeEmpty(value)`
- `Should(value).HaveLen(expectedLen)` / `should.ShouldHaveLen(value, expectedLen)`

### Slices/Arrays/Strings

- `Should(value).Contain(element)` / `should.ShouldContain(value, element)` (uses `reflect.DeepEqual` for elements, `strings.Contains` for strings)
- `Should(value).NotContain(element)` / `should.ShouldNotContain(value, element)`

### Strings

- `Should(value).ContainSubstring(substring)` / `should.ShouldContainSubstring(value, substring)`
- `Should(value).HavePrefix(prefix)` / `should.ShouldHavePrefix(value, prefix)`
- `Should(value).HaveSuffix(suffix)` / `should.ShouldHaveSuffix(value, suffix)`

### Numbers (Integers, Floats)

- `Should(value).BeGreaterThan(expected)` / `should.ShouldBeGreaterThan(value, expected)`
- `Should(value).BeLessThan(expected)` / `should.ShouldBeLessThan(value, expected)`
- `Should(value).BeGreaterThanOrEqualTo(expected)` / `should.ShouldBeGreaterThanOrEqualTo(value, expected)`
- `Should(value).BeLessThanOrEqualTo(expected)` / `should.ShouldBeLessThanOrEqualTo(value, expected)`

### Errors

- `Should(value).BeError()` / `should.ShouldBeError(value)` (checks `err != nil` and is `error` type)
- `Should(value).NotBeError()` / `should.ShouldNotBeError(value)` (checks `err == nil` or is not `error` type)
- `Should(value).BeErrorOfType(target)` / `should.ShouldBeErrorOfType(value, target)` (uses `errors.As`, `target` must be `*MyErrorType`)
- `Should(value).BeErrorWithMessage(substring)` / `should.ShouldBeErrorWithMessage(value, substring)`

### Functions

- `Should(value).TrueFn(func() bool)` / `should.ShouldTrueFn(f)` (Asserts the function `f` returns true. The initial `value` is ignored.)
