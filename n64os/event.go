package n64os

import "sync"

// EventRouter routes hardware interrupt events to message queues,
// matching the osSetEventMesg pattern.
type EventRouter struct {
	mu       sync.Mutex
	handlers map[MessageType]*eventBinding
}

type eventBinding struct {
	queue *MessageQueue
	msg   Message
}

// NewEventRouter creates an event router.
func NewEventRouter() *EventRouter {
	return &EventRouter{
		handlers: make(map[MessageType]*eventBinding),
	}
}

// SetEventMsg registers a queue to receive messages when an event occurs.
// This matches osSetEventMesg(event, queue, msg).
func (er *EventRouter) SetEventMsg(event MessageType, queue *MessageQueue, msg Message) {
	er.mu.Lock()
	defer er.mu.Unlock()
	er.handlers[event] = &eventBinding{queue: queue, msg: msg}
}

// Signal fires an event, sending the registered message to its queue.
func (er *EventRouter) Signal(event MessageType) {
	er.mu.Lock()
	binding, ok := er.handlers[event]
	er.mu.Unlock()
	if ok {
		binding.queue.Send(binding.msg)
	}
}

// Clear removes the handler for an event.
func (er *EventRouter) Clear(event MessageType) {
	er.mu.Lock()
	defer er.mu.Unlock()
	delete(er.handlers, event)
}
