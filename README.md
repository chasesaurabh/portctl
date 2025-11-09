# portctl

portctl is a small cross-platform CLI to discover and (optionally) kill processes listening on a TCP port.

Quick install (development)

```bash
# build locally
make build

# run help
./portctl --help
```

Basic usage

- Dry-run: show processes listening on port 8080

```bash
./portctl 8080 --dry-run
```

- Force kill with SIGKILL after a 5s grace period

```bash
./portctl 8080 --signal KILL --timeout 5s --force --kill-after
```

Flags

- `--signal` (signal to send; available: TERM, KILL, INT, HUP)
- `--dry-run` (show targets, do not send signals)
- `--force` (do not prompt, force action)
- `--timeout` (wait for graceful shutdown, e.g. 5s)
- `--kill-after` (send SIGKILL after timeout if process still exists)
- `--verbose` (verbose output)
- `--format` (output format: text|json|csv)

Note: The port is a required positional argument (not a flag).

Formats

- Human text (default)
- JSON: `--format json`
- CSV: `--format csv`

Contributing

See `docs/packaging.md` for packaging and release workflows.
