Deb packaging helpers

This directory contains simple helpers to build a Debian (.deb) package for `portctl`.

Prerequisites
- go (for building the binary)
- dpkg-deb (for fallback)
- fpm (optional, recommended) — https://github.com/jordansissel/fpm

Usage

Build a .deb using the helper script:

```bash
# from repository root
pack/build_deb.sh 0.1.0
```

If `fpm` is installed the script will use it to produce a proper .deb. Otherwise it will package a minimal dpkg structure and run `dpkg-deb --build`.

Files
- `control` — minimal DEBIAN/control template used by dpkg-deb fallback. Adjust `Maintainer`, `Depends` and `Architecture` as needed.
- `build_deb.sh` — build helper script.
