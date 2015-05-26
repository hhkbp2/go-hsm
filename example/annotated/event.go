package annotated

import hsm "github.com/hhkbp2/go-hsm"

const (
	EventA hsm.EventType = hsm.EventUser + 1 + iota
	EventB
	EventC
	EventD
	EventE
	EventF
	EventG
	EventH
)

var EventsToStr = map[hsm.EventType]string{
	EventA: "A",
	EventB: "B",
	EventC: "C",
	EventD: "D",
	EventE: "E",
	EventF: "F",
	EventG: "G",
	EventH: "H",
}

func PrintEvent(eventType hsm.EventType) string {
	return EventsToStr[eventType]
}

type AnnotatedEvent interface {
	hsm.Event
}

type GeneralEvent struct {
	*hsm.StdEvent
}

func NewEvent(eventType hsm.EventType) *GeneralEvent {
	return &GeneralEvent{
		hsm.NewStdEvent(eventType),
	}
}
