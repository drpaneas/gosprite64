# Getting Started

This repository supports one setup path: build natively with `go1.24.5-embedded` and `n64go@v0.1.2`. If that native bootstrap is unavailable on your host, use the Linux fallback below.

## Installation

You need a working `go` command on your host so you can install the embedded toolchain launcher and the ROM tool.

1. Clone the repository:

```bash
git clone https://github.com/drpaneas/gosprite64.git
cd gosprite64
```

2. Install the embedded Go toolchain launcher and download the toolchain:

```bash
go install github.com/embeddedgo/dl/go1.24.5-embedded@latest
go1.24.5-embedded download
```

3. Install `n64go`:

```bash
go install github.com/clktmr/n64/tools/n64go@v0.1.2
```

4. Build the examples:

```bash
./build_examples.sh
```

The repository tracks toolchain settings only in `go.env`:

```bash
GOTOOLCHAIN=go1.24.5-embedded
GOOS=noos
GOARCH=mips64
GOFLAGS='-tags=n64' '-trimpath' '-ldflags=-M=0x00000000:8M -F=0x00000400:8M -stripfn=1'
```

`./build_examples.sh` uses that file to build every example and generate ROMs under `examples/`, such as `examples/clearscreen/game.z64`.

## Linux Fallback

If `go1.24.5-embedded` cannot run natively on your host, run the verified Linux fallback from the repository root:

```bash
docker run --rm --platform linux/arm64 \
  -v "$PWD:/workspace/gosprite64" \
  -v gosprite64-gomod:/go/pkg/mod \
  -v gosprite64-gobuild:/root/.cache/go-build \
  -v gosprite64-sdk:/root/sdk \
  -w /workspace/gosprite64 \
  golang:1.26-bookworm \
  bash ./scripts/dev-linux-build.sh
```

This fallback produces the same `game.z64` files under `examples/`.
