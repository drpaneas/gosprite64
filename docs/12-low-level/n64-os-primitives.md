# N64 OS Primitives

The `n64os` package provides Go implementations of the N64's core OS primitives: message queues, event routing, task scheduling, and timers. These map to the libultra/libdragon OS layer (OSMesgQueue, osSetEventMesg, osSetTimer) but use Go idioms like channels and goroutines.

## MessageQueue

A message queue is the fundamental communication primitive on the N64. Hardware interrupts, RSP completion signals, and game-logic events all flow through message queues.

```go
type MessageQueue struct {
    ch   chan Message
    size int
}
```

### Message types

```go
type Message struct {
    Type MessageType
    Data uint32
}

const (
    MsgNone    MessageType = iota
    MsgVBlank              // VI vertical retrace
    MsgSPDone              // RSP task complete
    MsgDPDone              // RDP rendering complete
    MsgSIDone              // Serial interface transfer complete
    MsgPIDone              // Parallel interface DMA complete
    MsgPreNMI              // Pre-reset signal
    MsgTimer               // Timer expiration
    MsgUser                // User-defined message
)
```

### Creating a queue

```go
q := n64os.NewMessageQueue(8) // capacity of 8 messages
```

Matches `osCreateMesgQueue` with the given buffer size.

### Sending messages

```go
ok := q.Send(n64os.Message{Type: n64os.MsgUser, Data: 42})
// Returns false if the queue is full (non-blocking)
```

`Jam` is a high-priority send. In the Go channel implementation, it behaves identically to `Send` since channels are FIFO. For true priority dispatch, use `select` on multiple queues.

### Receiving messages

Blocking receive:

```go
msg := q.Recv()
fmt.Println(msg.Type, msg.Data)
```

Non-blocking receive:

```go
msg, ok := q.TryRecv()
if ok {
    // process msg
}
```

### Utility methods

```go
q.Len()       // number of pending messages
ch := q.Channel() // underlying channel, for use in select
```

Using the channel directly in a `select`:

```go
select {
case msg := <-vblankQueue.Channel():
    // VBlank received
case msg := <-timerQueue.Channel():
    // Timer fired
}
```

## EventRouter

The `EventRouter` routes hardware interrupts to message queues, matching the `osSetEventMesg` pattern. When a hardware event fires, the router sends a pre-registered message to the associated queue.

```go
router := n64os.NewEventRouter()
```

### Registering events

```go
vblankQ := n64os.NewMessageQueue(1)
router.SetEventMsg(
    n64os.MsgVBlank,
    vblankQ,
    n64os.Message{Type: n64os.MsgVBlank},
)
```

This says: "When a VBlank interrupt occurs, send this message to this queue." Matches `osSetEventMesg(OS_EVENT_SP, queue, msg)` in libultra.

### Signaling events

When a hardware interrupt handler fires, it signals the router:

```go
router.Signal(n64os.MsgVBlank)
```

The router looks up the registered binding and sends the message to the queue. If no handler is registered for the event, the signal is silently dropped.

### Clearing events

```go
router.Clear(n64os.MsgVBlank)
```

Removes the handler for the given event type.

## Scheduler

The `Scheduler` manages RSP task execution, handling the priority relationship between audio and graphics tasks. Audio tasks preempt graphics tasks to meet real-time deadlines.

```go
sched := n64os.NewScheduler(eventRouter, func(task *n64os.Task) {
    // Called when a task starts executing on the RSP
})
```

### Task structure

```go
type Task struct {
    Type     TaskType    // TaskGraphics or TaskAudio
    Priority int
    Data     interface{} // task-specific payload
    Done     *MessageQueue // receives MsgSPDone when task completes
}
```

### Submitting tasks

```go
gfxDone := n64os.NewMessageQueue(1)
sched.SubmitGraphics(&n64os.Task{
    Data: displayList,
    Done: gfxDone,
})

audDone := n64os.NewMessageQueue(1)
sched.SubmitAudio(&n64os.Task{
    Data: audioCommandList,
    Done: audDone,
})
```

### Running the scheduler

The scheduler runs in its own goroutine:

```go
go sched.Run()
```

`Run` reads from the submit channel and dispatches tasks with this priority order:

1. **Audio tasks** always run first. If a graphics task is currently running when an audio task arrives, the graphics task is yielded.
2. **Yielded graphics tasks** resume after all audio tasks complete.
3. **Queued graphics tasks** run when nothing else is pending.

### Task completion

When the RSP finishes a task, call `CompleteTask`:

```go
sched.CompleteTask()
```

This sends `MsgSPDone` to the completed task's `Done` queue and dispatches the next pending task.

## Timer

The `Timer` provides periodic or one-shot timer events, matching `osSetTimer`.

### Periodic timer

```go
timerQ := n64os.NewMessageQueue(4)
t := n64os.NewTimer(
    timerQ,
    n64os.Message{Type: n64os.MsgTimer, Data: 1},
    time.Second / 60, // fire every frame at 60 FPS
)
t.Start()
```

The timer sends the configured message to the queue at every interval. It runs in a goroutine using `time.Ticker`.

### One-shot timer

```go
t := n64os.NewOneShotTimer(
    timerQ,
    n64os.Message{Type: n64os.MsgTimer, Data: 99},
    3 * time.Second, // fire once after 3 seconds
)
t.Start()
```

### Stopping a timer

```go
t.Stop()
```

Cancels the timer. A stopped timer cannot be restarted (the stop channel is closed).

## Complete example: game loop

A typical N64 game loop uses message queues to synchronize with VBlank and task completion:

```go
// Set up event routing
router := n64os.NewEventRouter()
vblankQ := n64os.NewMessageQueue(1)
router.SetEventMsg(n64os.MsgVBlank, vblankQ, n64os.Message{Type: n64os.MsgVBlank})

spDoneQ := n64os.NewMessageQueue(1)
router.SetEventMsg(n64os.MsgSPDone, spDoneQ, n64os.Message{Type: n64os.MsgSPDone})

// Set up scheduler
sched := n64os.NewScheduler(router, func(task *n64os.Task) {
    // Load microcode and start RSP
})
go sched.Run()

// Game loop
for {
    // Wait for VBlank
    vblankQ.Recv()

    // Update game logic
    update()

    // Build display list
    dl := buildDisplayList()

    // Submit graphics task
    taskDone := n64os.NewMessageQueue(1)
    sched.SubmitGraphics(&n64os.Task{
        Data: dl,
        Done: taskDone,
    })

    // Wait for RSP to finish
    taskDone.Recv()

    // Swap framebuffers
    swapBuffers()
}
```
