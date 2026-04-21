#!/usr/bin/env bash
set -euo pipefail

if [[ "$(uname -s)" != "Linux" ]]; then
  echo "error: scripts/dev-linux-build.sh must run inside a Linux container or Linux CI runner" >&2
  exit 1
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
required_n64_version="v0.1.2"

clean_go_env() {
  env -u GOENV -u GOOS -u GOARCH -u GOFLAGS -u GOTOOLCHAIN -u GOPATH -u GOBIN "$@"
}

if [[ -f "$repo_root/.envrc" ]]; then
  echo "error: stale .envrc detected; remove it and use go.env only" >&2
  exit 1
fi

stale_example_modules="$(
  find "$repo_root/examples" -mindepth 2 -maxdepth 2 \( -name go.mod -o -name go.sum \) -print | sort
)"

if [[ -n "$stale_example_modules" ]]; then
  echo "error: remove nested example go.mod/go.sum files before building" >&2
  printf '  %s\n' $stale_example_modules >&2
  exit 1
fi

export PATH="/usr/local/go/bin:$PATH"
export PATH="$(clean_go_env /usr/local/go/bin/go env GOPATH)/bin:$PATH"

cd "$repo_root"

current_n64_version="$(
  clean_go_env /usr/local/go/bin/go list -m -f '{{.Version}}' github.com/clktmr/n64
)"

if [[ "$current_n64_version" != "$required_n64_version" ]]; then
  echo "error: expected github.com/clktmr/n64 ${required_n64_version}, got ${current_n64_version}" >&2
  echo "error: run this script from a checkout that already contains the v0.1.2 dependency bump" >&2
  exit 1
fi

if ! command -v go1.24.5-embedded >/dev/null 2>&1; then
  clean_go_env /usr/local/go/bin/go install github.com/embeddedgo/dl/go1.24.5-embedded@latest
  hash -r
fi

if ! clean_go_env go1.24.5-embedded version >/dev/null 2>&1; then
  clean_go_env go1.24.5-embedded download
fi

if ! command -v n64go >/dev/null 2>&1; then
  clean_go_env /usr/local/go/bin/go install github.com/clktmr/n64/tools/n64go@v0.1.2
  hash -r
fi

for example_dir in "$repo_root"/examples/*; do
  [[ -d "$example_dir" ]] || continue
  example_name="$(basename "$example_dir")"
  elf="$example_dir/game.elf"
  rom="$example_dir/game.z64"

  rm -f "$elf" "$rom"

  clean_go_env GOENV="$repo_root/go.env" go1.24.5-embedded build -o "$elf" "./examples/$example_name"
  clean_go_env GOENV="$repo_root/go.env" n64go rom "$elf"

  test -f "$elf"
  test -f "$rom"
done

echo "All examples built successfully!"
