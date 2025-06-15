#!/bin/bash
# Cross-platform build script for GoSprite64 examples
# Works on Windows (Git Bash), Linux, and macOS

# Exit immediately if a command exits with a non-zero status
set -e

# Function to print error message and exit
error_exit() {
    echo "which go: $(which go)" >&2
    echo "go version: $(go version)" >&2
    echo "GOTOOLCHAIN: $GOTOOLCHAIN" >&2
    echo "GOROOT: $GOROOT" >&2
    echo "GOBIN: $GOBIN" >&2
    echo "GOPATH: $GOPATH" >&2
    echo "GO111MODULE: $GO111MODULE" >&2
    echo "GOFLAGS: $GOFLAGS" >&2
    echo "GOOS: $GOOS" >&2
    echo "GOARCH: $GOARCH" >&2
    echo 'ls -l $(go env GOROOT)/bin >&2'
    ls -l $(go env GOROOT)/bin
    echo 'ls -l $(go env GOPATH)/bin >&2'
    ls -l $(go env GOPATH)/bin
    echo 'ls -l $(go env GOBIN) >&2'
    echo "ERROR: $1" >&2
    echo "cat .envrc"
    cat .envrc
    echo "direnv status"
    direnv status
    exit 1
}

# Function to build example in a directory
build_example() {
    local dir="$1"
    echo "Building example in $dir"
    
    # Change to the directory
    cd "$dir" || error_exit "Failed to change to directory $dir"
    
    # Run go build
    echo "  Running go build -o game.elf ."
    go build -o game.elf . || error_exit "Failed to build $dir"
    
    # Run mkrom
    echo "  Running mkrom game.elf"
    mkrom game.elf || error_exit "Failed to create ROM for $dir"
    
    # Check if files exist
    if [ ! -f "game.elf" ]; then
        error_exit "game.elf not found in $dir"
    fi
    
    if [ ! -f "game.z64" ]; then
        error_exit "game.z64 not found in $dir"
    fi
    
    echo "  Successfully built $dir"
    
    # Return to the original directory
    cd - >/dev/null
}

# Store the starting directory
start_dir=$(pwd)

# Navigate to the examples directory
cd "$(dirname "$0")/examples" || error_exit "Failed to navigate to examples directory"

# Find all directories in the examples folder
echo "Finding example directories..."
for dir in */; do
    # Remove trailing slash
    dir=${dir%/}
    
    # Skip if not a directory
    if [ ! -d "$dir" ]; then
        continue
    fi
    
    echo "Found example: $dir"
    build_example "$dir"
done

# Return to the starting directory
cd "$start_dir" || error_exit "Failed to return to starting directory"

echo "All examples built successfully!"
