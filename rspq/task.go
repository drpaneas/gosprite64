package rspq

// TaskType identifies the RSP task category.
type TaskType uint32

const (
	// GfxTask is a graphics (display list) task.
	GfxTask TaskType = 1
	// AudTask is an audio synthesis task.
	AudTask TaskType = 2
)

// TaskFlags control task execution behavior.
type TaskFlags uint32

// Task describes an RSP task to execute, matching the OSTask structure.
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
