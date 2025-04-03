# bumert

`bumert` is a simple, fluent assertion library for Go that does nothing... in production builds.

**Core Feature:** Assertions only execute and are only compiled into your binary when the `debug` build tag (or the `bumert` build tag) is specified. In standard/production builds (when the tag is _not_ provided), the assertion calls become no-ops, incurring zero runtime overhead.

> Yes I know, The name could have been **arsert**

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
	// These would panic if compiled with `-tags debug`
	bumert.Should(nilData).NotBeNil()
	bumert.Should(0).BeGreaterThan(0)
	// but they would be no-ops in production builds
	// and/or if no build tags are provided
}
```

### Legacy API (Recommended for oldtimers)

Use the `Assert` function directly:

```go
bumert.Assert(false)
```

or to provide a custom message:

```go
truth := stocksOnlyGoUp()
bumert.Assertf(truth, "Stocks only go up, right?")
```

## Conditional Compilation: Enabling Assertions (`debug` Tag)

`bumert` relies on Go's build tag system. By default (without any special tags), all assertion calls are effectively removed during compilation. They have basically **zero performance impact** on your production builds.

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

When the `debug` tag is enabled, a **failed assertion will cause a panic**. This is intentional and follows the common pattern for assertion libraries in Go (and other languages). It immediately halts execution at the point of the violated assumption during development or debugging, making issues easy to spot. The panic message includes the file and line number of the failed assertion, along with the failure reason (e.g., "Expected <foo> to be greater than <bar>").

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

- `Should(value).BeNil()`
- `Should(value).NotBeNil()`
- `Should(value).BeEqual(expected)` (uses `reflect.DeepEqual`)
- `Should(value).NotBeEqual(unexpected)` (uses `reflect.DeepEqual`)
- `Should(value).BeZero()` (checks against the type's zero value: `0`, `""`, `nil`, `false`, etc.)
- `Should(value).NotBeZero()` (checks if the value is not the type's zero value)

### Booleans

- `Should(value).BeTrue()`
- `Should(value).BeFalse()`

### Collections (Slices, Arrays, Maps, Strings, Channels)

- `Should(value).BeEmpty()` (nil or zero length)
- `Should(value).NotBeEmpty()`
- `Should(value).HaveLen(expectedLen)`

### Slices/Arrays/Strings

- `Should(value).Contain(element)` (uses `reflect.DeepEqual` for elements, `strings.Contains` for strings)
- `Should(value).NotContain(element)`

### Strings

- `Should(value).ContainSubstring(substring)`
- `Should(value).HavePrefix(prefix)`
- `Should(value).HaveSuffix(suffix)`

### Numbers (Integers, Floats)

- `Should(value).BeGreaterThan(expected)`
- `Should(value).BeLessThan(expected)`
- `Should(value).BeGreaterThanOrEqualTo(expected)`
- `Should(value).BeLessThanOrEqualTo(expected)`

### Errors

- `Should(value).BeError()` (checks `err != nil` and is `error` type)
- `Should(value).NotBeError()` (checks `err == nil` or is not `error` type)
- `Should(value).BeErrorOfType(target)` (uses `errors.As`, `target` must be `*MyErrorType`)
- `Should(value).BeErrorWithMessage(substring)`

### Functions

- `Should(value).TrueFn(func() bool)` (Asserts the function `f` returns true. The initial `value` is ignored.)
