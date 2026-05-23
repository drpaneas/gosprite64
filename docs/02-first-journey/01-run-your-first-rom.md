# Run Your First ROM

## What you are about to achieve

Build `examples/clearscreen/game.z64` and see a solid blue screen in your emulator.

## Expected result

![A solid blue screen rendered by the first ROM](../images/beginner/first-rom-blue-screen.png)

## Minimal commands

```bash
chmod +x ./build_examples.sh
./build_examples.sh
```

Open `examples/clearscreen/game.z64` after the build.

## What changed

You did not write game code yet. You proved that your toolchain can build the repository examples into ROMs, starting from `examples/clearscreen/main.go`.

## Why it matters

This is the fastest confirmation that your machine can produce working N64 output before you start editing code.

## If this failed

Go to [Installation](../02-getting-started/installation.md) and work through the toolchain setup notes first.

## Next step

Go to [Change One Thing](./02-change-one-thing.md).
