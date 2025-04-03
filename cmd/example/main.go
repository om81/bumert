package main

import (
	"errors"
	"fmt"

	"github.com/deblasis/bumert"
)

// config represents some application configuration.
// We use bumert to assert preconditions during development.
type config struct {
	Host string
	Port int
	User string
	Tags []string
}

// newConfig creates a config, potentially returning an error.
func newConfig(host string, port int, user string) (*config, error) {
	// --- Debug Assertions --- (These only run with `-tags debug` or `-tags bumert`)
	bumert.Assert(false)
	bumert.Should(host).NotBeEmpty()     // Ensure host is provided
	bumert.Should(port).BeGreaterThan(0) // Ensure port is positive
	// --- End Debug Assertions ---

	if host == "" {
		// Real validation still needed for production
		return nil, errors.New("host cannot be empty")
	}
	if port <= 0 {
		return nil, fmt.Errorf("invalid port: %d", port)
	}

	cfg := &config{
		Host: host,
		Port: port,
		User: user,
		Tags: []string{"default", "example"},
	}

	// --- More Debug Assertions ---
	bumert.Should(cfg).NotBeNil()
	bumert.Should(cfg.Tags).HaveLen(2)
	bumert.Should(cfg.Tags).Contain("example")
	// --- End Debug Assertions ---

	return cfg, nil
}

func main() {
	fmt.Println("--- Running Example ---")

	// --- Case 1: Valid Config ---
	fmt.Println("\nAttempting valid configuration...")
	cfg1, err1 := newConfig("localhost", 8080, "admin")

	// --- Debug Assertion on Error --- (Only runs with `-tags debug`)
	bumert.Should(err1).NotBeError() // Assert no error occurred here
	// --- End Debug Assertion ---

	if err1 != nil {
		fmt.Println("  Error creating config 1:", err1)
	} else {
		fmt.Printf("  Config 1 created successfully: %+v\n", cfg1)
	}

	// --- Case 2: Invalid Config (Empty Host) ---
	// This demonstrates how debug assertions can catch issues early.
	fmt.Println("\nAttempting invalid configuration (empty host)...")
	cfg2, err2 := newConfig("", 9000, "guest")

	// --- Debug Assertion on Error --- (Only runs with `-tags debug`)
	bumert.Should(err2).BeError() // Assert an error DID occur here
	// --- End Debug Assertion ---

	if err2 != nil {
		fmt.Println("  Error creating config 2:", err2) // Expected error
	} else {
		fmt.Printf("  Config 2 created unexpectedly: %+v\n", cfg2)
	}

	// --- Case 3: Assertion Failure (if debug enabled) ---
	fmt.Println("\nDemonstrating a failing assertion (if debug enabled)...")
	shouldFailValue := 0
	fmt.Println("  About to assert shouldFailValue > 0...")

	// --- Failing Debug Assertion ---
	// This line will PANIC ONLY when run with `-tags debug`.
	// In production builds, it does nothing.
	bumert.Should(shouldFailValue).BeGreaterThan(0)
	// --- End Failing Debug Assertion ---

	// This line will only be reached if NOT running with `-tags debug`
	fmt.Println("  Assertion passed (or was skipped in production build).")

	fmt.Println("\n--- Example Finished ---")
}
