# Change One Thing

## What you are about to achieve

Change one visible detail in `examples/clearscreen/main.go` and rebuild `examples/clearscreen/game.z64` so you know the code you edit affects the ROM you run.

## Expected result

![The updated screen after changing one line](../images/beginner/change-one-thing.png)

## Minimal change

In `examples/clearscreen/main.go`, change:

```go
gosprite64.ClearScreenWith(gosprite64.Blue)
```

To:

```go
gosprite64.ClearScreenWith(gosprite64.Red)
```

From the repository root, rebuild the ROM:

```bash
./build_examples.sh
```

Then reopen `examples/clearscreen/game.z64`.

## What changed

You changed one rendering line in `examples/clearscreen/main.go` and rebuilt the ROM.

## Why it matters

Beginners gain confidence faster when they can make one tiny edit and immediately see the result.

## If this failed

Make sure you saved `examples/clearscreen/main.go`, ran `./build_examples.sh` from the repository root, and reopened `examples/clearscreen/game.z64`.

## Next step

Go to [Make Something Move](./03-make-something-move.md).
