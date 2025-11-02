Integration tests
=================

These tests require `go` and either the built `cmd/portctl/portctl` binary or the ability to `go run` the CLI.

Run locally:

```bash
cd tests/integration
./run_tests.sh
```

Notes:
- Tests attempt to start a `python3 -m http.server` on an ephemeral port and then exercise `portctl` in dry-run and forced-kill mode.
- CI runners must have `python3`, `go`, and `lsof`/`ss` available for the detection steps.
