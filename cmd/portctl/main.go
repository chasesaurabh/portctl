package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Minimal entrypoint for PR 001
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "portctl - discover and free TCP ports\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <port>\n", os.Args[0])
		flag.PrintDefaults()
	}

	dryRun := flag.Bool("dry-run", false, "Show targets but do not send signals")
	signal := flag.String("signal", "TERM", "Signal to send (name or number)")
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	port := flag.Arg(0)
	fmt.Printf("portctl (stub): port=%s dry-run=%v signal=%s\n", port, *dryRun, *signal)
}
