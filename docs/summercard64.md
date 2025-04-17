# 🔌 Loading Your Game on Real Hardware

You’ve built your `.z64` file using GoSprite64. You’ve tested it in the emulator.  
Now it’s time for the real deal—**running your game on an actual Nintendo 64 using SummerCart64** (aka SC64).

Whether you’re writing directly to the SD card or using the official `sc64deployer` tool via USB, we’ve got you covered.

---

## 🗂 Option 1: Copy to SD Card

The simplest way to test your game is to copy the ROM directly to the SD card used by the SummerCart64.

### 📥 Steps

1. Build your ROM:

   ```sh
   emgo build
   ```

    This will give you a .z64 file (e.g., clearscreen.z64).

2. Insert the SD card into your computer.

3. Copy the .z64 file onto it:

    ```sh
    cp clearscreen.z64 /Volumes/SC64/  # macOS
    cp clearscreen.z64 /run/media/yourname/SC64/  # Linux
    ```

4. Safely eject the SD card and put it back into your SC64.

5. Power on your N64. Your game should appear in the menu.

That’s it! You’re running Go code on a real N64.

## 🔌 Option 2: Upload via USB (sc64deployer)

If you’ve connected the SC64 via USB, you can upload ROMs directly using sc64deployer, the official SC64 control tool.

### 📦 Install sc64deployer

Clone the repo and build it:

```sh
./sc64deployer upload clearscreen.z64
```

Example output:

```sh
Uploading ROM [clearscreen.z64]... done
Save type set to [None]
Boot mode set to [Bootloader -> ROM]
```

### ⚠ Important: Power Cycle, Don’t Reset

After uploading:

> DO NOT press reset on your N64.

Instead: Power OFF the N64. Wait a moment. Power it back ON.

This ensures the ROM is correctly loaded and avoids leaving the SC64 in a weird state. Using the reset button can cause crashes or corrupted memory states.
