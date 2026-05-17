package rspq

import "sync"

// Status represents RSP task execution status.
type Status int

const (
	StatusIdle Status = iota
	StatusRunning
	StatusYielded
	StatusDone
)

// Queue manages RSP task submission and completion.
// It handles the interleaving of graphics and audio tasks,
// where audio tasks can preempt graphics tasks via the yield mechanism.
type Queue struct {
	mu          sync.Mutex
	status      Status
	currentTask *Task
	pending     []*Task
	onComplete  func(*Task)
}

// NewQueue creates an RSP task queue.
func NewQueue() *Queue {
	return &Queue{
		status: StatusIdle,
	}
}

// SetCompletionHandler sets a callback invoked when a task completes.
func (q *Queue) SetCompletionHandler(fn func(*Task)) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.onComplete = fn
}

// Submit adds a task to the queue. If the RSP is idle, it starts immediately.
// If a graphics task is running and an audio task is submitted, the graphics
// task is yielded.
func (q *Queue) Submit(task *Task) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.status == StatusIdle {
		q.currentTask = task
		q.status = StatusRunning
		return
	}

	if task.Type == AudTask && q.currentTask != nil && q.currentTask.Type == GfxTask {
		q.status = StatusYielded
		q.pending = append([]*Task{task}, q.pending...)
		return
	}

	q.pending = append(q.pending, task)
}

// Complete marks the current task as done and starts the next pending task.
func (q *Queue) Complete() {
	q.mu.Lock()
	defer q.mu.Unlock()

	completed := q.currentTask
	q.currentTask = nil
	q.status = StatusIdle

	if len(q.pending) > 0 {
		q.currentTask = q.pending[0]
		q.pending = q.pending[1:]
		q.status = StatusRunning
	}

	if completed != nil && q.onComplete != nil {
		q.onComplete(completed)
	}
}

// Status returns the current RSP status.
func (q *Queue) Status() Status {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.status
}

// Current returns the currently executing task, or nil.
func (q *Queue) Current() *Task {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.currentTask
}
