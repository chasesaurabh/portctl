#!/usr/bin/env bash
# pack/build_rpm.sh
# Small helper to build an RPM for portctl using fpm or rpmbuild fallback.
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUTDIR="$HERE/dist"
PKGNAME="portctl"
VERSION="${1:-0.0.0}"
ARCH="x86_64"

mkdir -p "$OUTDIR"

# Build the Go binary if not present
if [ ! -f "$HERE/portctl" ]; then
  echo "Building portctl binary..."
  (cd "$HERE" && go build -o portctl ./cmd/portctl)
fi

if command -v fpm >/dev/null 2>&1; then
  echo "Using fpm to build .rpm"
  fpm -s dir -t rpm \
    -n "$PKGNAME" \
    -v "$VERSION" \
    --architecture "$ARCH" \
    --description "portctl — discover and kill processes binding a TCP port" \
    --maintainer "Saurabh Chase <chasesaurabh@gmail.com>" \
    --license MIT \
    --url "https://github.com/chasesaurabh/portctl" \
    --vendor "chasesaurabh" \
    --prefix /usr/local/bin \
    "$HERE/portctl"=/usr/local/bin/portctl \
    -C "$HERE"
  mv *.rpm "$OUTDIR"/
else
  echo "fpm not found; using rpmbuild fallback"
  SRCDIR="$OUTDIR/src"
  mkdir -p "$SRCDIR"
  TAR="$SRCDIR/${PKGNAME}-${VERSION}.tar.gz"
  cp "$HERE/portctl" "$SRCDIR/"
  (cd "$SRCDIR" && tar czf "$TAR" .)
  rpmbuild -ta "$TAR" --define "_topdir $OUTDIR/rpmbuild"
  echo "Look in $OUTDIR/rpmbuild for RPM artifacts"
fi

echo "Built packages in $OUTDIR"
