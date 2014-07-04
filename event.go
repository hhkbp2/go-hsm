package hsm

type EventType uint32

const (
    EventEmpty = iota
    EventInit
    EventEntry
    EventExit
    EventUser
)

var StdEvents = map[EventType]*StdEvent{
    EventEmpty: NewStdEvent(EventEmpty),
    EventInit:  NewStdEvent(EventInit),
    EventEntry: NewStdEvent(EventEntry),
    EventExit:  NewStdEvent(EventExit),
}

type Event interface {
    Type() EventType
}

type StdEvent struct {
    EventType EventType
}

func NewStdEvent(eventType EventType) *StdEvent {
    return &StdEvent{eventType}
}

func (stdEvent *StdEvent) Type() EventType {
    return stdEvent.EventType
}
