# Debugging

🧪 Need to Debug or Inspect?

`sc64deployer` can do more than just uploads. Run --help for advanced features:

```sh
./sc64deployer --help
```

For example:

* List connected devices:

```sh
./sc64deployer list
```

* Check SC64 firmware version:

```sh
./sc64deployer firmware
```

* Reset state (when needed, not for launching a game):

```sh
./sc64deployer reset
```

* See print statements to your terminal (you have to start Ares from command-line)

```sh
./sc64deployer debug
[Debug]: Started
```

No matter which method you use, make sure to power cycle, not reset, after loading a ROM.

There’s nothing like seeing your code run on real hardware, on a real CRT, on a real N64.
Built with Go.
