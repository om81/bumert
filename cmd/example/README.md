# Bumert Example Usage

This directory contains a simple example Go program demonstrating how to use the `bumert` assertion library.

## Purpose

The example shows:

1.  Using `bumert.Should()` for fluent assertions.
2.  How assertions act as checks during development (when compiled with `-tags debug`).
3.  How assertions are completely skipped in production builds (compiled without the tag).
4.  The panic behavior of a failed assertion in a debug build.

## Running the Example

### Production Mode (Assertions Skipped)

Run without the `debug` tag. Assertions will be compiled out, and the program will run to completion without panicking, even though one assertion is logically incorrect.

```bash
# Navigate to the root of the bumert repository first
cd ../../
go run ./cmd/example/main.go
```

**Expected Output (Production):**

```
--- Running Example ---

Attempting valid configuration...
  Config 1 created successfully: &{Host:localhost Port:8080 User:admin Tags:[default example]}

Attempting invalid configuration (empty host)...
  Error creating config 2: host cannot be empty

Demonstrating a failing assertion (if debug enabled)...
  About to assert shouldFailValue > 0...
  Assertion passed (or was skipped in production build).

--- Example Finished ---
```

### Debug Mode (Assertions Enabled)

Run _with_ the `debug` tag. Assertions will execute. The final assertion (`bumert.Should(0).BeGreaterThan(0)`) will fail and cause the program to panic.

```bash
# Navigate to the root of the bumert repository first
cd ../../
go run -tags debug ./cmd/example/main.go
```

**Expected Output (Debug):**

The program will print output similar to the production run up until the failing assertion, then it will panic. The exact panic message will include the file (`main.go`), line number, and the assertion failure reason.

```
--- Running Example ---

Attempting valid configuration...
  Config 1 created successfully: &{Host:localhost Port:8080 User:admin Tags:[default example]}

Attempting invalid configuration (empty host)...
  Error creating config 2: host cannot be empty

Demonstrating a failing assertion (if debug enabled)...
  About to assert shouldFailValue > 0...
panic: main.go:84: assertion failed: should be greater than 0, but got 0

goroutine 1 [running]:
main.main()
	/home/user/path/to/bumert/cmd/example/main.go:84 +0x...
exit status 2
```

_(Note: The exact file path and line number in the panic message may vary slightly)_
