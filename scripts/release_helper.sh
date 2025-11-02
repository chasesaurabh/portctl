#!/usr/bin/env bash
set -euo pipefail

# Helper to create a Git tag and run goreleaser locally (snapshot mode by default)
if [ "$#" -lt 1 ]; then
  echo "usage: $0 <version> [--snapshot]"
  exit 1
fi

version="$1"
mode="--snapshot"
if [ "${2:-}" = "--release" ]; then
  mode=""
fi

echo "Creating tag v${version} and running goreleaser ${mode}"

git tag -a "v${version}" -m "release v${version}"

goreleaser release ${mode} --rm-dist
