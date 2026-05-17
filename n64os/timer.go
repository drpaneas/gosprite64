package n64os

import "time"

// Timer provides periodic or one-shot timer events,
// matching osSetTimer functionality.
type Timer struct {
	queue    *MessageQueue
	msg      Message
	interval time.Duration
	oneShot  bool
	stop     chan struct{}
	running  bool
}

// NewTimer creates a periodic timer that sends msg to queue at the given interval.
func NewTimer(queue *MessageQueue, msg Message, interval time.Duration) *Timer {
	return &Timer{
		queue:    queue,
		msg:      msg,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// NewOneShotTimer creates a timer that fires once after delay.
func NewOneShotTimer(queue *MessageQueue, msg Message, delay time.Duration) *Timer {
	return &Timer{
		queue:    queue,
		msg:      msg,
		interval: delay,
		oneShot:  true,
		stop:     make(chan struct{}),
	}
}

// Start begins the timer in a goroutine.
func (t *Timer) Start() {
	if t.running {
		return
	}
	t.running = true
	go func() {
		if t.oneShot {
			select {
			case <-time.After(t.interval):
				t.queue.Send(t.msg)
			case <-t.stop:
			}
		} else {
			ticker := time.NewTicker(t.interval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					t.queue.Send(t.msg)
				case <-t.stop:
					return
				}
			}
		}
		t.running = false
	}()
}

// Stop cancels the timer.
func (t *Timer) Stop() {
	if t.running {
		close(t.stop)
		t.running = false
	}
}
