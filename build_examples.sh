#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

clean_go_env() {
  env -u GOENV -u GOOS -u GOARCH -u GOFLAGS -u GOTOOLCHAIN -u GOPATH -u GOBIN "$@"
}

retry_macos_bootstrap() {
  if [[ "$(uname -s)" != "Darwin" ]]; then
    return 1
  fi

  echo "go1.24.5-embedded failed to start; retrying download with BOOT_GO_LDFLAGS=-w" >&2
  clean_go_env env BOOT_GO_LDFLAGS=-w go1.24.5-embedded download || return 1
  clean_go_env go1.24.5-embedded version >/dev/null 2>&1
}

fallback_instructions() {
  cat >&2 <<'EOF'
Use the Linux fallback:
  docker run --rm --platform linux/arm64 \
    -v "$PWD:/workspace/gosprite64" \
    -v gosprite64-gomod:/go/pkg/mod \
    -v gosprite64-gobuild:/root/.cache/go-build \
    -v gosprite64-sdk:/root/sdk \
    -w /workspace/gosprite64 \
    golang:1.26-bookworm \
    bash ./scripts/dev-linux-build.sh
EOF
}

if [[ -f "$repo_root/.envrc" ]]; then
  echo "error: stale .envrc detected; remove it and use go.env only" >&2
  exit 1
fi

mapfile -t nested_example_modules < <(
  find "$repo_root/examples" -mindepth 2 -maxdepth 2 \( -name go.mod -o -name go.sum \) -print | sort
)

if (( ${#nested_example_modules[@]} > 0 )); then
  echo "error: remove nested example go.mod/go.sum files before building" >&2
  printf '  %s\n' "${nested_example_modules[@]}" >&2
  exit 1
fi

if ! command -v go >/dev/null 2>&1; then
  echo "error: host go command not found" >&2
  exit 1
fi

if ! command -v go1.24.5-embedded >/dev/null 2>&1; then
  echo "error: go1.24.5-embedded not found" >&2
  echo "install it with:" >&2
  echo "  go install github.com/embeddedgo/dl/go1.24.5-embedded@latest" >&2
  echo "  go1.24.5-embedded download" >&2
  echo "on macOS, if download aborts with the __DATA / __DWARF dyld error, retry with:" >&2
  echo "  BOOT_GO_LDFLAGS=-w go1.24.5-embedded download" >&2
  fallback_instructions
  exit 1
fi

if ! clean_go_env go1.24.5-embedded version >/dev/null 2>&1; then
  if ! retry_macos_bootstrap; then
    echo "error: go1.24.5-embedded is installed but failed to start on this host" >&2
    echo "if macOS failed during toolchain bootstrap, retry manually with:" >&2
    echo "  BOOT_GO_LDFLAGS=-w go1.24.5-embedded download" >&2
    fallback_instructions
    exit 1
  fi
fi

if ! command -v n64go >/dev/null 2>&1; then
  echo "error: n64go not found" >&2
  echo "install it with:" >&2
  echo "  go install github.com/clktmr/n64/tools/n64go@v0.1.2" >&2
  exit 1
fi

if ! clean_go_env n64go -h >/dev/null 2>&1; then
  echo "error: n64go is installed but failed to start on this host" >&2
  exit 1
fi

cd "$repo_root"

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
