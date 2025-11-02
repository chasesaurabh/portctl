Packaging & release guide (portctl)

This document explains how releases, packaging, and Homebrew formula generation are intended to work for `portctl`.

Prerequisites
- goreleaser (https://goreleaser.com/) — used for cross-platform builds and packaging
- fpm (optional) — helpful for producing nice .deb/.rpm packages
- gh (GitHub CLI) (optional) — helpful for creating releases

High-level flow

1. Prepare a release branch and update `CHANGELOG.md` (if present) and the `README.md`/`LICENSE`.
2. Tag a release: `git tag -a vX.Y.Z -m "release vX.Y.Z" && git push --tags`
3. CI release workflow (GitHub Actions) triggers on tag pushes and runs goreleaser to produce artifacts and publish to GitHub Releases.
4. After GitHub release is published, the Homebrew tap (separate repo) can be updated with the new formula pointing to the release tarball.

Goreleaser notes

- The included `.goreleaser.yml` builds artifacts for `darwin` and `linux` across `amd64` and `arm64` and generates archives and checksums.
- To test locally: `goreleaser --snapshot --skip-publish --rm-dist`

Deb and RPM helper scripts

- `pack/build_deb.sh` — builds a .deb using `fpm` when available, otherwise uses `dpkg-deb` fallback.
- `pack/build_rpm.sh` — builds an .rpm using `fpm`, otherwise creates a tarball and uses `rpmbuild -ta` as a fallback.

Homebrew tap

- The Homebrew tap is a separate repository (e.g., `github.com/<user>/homebrew-portctl`). The tap contains a single formula `Formula/portctl.rb` which downloads the GitHub release tarball and installs the binary.
- The `goreleaser` config can be adjusted to also auto-generate and update the Homebrew formula, but I recommend managing the tap separately in its own repo for easier review.

Local development

- Build locally: `make build`
- Run unit tests: `make test`
- Run integration tests: `make integration`

CI notes

- The CI workflow runs unit tests and the integration driver on Linux and macOS runners.
- The release workflow triggers goreleaser on tag pushes; ensure `README.md` and `LICENSE` exist before running goreleaser so archives include them.

Troubleshooting

- If goreleaser fails due to missing files, ensure `README.md` and `LICENSE` are present.
- For packaging failures, install `fpm` and `rpmbuild`/`dpkg-deb` on your runner.

Contact

Open an issue or PR in this repository if you need help with packaging specifics.
