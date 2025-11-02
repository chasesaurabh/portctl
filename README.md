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
./portctl --port 8080 --dry-run
```

- Force kill with SIGKILL after a 5s grace period

```bash
./portctl --port 8080 --signal KILL --timeout 5s --force
```

Formats

- Human text (default)
- JSON: `--format json`
- CSV: `--format csv`

Contributing

See `docs/packaging.md` for packaging and release workflows.
