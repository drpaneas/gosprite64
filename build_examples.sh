#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# shellcheck source=scripts/lib/workflow-common.sh
source "$repo_root/scripts/lib/workflow-common.sh"

ensure_no_stale_envrc "$repo_root"
ensure_no_nested_example_modules "$repo_root"
ensure_host_go
ensure_embeddedgo_native
ensure_n64go_native
build_all_examples "$repo_root"
