package n64os

type TaskType int

const (
	TaskGraphics TaskType = iota
	TaskAudio
)

type Task struct {
	Type     TaskType
	Priority int
	Data     interface{}
	Done     *MessageQueue
}

type Scheduler struct {
	events    *EventRouter
	submit    chan *Task
	gfxQueue  []*Task
	audQueue  []*Task
	running   *Task
	yielded   *Task
	onExecute func(*Task)
}

func NewScheduler(events *EventRouter, onExecute func(*Task)) *Scheduler {
	return &Scheduler{
		events:    events,
		submit:    make(chan *Task, 16),
		onExecute: onExecute,
	}
}

func (s *Scheduler) SubmitGraphics(task *Task) {
	task.Type = TaskGraphics
	s.submit <- task
}

func (s *Scheduler) SubmitAudio(task *Task) {
	task.Type = TaskAudio
	s.submit <- task
}

func (s *Scheduler) Run() {
	for task := range s.submit {
		switch task.Type {
		case TaskAudio:
			if s.running != nil && s.running.Type == TaskGraphics {
				s.yielded = s.running
				s.running = nil
			}
			s.audQueue = append(s.audQueue, task)
		case TaskGraphics:
			s.gfxQueue = append(s.gfxQueue, task)
		}
		s.dispatch()
	}
}

func (s *Scheduler) CompleteTask() {
	if s.running != nil {
		if s.running.Done != nil {
			s.running.Done.Send(Message{Type: MsgSPDone})
		}
		s.running = nil
	}
	s.dispatch()
}

func (s *Scheduler) dispatch() {
	if s.running != nil {
		return
	}
	if len(s.audQueue) > 0 {
		s.running = s.audQueue[0]
		s.audQueue = s.audQueue[1:]
		if s.onExecute != nil {
			s.onExecute(s.running)
		}
		return
	}
	if s.yielded != nil {
		s.running = s.yielded
		s.yielded = nil
		if s.onExecute != nil {
			s.onExecute(s.running)
		}
		return
	}
	if len(s.gfxQueue) > 0 {
		s.running = s.gfxQueue[0]
		s.gfxQueue = s.gfxQueue[1:]
		if s.onExecute != nil {
			s.onExecute(s.running)
		}
		return
	}
}
