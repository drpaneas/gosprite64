#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
bash_path="$(command -v bash)"

# shellcheck source=scripts/lib/workflow-common.sh
source "$repo_root/scripts/lib/workflow-common.sh"

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

mkdir -p "$tmpdir/examples/demo" "$tmpdir/bin" "$tmpdir/state"

touch "$tmpdir/.envrc"
if ensure_no_stale_envrc "$tmpdir" >"$tmpdir/envrc.out" 2>&1; then
  echo "expected ensure_no_stale_envrc to fail" >&2
  exit 1
fi
grep -q 'stale \.envrc detected' "$tmpdir/envrc.out"
rm -f "$tmpdir/.envrc"

printf 'module stale.example\n\ngo 1.24.3\n' >"$tmpdir/examples/demo/go.mod"
if ensure_no_nested_example_modules "$tmpdir" >"$tmpdir/nested.out" 2>&1; then
  echo "expected ensure_no_nested_example_modules to fail" >&2
  exit 1
fi
grep -q 'remove nested example go\.mod/go\.sum files before building' "$tmpdir/nested.out"
grep -q 'examples/demo/go\.mod' "$tmpdir/nested.out"
rm -f "$tmpdir/examples/demo/go.mod"

cat >"$tmpdir/bin/uname" <<'EOF'
#!/usr/bin/env bash
printf 'Darwin\n'
EOF
chmod +x "$tmpdir/bin/uname"

cat >"$tmpdir/bin/go1.24.5-embedded" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
log_file="${WORKFLOW_TEST_LOG:?}"
state_dir="${WORKFLOW_TEST_STATE:?}"
printf 'cmd=%s boot=%s\n' "$*" "${BOOT_GO_LDFLAGS-}" >> "$log_file"
case "${1:-}" in
  version)
    [[ -f "$state_dir/ok" ]] && exit 0
    exit 1
    ;;
  download)
    [[ "${BOOT_GO_LDFLAGS-}" == "-w" ]] || exit 1
    touch "$state_dir/ok"
    exit 0
    ;;
  *)
    exit 0
    ;;
esac
EOF
chmod +x "$tmpdir/bin/go1.24.5-embedded"

cat >"$tmpdir/bin/n64go" <<'EOF'
#!/usr/bin/env bash
[[ "${1:-}" == "-h" ]] && exit 0
exit 1
EOF
chmod +x "$tmpdir/bin/n64go"

PATH="$tmpdir/bin:$PATH" \
WORKFLOW_TEST_LOG="$tmpdir/retry.log" \
WORKFLOW_TEST_STATE="$tmpdir/state" \
ensure_embeddedgo_native

test "$(grep -c 'cmd=download boot=-w' "$tmpdir/retry.log")" -eq 1
grep -q 'cmd=version boot=' "$tmpdir/retry.log"

PATH="$tmpdir/bin:$PATH" ensure_n64go_native

# Verify the top-level wrapper scripts delegate to the shared workflow helpers.
wrapper_root="$tmpdir/wrapper-root"
mkdir -p "$wrapper_root/scripts/lib" "$wrapper_root/scripts" "$wrapper_root/examples/demo" "$wrapper_root/bin"

cp "$repo_root/build_examples.sh" "$wrapper_root/build_examples.sh"
cp "$repo_root/scripts/dev-linux-build.sh" "$wrapper_root/scripts/dev-linux-build.sh"

cat >"$wrapper_root/scripts/lib/workflow-common.sh" <<'EOF'
record_call() {
  local name="$1"
  shift

  printf '%s' "$name" >>"$WORKFLOW_WRAPPER_TEST_LOG"

  local arg
  for arg in "$@"; do
    printf '|%s' "$arg" >>"$WORKFLOW_WRAPPER_TEST_LOG"
  done

  printf '\n' >>"$WORKFLOW_WRAPPER_TEST_LOG"
}

ensure_no_stale_envrc() { record_call ensure_no_stale_envrc "$@"; }
ensure_no_nested_example_modules() { record_call ensure_no_nested_example_modules "$@"; }
ensure_host_go() { record_call ensure_host_go "$@"; }
ensure_embeddedgo_native() { record_call ensure_embeddedgo_native "$@"; }
ensure_n64go_native() { record_call ensure_n64go_native "$@"; }
ensure_embeddedgo_linux() { record_call ensure_embeddedgo_linux "$@"; }
ensure_n64go_linux() { record_call ensure_n64go_linux "$@"; }
ensure_n64_module_version() { record_call ensure_n64_module_version "$@"; }
build_all_examples() { record_call build_all_examples "$@"; }
EOF

cat >"$wrapper_root/bin/dirname" <<EOF
#!$bash_path
printf '%s\n' "\${1%/*}"
EOF
chmod +x "$wrapper_root/bin/dirname"

cat >"$wrapper_root/bin/uname" <<EOF
#!$bash_path
printf 'Linux\n'
EOF
chmod +x "$wrapper_root/bin/uname"

build_wrapper_log="$tmpdir/build-wrapper.log"
PATH="$wrapper_root/bin" WORKFLOW_WRAPPER_TEST_LOG="$build_wrapper_log" "$bash_path" "$wrapper_root/build_examples.sh"

mapfile -t build_wrapper_calls <"$build_wrapper_log"
test "${#build_wrapper_calls[@]}" -eq 6
test "${build_wrapper_calls[0]}" = "ensure_no_stale_envrc|$wrapper_root"
test "${build_wrapper_calls[1]}" = "ensure_no_nested_example_modules|$wrapper_root"
test "${build_wrapper_calls[2]}" = "ensure_host_go"
test "${build_wrapper_calls[3]}" = "ensure_embeddedgo_native"
test "${build_wrapper_calls[4]}" = "ensure_n64go_native"
test "${build_wrapper_calls[5]}" = "build_all_examples|$wrapper_root"

linux_wrapper_log="$tmpdir/linux-wrapper.log"
PATH="$wrapper_root/bin" WORKFLOW_WRAPPER_TEST_LOG="$linux_wrapper_log" "$bash_path" "$wrapper_root/scripts/dev-linux-build.sh"

mapfile -t linux_wrapper_calls <"$linux_wrapper_log"
test "${#linux_wrapper_calls[@]}" -eq 6
test "${linux_wrapper_calls[0]}" = "ensure_no_stale_envrc|$wrapper_root"
test "${linux_wrapper_calls[1]}" = "ensure_no_nested_example_modules|$wrapper_root"
test "${linux_wrapper_calls[2]}" = "ensure_embeddedgo_linux"
test "${linux_wrapper_calls[3]}" = "ensure_n64go_linux"
test "${linux_wrapper_calls[4]}" = "ensure_n64_module_version|$wrapper_root|v0.1.2"
test "${linux_wrapper_calls[5]}" = "build_all_examples|$wrapper_root"

echo "workflow guard tests passed"
