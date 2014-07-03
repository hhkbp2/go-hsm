package hsm

const (
    EventEmpty = iota
    EventInit
    EventEntry
    EventExit
    EventUser
)

type Event interface {
    Type() uint32
}

type StdEvent struct {
    EventType uint32
}

func NewStdEvent(eventType uint32) *StdEvent {
    return &StdEvent{eventType}
}

func (stdEvent *StdEvent) Type() uint32 {
    return stdEvent.EventType
}
