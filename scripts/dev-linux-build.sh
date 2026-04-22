#!/usr/bin/env bash
set -euo pipefail

if [[ "$(uname -s)" != "Linux" ]]; then
  echo "error: scripts/dev-linux-build.sh must run inside a Linux container or Linux CI runner" >&2
  exit 1
fi

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# shellcheck source=scripts/lib/workflow-common.sh
source "$repo_root/scripts/lib/workflow-common.sh"

ensure_no_stale_envrc "$repo_root"
ensure_no_nested_example_modules "$repo_root"
ensure_embeddedgo_linux
ensure_n64go_linux
ensure_n64_module_version "$repo_root" "v0.1.2"
build_all_examples "$repo_root"
