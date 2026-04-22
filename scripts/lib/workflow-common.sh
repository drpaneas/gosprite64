#!/usr/bin/env bash

clean_go_env() {
  env -u GOENV -u GOOS -u GOARCH -u GOFLAGS -u GOTOOLCHAIN -u GOPATH -u GOBIN "$@"
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

ensure_no_stale_envrc() {
  local repo_root="$1"
  if [[ -f "$repo_root/.envrc" ]]; then
    echo "error: stale .envrc detected; remove it and use n64.env only" >&2
    return 1
  fi
}

ensure_no_nested_example_modules() {
  local repo_root="$1"
  local -a nested_example_modules=()
  mapfile -t nested_example_modules < <(
    find "$repo_root/examples" -mindepth 2 -maxdepth 2 \( -name go.mod -o -name go.sum \) -print | sort
  )

  if (( ${#nested_example_modules[@]} > 0 )); then
    echo "error: remove nested example go.mod/go.sum files before building" >&2
    printf '  %q\n' "${nested_example_modules[@]}" >&2
    return 1
  fi
}

ensure_host_go() {
  if ! command -v go >/dev/null 2>&1; then
    echo "error: host go command not found" >&2
    return 1
  fi
}

retry_macos_bootstrap() {
  if [[ "$(uname -s)" != "Darwin" ]]; then
    return 1
  fi

  echo "go1.24.5-embedded failed to start; retrying download with BOOT_GO_LDFLAGS=-w" >&2
  clean_go_env env BOOT_GO_LDFLAGS=-w go1.24.5-embedded download || return 1
  clean_go_env go1.24.5-embedded version >/dev/null 2>&1
}

ensure_embeddedgo_native() {
  if ! command -v go1.24.5-embedded >/dev/null 2>&1; then
    echo "error: go1.24.5-embedded not found" >&2
    echo "install it with:" >&2
    echo "  go install github.com/embeddedgo/dl/go1.24.5-embedded@latest" >&2
    echo "  go1.24.5-embedded download" >&2
    echo "on macOS, if download aborts with the __DATA / __DWARF dyld error, retry with:" >&2
    echo "  BOOT_GO_LDFLAGS=-w go1.24.5-embedded download" >&2
    fallback_instructions
    return 1
  fi

  if ! clean_go_env go1.24.5-embedded version >/dev/null 2>&1; then
    if ! retry_macos_bootstrap; then
      echo "error: go1.24.5-embedded is installed but failed to start on this host" >&2
      echo "if macOS failed during toolchain bootstrap, retry manually with:" >&2
      echo "  BOOT_GO_LDFLAGS=-w go1.24.5-embedded download" >&2
      fallback_instructions
      return 1
    fi
  fi
}

ensure_embeddedgo_linux() {
  ensure_host_go
  local gopath_bin
  gopath_bin="$(clean_go_env go env GOPATH)/bin"
  export PATH="$gopath_bin:$PATH"

  if ! command -v go1.24.5-embedded >/dev/null 2>&1; then
    clean_go_env go install github.com/embeddedgo/dl/go1.24.5-embedded@latest
    hash -r
  fi

  if ! clean_go_env go1.24.5-embedded version >/dev/null 2>&1; then
    clean_go_env go1.24.5-embedded download
  fi
}

ensure_n64go_native() {
  if ! command -v n64go >/dev/null 2>&1; then
    echo "error: n64go not found" >&2
    echo "install it with:" >&2
    echo "  go install github.com/clktmr/n64/tools/n64go@v0.1.2" >&2
    return 1
  fi

  if ! clean_go_env n64go -h >/dev/null 2>&1; then
    echo "error: n64go is installed but failed to start on this host" >&2
    return 1
  fi
}

ensure_n64go_linux() {
  if ! command -v n64go >/dev/null 2>&1; then
    clean_go_env go install github.com/clktmr/n64/tools/n64go@v0.1.2
    hash -r
  fi
}

ensure_n64_module_version() {
  local repo_root="$1"
  local required_version="$2"
  local current_version

  current_version="$(
    cd "$repo_root" &&
      clean_go_env go list -m -f '{{.Version}}' github.com/clktmr/n64
  )"

  if [[ "$current_version" != "$required_version" ]]; then
    echo "error: expected github.com/clktmr/n64 ${required_version}, got ${current_version}" >&2
    echo "error: run this script from a checkout that already contains the v0.1.2 dependency bump" >&2
    return 1
  fi
}

build_all_examples() {
  local repo_root="$1"

  cd "$repo_root" || return 1

  for example_dir in "$repo_root"/examples/*; do
    [[ -d "$example_dir" ]] || continue

    local example_name
    example_name="$(basename "$example_dir")"
    local elf="$example_dir/game.elf"
    local rom="$example_dir/game.z64"

    rm -f "$elf" "$rom"

    clean_go_env GOENV="$repo_root/n64.env" go1.24.5-embedded build -o "$elf" "./examples/$example_name"
    clean_go_env GOENV="$repo_root/n64.env" n64go rom "$elf"

    test -f "$elf"
    test -f "$rom"
  done

  echo "All examples built successfully!"
}
