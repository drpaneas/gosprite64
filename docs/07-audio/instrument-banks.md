# Instrument Banks

Define and load custom instrument samples for the sequence player.

```go
import "github.com/drpaneas/gosprite64/audio/bank"
```

## What is an instrument bank?

An instrument bank is a collection of instrument definitions paired with raw audio sample data. When the [sequence player](sequence-player.md) triggers a note event, it looks up the instrument assigned to that channel, selects the appropriate sample based on the note's pitch, and plays it through the audio hardware.

This is the same approach used by classic N64 games: a compact bank of instrument definitions stored alongside compressed sample data in the ROM. One bank might contain a piano, strings, drums, and bass - everything needed for a song.

## Types

### Instrument

Each instrument has metadata and one or more sounds that cover different key ranges:

```go
type Instrument struct {
    ID         uint8
    Volume     uint8    // default volume (0-127)
    Pan        uint8    // stereo panning (0=left, 64=center, 127=right)
    Priority   uint8    // voice allocation priority
    SampleRate uint16   // native sample rate in Hz
    KeyLow     uint8    // lowest MIDI note this instrument responds to
    KeyHigh    uint8    // highest MIDI note this instrument responds to
    Sounds     []Sound  // sample entries for different key ranges
}
```

### Sound

A sound is a single audio sample within an instrument. Instruments with multiple sounds use key-splitting to play different samples depending on the note pitch:

```go
type Sound struct {
    SampleAddr uint32   // offset into sample data
    SampleLen  uint32   // sample length in bytes
    LoopStart  uint32   // loop region start
    LoopEnd    uint32   // loop region end
    LoopCount  int32    // -1 for infinite, 0 for no loop
    KeyBase    uint8    // the note at which this sample plays at native pitch
    Tuning     float32  // fine-tuning adjustment
}
```

## Creating a bank

### Empty bank

```go
b := bank.NewBank()
```

Creates a bank with no instruments. You can populate it manually by appending to `b.Instruments` and setting `b.SampleData`.

### Loading from binary data

```go
b, err := bank.LoadBank(bankData)
if err != nil {
    // handle invalid bank data
}
```

`LoadBank` parses a binary bank file and returns a fully populated `Bank`. The binary format stores instrument definitions followed by per-instrument sound entries, all in big-endian byte order (matching N64 hardware).

The `bank.ErrInvalidBank` error is returned if the data is truncated or malformed.

## Querying a bank

### GetInstrument

Retrieve an instrument by its zero-based index:

```go
inst := b.GetInstrument(0)
if inst == nil {
    // index out of range
}
fmt.Printf("Instrument: rate=%d, sounds=%d\n", inst.SampleRate, len(inst.Sounds))
```

### InstrumentCount

Returns the total number of instruments in the bank:

```go
count := b.InstrumentCount()
for i := 0; i < count; i++ {
    inst := b.GetInstrument(uint8(i))
    fmt.Printf("[%d] keys %d-%d\n", inst.ID, inst.KeyLow, inst.KeyHigh)
}
```

## Complete example

```go
func NewGame() *Game {
    g := &Game{}

    // Load the instrument bank from ROM
    bankData := loadAsset("instruments.bnk")
    b, err := bank.LoadBank(bankData)
    if err != nil {
        panic("failed to load instrument bank: " + err.Error())
    }
    g.bank = b

    // Set up the sequence player
    g.music = sequence.NewPlayer()
    g.music.Data = loadAsset("overworld.seq")
    g.music.SetLoop(0, -1)
    g.music.Play()

    return g
}
```

## Bank binary format reference

The binary format parsed by `LoadBank`:

| Field | Type | Description |
|---|---|---|
| instrumentCount | u8 | Number of instruments |

Per instrument:

| Field | Type | Description |
|---|---|---|
| id | u8 | Instrument ID |
| volume | u8 | Default volume |
| pan | u8 | Stereo panning |
| priority | u8 | Voice priority |
| sampleRate | u16 BE | Native sample rate |
| keyLow | u8 | Lowest responding key |
| keyHigh | u8 | Highest responding key |
| soundCount | u8 | Number of sounds |

Per sound (28 bytes):

| Field | Type | Description |
|---|---|---|
| sampleAddr | u32 BE | Offset into sample data |
| sampleLen | u32 BE | Sample length in bytes |
| loopStart | u32 BE | Loop start offset |
| loopEnd | u32 BE | Loop end offset |
| loopCount | s32 BE | Loop count (-1 = infinite) |
| keyBase | u8 | Note at native pitch |
| (padding) | 3 bytes | Alignment padding |
| tuning | f32 BE | Fine-tuning factor |
