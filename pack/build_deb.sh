#!/usr/bin/env bash
# pack/build_deb.sh
# Small helper to build a .deb for portctl. Prefers fpm if available, falls back to simple dpkg-deb layout.
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUTDIR="$HERE/dist"
PKGNAME="portctl"
VERSION="${1:-0.0.0}"
ARCH="amd64"

mkdir -p "$OUTDIR"/pkg

# Build the Go binary if not present
if [ ! -f "$HERE/portctl" ]; then
  echo "Building portctl binary..."
  (cd "$HERE" && go build -o portctl ./cmd/portctl)
fi

if command -v fpm >/dev/null 2>&1; then
  echo "Using fpm to build .deb"
  fpm -s dir -t deb \
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
  mv *.deb "$OUTDIR"/
else
  echo "fpm not found; using dpkg-deb fallback"
  PKGDIR="$OUTDIR/pkg/${PKGNAME}_${VERSION}_${ARCH}"
  mkdir -p "$PKGDIR"/DEBIAN
  mkdir -p "$PKGDIR"/usr/local/bin
  cp "$HERE/portctl" "$PKGDIR"/usr/local/bin/portctl
  cp "$HERE/pack/deb/control" "$PKGDIR"/DEBIAN/control
  dpkg-deb --build "$PKGDIR"
  mv "${PKGDIR}.deb" "$OUTDIR/"
fi

echo "Built packages in $OUTDIR"
