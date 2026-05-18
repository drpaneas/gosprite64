#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "Building docs with playable demos..."

# Build all examples first to produce ROMs
chmod +x "$repo_root/build_examples.sh"
"$repo_root/build_examples.sh"

# Create emulator roms directory
mkdir -p "$repo_root/docs/emulator/roms"

# Copy all built ROMs
for example_dir in "$repo_root"/examples/*/; do
    rom="$example_dir/game.z64"
    if [ -f "$rom" ]; then
        name="$(basename "$example_dir")"
        cp "$rom" "$repo_root/docs/emulator/roms/${name}.z64"
    fi
done

echo "ROMs copied: $(ls "$repo_root/docs/emulator/roms/" | wc -l) files"

# Build the mdbook
cd "$repo_root"
mdbook build

echo "Docs built successfully with playable demos!"
echo "Open book/index.html to view, or deploy book/ to your hosting."
