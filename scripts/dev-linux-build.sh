#!/usr/bin/env bash
set -euo pipefail

if [[ "$(uname -s)" != "Linux" ]]; then
  echo "error: scripts/dev-linux-build.sh must run inside a Linux container" >&2
  exit 1
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
required_n64_version="v0.1.2"

export PATH="$(go env GOPATH)/bin:$PATH"

current_n64_version="$(
  env GOTOOLCHAIN=local GOOS= GOARCH= GOFLAGS= GOENV= \
    go list -m -f '{{.Version}}' github.com/clktmr/n64
)"

if [[ "$current_n64_version" != "$required_n64_version" ]]; then
  echo "error: expected github.com/clktmr/n64 ${required_n64_version}, got ${current_n64_version}" >&2
  echo "error: run this script from the checkout/worktree that already contains the dependency bump" >&2
  exit 1
fi

if ! command -v go1.24.5-embedded >/dev/null 2>&1; then
  go install github.com/embeddedgo/dl/go1.24.5-embedded@latest
  hash -r
fi

if ! go1.24.5-embedded version >/dev/null 2>&1; then
  go1.24.5-embedded download
fi

if ! command -v n64go >/dev/null 2>&1; then
  go install github.com/clktmr/n64/tools/n64go@v0.1.2
  hash -r
fi

export GOENV="$repo_root/go.env"

go env GOTOOLCHAIN GOOS GOARCH GOFLAGS

cd "$repo_root/examples"

for example_dir in */; do
  (
    cd "$example_dir"
    go1.24.5-embedded build -o game.elf .
    n64go rom game.elf
  )
done
