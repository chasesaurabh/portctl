package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chasesaurabh/portctl/pkg/discover"
	"github.com/chasesaurabh/portctl/pkg/kill"
)

func runCLI(args []string) int {
	cfg := defaultConfig()

	fs := flag.NewFlagSet("portctl", flag.ContinueOnError)
	portFlag := fs.String("port", "", "TCP port to operate on (alternative to positional argument)")
	signal := fs.String("signal", cfg.Signal, "Signal to send (name or number)")
	dryRun := fs.Bool("dry-run", cfg.DryRun, "Show targets but do not send signals")
	force := fs.Bool("force", cfg.Force, "Do not prompt; force action")
	timeout := fs.Duration("timeout", cfg.Timeout, "Seconds to wait for graceful shutdown (e.g., 5s)")
	killAfter := fs.Bool("kill-after", cfg.KillAfter, "Send SIGKILL after timeout if process still exists")
	verbose := fs.Bool("verbose", cfg.Verbose, "Verbose output")
	format := fs.String("format", cfg.Format, "Output format: text|json|csv")

	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "error parsing flags:", err)
		return 1
	}

       // Accept port from --port or positional argument
       if *portFlag != "" {
	       cfg.Port = *portFlag
       } else if fs.NArg() >= 1 {
	       cfg.Port = fs.Arg(0)
       } else {
	       fmt.Fprintln(os.Stderr, "Usage: portctl [--port PORT] [flags] PORT")
	       fs.Usage()
	       return 1
       }
	cfg.Signal = *signal
	cfg.DryRun = *dryRun
	cfg.Force = *force
	cfg.Timeout = *timeout
	cfg.KillAfter = *killAfter
	cfg.Verbose = *verbose
	cfg.Format = *format

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "config: %+v\n", cfg)
	}

	listeners, err := discover.Discover(cfg.Port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "discover error:", err)
		return 2
	}
	// Simple output for now
	for _, l := range listeners {
		fmt.Printf("PID=%d CMD=%s ADDR=%s\n", l.PID, l.Cmd, l.Address)
	}

	if cfg.DryRun {
		return 0
	}

	// Ask for confirmation unless forced
	if !cfg.Force {
		var resp string
		fmt.Fprint(os.Stderr, "Proceed to send signal(s)? [y/N]: ")
		if _, err := fmt.Scanln(&resp); err != nil {
			fmt.Fprintln(os.Stderr, "input error:", err)
			return 1
		}
		if resp != "y" && resp != "Y" {
			fmt.Fprintln(os.Stderr, "aborted by user")
			return 0
		}
	}

	// Send signals
	for _, l := range listeners {
		if err := kill.SendSignal(l.PID, cfg.Signal); err != nil {
			fmt.Fprintf(os.Stderr, "failed to send signal to %d: %v\n", l.PID, err)
			// continue to next
			continue
		}
		if err := kill.WaitForExit(l.PID, cfg.Timeout); err != nil {
			fmt.Fprintf(os.Stderr, "pid %d did not exit: %v\n", l.PID, err)
			if cfg.KillAfter {
				if err := kill.SendSignal(l.PID, "KILL"); err != nil {
					fmt.Fprintf(os.Stderr, "failed to send KILL to %d: %v\n", l.PID, err)
				}
			}
		}
	}

	// Re-check
	newListeners, _ := discover.Discover(cfg.Port)
	if len(newListeners) == 0 {
		fmt.Println("port freed")
		return 0
	}
	fmt.Fprintln(os.Stderr, "some processes still listening")
	return 3
}

	// main is defined in cmd/portctl/main.go; this package provides runCLI
