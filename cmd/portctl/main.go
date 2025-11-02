package main

import (
	"os"
)

func main() {
	// Delegate to runCLI which contains the real CLI logic. Keep main small so
	// tests can call runCLI directly.
	os.Exit(runCLI(os.Args[1:]))
}
