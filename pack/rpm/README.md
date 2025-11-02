RPM packaging helpers

This directory contains simple helpers to build an RPM package for `portctl`.

Prerequisites
- go (for building the binary)
- rpmbuild (for the fallback)
- fpm (optional, recommended) — https://github.com/jordansissel/fpm

Usage

Build an RPM using the helper script:

```bash
# from repository root
pack/build_rpm.sh 0.1.0
```

If `fpm` is installed the script will use it to produce a proper .rpm. Otherwise it will create a tarball and invoke `rpmbuild -ta` as a simple fallback.

Files
- `portctl.spec` — minimal RPM spec file used by `rpmbuild` fallback. Adjust `Summary`, `License`, and `BuildArch` as needed.
- `build_rpm.sh` — build helper script.
