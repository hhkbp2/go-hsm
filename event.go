package hsm

type EventType uint32

// The types of predefined events.
const (
	EventEmpty EventType = iota
	EventInit
	EventEntry
	EventExit
	EventUser
)

// The default events(used in state transfer procedure).
// They are defined as global const for optimization.
var StdEvents = map[EventType]*StdEvent{
	EventEmpty: NewStdEvent(EventEmpty),
	EventInit:  NewStdEvent(EventInit),
	EventEntry: NewStdEvent(EventEntry),
	EventExit:  NewStdEvent(EventExit),
}

// Event represents the interface that every message which is dispatched
// to state machine should implements.
type Event interface {
	// Returns the type of this event
	Type() EventType
}

// StdEvent is the default Event implementation.
type StdEvent struct {
	EventType EventType
}

func NewStdEvent(eventType EventType) *StdEvent {
	return &StdEvent{eventType}
}

func (stdEvent *StdEvent) Type() EventType {
	return stdEvent.EventType
}
