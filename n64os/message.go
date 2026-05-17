package n64os

// Message represents an OS message, similar to OSMesg.
type Message struct {
	Type MessageType
	Data uint32
}

// MessageType identifies the kind of message.
type MessageType int

const (
	MsgNone      MessageType = iota
	MsgVBlank                // VI vertical retrace
	MsgSPDone                // RSP task complete
	MsgDPDone                // RDP rendering complete
	MsgSIDone                // Serial interface transfer complete
	MsgPIDone                // Parallel interface DMA complete
	MsgPreNMI                // Pre-reset signal
	MsgTimer                 // Timer expiration
	MsgUser                  // User-defined message
)

// MessageQueue is a Go channel-based message queue matching OSMesgQueue semantics.
type MessageQueue struct {
	ch   chan Message
	size int
}

// NewMessageQueue creates a message queue with the given capacity.
func NewMessageQueue(size int) *MessageQueue {
	return &MessageQueue{
		ch:   make(chan Message, size),
		size: size,
	}
}

// Send adds a message to the queue. Returns false if the queue is full (non-blocking).
func (q *MessageQueue) Send(msg Message) bool {
	select {
	case q.ch <- msg:
		return true
	default:
		return false
	}
}

// Jam pushes a message to the front of the queue (high priority).
// In Go channels this inserts normally since channels are FIFO;
// for true priority, use Recv with select on multiple queues.
func (q *MessageQueue) Jam(msg Message) bool {
	return q.Send(msg)
}

// Recv blocks until a message is available, then returns it.
func (q *MessageQueue) Recv() Message {
	return <-q.ch
}

// TryRecv returns a message if one is available, or false.
func (q *MessageQueue) TryRecv() (Message, bool) {
	select {
	case msg := <-q.ch:
		return msg, true
	default:
		return Message{}, false
	}
}

// Len returns the number of pending messages.
func (q *MessageQueue) Len() int {
	return len(q.ch)
}

// Channel returns the underlying channel for use in select statements.
func (q *MessageQueue) Channel() <-chan Message {
	return q.ch
}
