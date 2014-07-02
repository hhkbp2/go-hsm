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
    Type_ uint32
}

func NewStdEvent(type_ uint32) *StdEvent {
    return &StdEvent{type_}
}

func (stdEvent *StdEvent) Type() uint32 {
    return stdEvent.Type_
}
