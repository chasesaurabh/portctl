#!/usr/bin/env bash
# small integration test: start a Python http.server, run portctl in dry-run and force modes
set -euo pipefail

tmpdir=$(mktemp -d)
cd "$tmpdir"

# start a background server on an ephemeral port
python3 -m http.server 0 &
server_pid=$!
# give it a moment
sleep 0.5

# find the port (use lsof or ss)
if command -v ss >/dev/null 2>&1; then
  port_line=$(ss -ltnp | grep "python3")
else
  port_line=$(lsof -nP -iTCP -sTCP:LISTEN | grep "python3" | head -n1)
fi

if [ -z "$port_line" ]; then
  echo "failed to detect server port"
  kill "$server_pid" || true
  exit 1
fi

# extract port
port=$(echo "$port_line" | sed -E 's/.*:([0-9]+) .*/\1/')

echo "server pid=$server_pid port=$port"

# dry-run should not kill
../../cmd/portctl/portctl -dry-run "$port" || true
if kill -0 "$server_pid" 2>/dev/null; then
  echo "dry-run preserved server: PASS"
else
  echo "dry-run killed server: FAIL"
  exit 1
fi

# force kill using the built binary (if exists), else use runCLI via go run
if [ -x "../../cmd/portctl/portctl" ]; then
  ../../cmd/portctl/portctl -force -timeout 3s -kill-after "$port"
else
  (cd ../.. && go run ./cmd/portctl -force -timeout 3s -kill-after "$port")
fi

sleep 1
if kill -0 "$server_pid" 2>/dev/null; then
  echo "server still running after kill: FAIL"
  exit 1
else
  echo "server terminated: PASS"
fi

# cleanup
rm -rf "$tmpdir"
exit 0
#!/usr/bin/env bash
# small integration test: start a Python http.server, run portctl in dry-run and force modes
set -euo pipefail

tmpdir=$(mktemp -d)
cd "$tmpdir"

# start a background server on an ephemeral port
python3 -m http.server 0 &
server_pid=$!
# give it a moment
sleep 0.5

# find the port (use lsof or ss)
if command -v ss >/dev/null 2>&1; then
  port_line=$(ss -ltnp | grep "python3")
else
  port_line=$(lsof -nP -iTCP -sTCP:LISTEN | grep "python3" | head -n1)
fi

if [ -z "$port_line" ]; then
  echo "failed to detect server port"
  kill "$server_pid" || true
  exit 1
fi

# extract port
port=$(echo "$port_line" | sed -E 's/.*:([0-9]+) .*/\1/')

echo "server pid=$server_pid port=$port"

# dry-run should not kill
../../cmd/portctl/portctl -dry-run "$port" || true
if kill -0 "$server_pid" 2>/dev/null; then
  echo "dry-run preserved server: PASS"
else
  echo "dry-run killed server: FAIL"
  exit 1
fi

# force kill using the built binary (if exists), else use runCLI via go run
if [ -x "../../cmd/portctl/portctl" ]; then
  ../../cmd/portctl/portctl -force -timeout 3s -kill-after "$port"
else
  (cd ../.. && go run ./cmd/portctl -force -timeout 3s -kill-after "$port")
fi

sleep 1
if kill -0 "$server_pid" 2>/dev/null; then
  echo "server still running after kill: FAIL"
  exit 1
else
  echo "server terminated: PASS"
fi

# cleanup
rm -rf "$tmpdir"
exit 0
#!/usr/bin/env bash
# small integration test: start a Python http.server, run portctl in dry-run and force modes
set -euo pipefail

tmpdir=$(mktemp -d)
cd "$tmpdir"

# start a background server on an ephemeral port
python3 -m http.server 0 &
server_pid=$!
# give it a moment
sleep 0.5

# find the port (use lsof or ss)
if command -v ss >/dev/null 2>&1; then
  port_line=$(ss -ltnp | grep "python3")
else
  port_line=$(lsof -nP -iTCP -sTCP:LISTEN | grep "python3" | head -n1)
fi

if [ -z "$port_line" ]; then
  echo "failed to detect server port"
  kill "$server_pid" || true
  exit 1
fi

# extract port
port=$(echo "$port_line" | sed -E 's/.*:([0-9]+) .*/\1/')

echo "server pid=$server_pid port=$port"

# dry-run should not kill
../../cmd/portctl/portctl -dry-run "$port" || true
if kill -0 "$server_pid" 2>/dev/null; then
  echo "dry-run preserved server: PASS"
else
  echo "dry-run killed server: FAIL"
  exit 1
fi

# force kill using the built binary (if exists), else use runCLI via go run
if [ -x "../../cmd/portctl/portctl" ]; then
  ../../cmd/portctl/portctl -force -timeout 3s -kill-after "$port"
else
  (cd ../.. && go run ./cmd/portctl -force -timeout 3s -kill-after "$port")
fi

sleep 1
if kill -0 "$server_pid" 2>/dev/null; then
  echo "server still running after kill: FAIL"
  exit 1
else
  echo "server terminated: PASS"
fi

# cleanup
rm -rf "$tmpdir"
exit 0
