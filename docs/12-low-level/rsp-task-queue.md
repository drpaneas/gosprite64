# RSP Task Queue

The `rspq` package manages submission and execution of tasks on the N64's RSP (Reality Signal Processor). The RSP is a MIPS-based coprocessor that runs microcode for graphics (display list processing) and audio (sample mixing). This package provides task queuing, microcode management, and the submit/wait protocol.

## Microcode

The RSP has 4KB of instruction memory (IMEM) and 4KB of data memory (DMEM). Microcode is loaded into these memories before task execution.

### MicrocodeType

Identifies which RSP microcode to use:

```go
const (
    Fast3D  MicrocodeType = iota // original SM64 microcode
    F3DEX                        // extended vertex buffer (32 slots)
    F3DEX2                       // most common N64 microcode
    AspMain                      // standard audio microcode
)
```

### Microcode struct

Holds the IMEM and DMEM blobs:

```go
type Microcode struct {
    Type MicrocodeType
    Code []byte  // IMEM content (up to 4KB)
    Data []byte  // DMEM content (up to 4KB)
}
```

Maximum sizes:

```go
const MaxIMEMSize = 4096
const MaxDMEMSize = 4096
```

### Load

Loads microcode into the RSP:

```go
rspq.Load(ucode)
```

On N64, this wraps the `rsp.Load` function from the hardware driver. On other platforms, it is a no-op.

### Start

Resumes RSP execution after loading microcode:

```go
rspq.Start()
```

## Task

A `Task` describes an RSP task to execute, matching the N64's `OSTask` structure:

```go
type Task struct {
    Type  TaskType
    Flags TaskFlags

    // Microcode pointers
    UcodeBootAddr uint32
    UcodeBootSize uint32
    UcodeAddr     uint32
    UcodeSize     uint32
    UcodeDataAddr uint32
    UcodeDataSize uint32

    // Task data (display list or audio command list)
    DataAddr uint32
    DataSize uint32

    // Yield buffer for task preemption
    YieldAddr uint32
    YieldSize uint32

    // Output buffer (for audio tasks)
    OutputAddr uint32
    OutputSize uint32

    // DRAM stack for RSP microcode
    DRAMStackAddr uint32
    DRAMStackSize uint32
}
```

### TaskType

```go
const (
    GfxTask TaskType = 1  // graphics (display list) task
    AudTask TaskType = 2  // audio synthesis task
)
```

## SubmitTask

Loads a task onto the RSP and starts execution. This replicates the `osSpTaskLoad` + `osSpTaskStartGo` protocol:

```go
func SubmitTask(task *OSTask, bootCode []byte)
```

The submission sequence:

1. Marshal the task struct to a 64-byte buffer
2. Writeback the buffer to ensure cache coherency
3. Clear RSP signal flags and enable interrupt-on-break
4. DMA the task descriptor to DMEM at offset 0xFC0
5. DMA the boot microcode to IMEM
6. Resume RSP execution

All address fields in the task must be physical RDRAM addresses.

## WaitTaskDone

Blocks until the RSP signals completion (break interrupt):

```go
func WaitTaskDone()
```

The RSP sets the break flag when it finishes executing a task. This function waits for that interrupt.

## Queue

The `Queue` manages task submission and interleaving. It handles the priority relationship between audio and graphics tasks.

```go
q := rspq.NewQueue()
```

### Task priority and preemption

Audio tasks have higher priority than graphics tasks. If an audio task is submitted while a graphics task is running, the graphics task is yielded:

```go
q.Submit(gfxTask)  // starts running (RSP was idle)
q.Submit(audTask)  // graphics task is yielded, audio runs next
```

The yielded graphics task resumes after the audio task completes.

### Queue status

```go
status := q.Status()
// StatusIdle     - RSP is idle
// StatusRunning  - a task is executing
// StatusYielded  - a graphics task was preempted
// StatusDone     - last task completed
```

### Submit

Adds a task to the queue:

```go
q.Submit(task)
```

If the RSP is idle, the task starts immediately. If a graphics task is running and an audio task is submitted, the current task is yielded and the audio task runs first. Otherwise, the task is queued.

### Complete

Marks the current task as done and dispatches the next pending task:

```go
q.Complete()
```

### Completion handler

Register a callback for task completion:

```go
q.SetCompletionHandler(func(task *rspq.Task) {
    if task.Type == rspq.GfxTask {
        // frame rendering complete, swap buffers
    }
})
```

### Current

Returns the currently executing task:

```go
task := q.Current()
if task != nil {
    fmt.Println(task.Type)
}
```

## Typical graphics frame

```go
// Build display list
dl := gfx.NewDisplayList(256)
// ... add commands ...
dl.SPEndDisplayList()

// Create and submit task
task := &rspq.Task{
    Type:          rspq.GfxTask,
    UcodeAddr:     ucodeAddr,
    UcodeSize:     ucodeSize,
    UcodeDataAddr: ucodeDataAddr,
    UcodeDataSize: ucodeDataSize,
    DataAddr:      displayListAddr,
    DataSize:      uint32(dl.Len() * 8),
}

rspq.SubmitTask(task, bootCode)
rspq.WaitTaskDone()
```

## Audio/graphics interleaving

The N64 runs audio and graphics on the same RSP. Audio tasks must complete within tight deadlines to avoid audio glitches, so they preempt graphics tasks:

```
Frame timeline:
  [GFX running]──>[yield]──>[AUD runs]──>[AUD done]──>[GFX resumes]──>[GFX done]
```

The `Queue` handles this automatically. When you submit an audio task while graphics is running, the queue sets `StatusYielded`, puts the audio task at the front of the pending list, and dispatches it when the current task checkpoint allows.
