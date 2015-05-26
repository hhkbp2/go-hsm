package c_comment

import hsm "github.com/hhkbp2/go-hsm"

// According to the state chart, there are three types of
// different events representing the input of our state machine.
const (
	// EventSlash represents a slash character as input
	EventSlash hsm.EventType = hsm.EventUser + iota
	// EventStar represents a star character is meet
	EventStar
	// EventChar represents a character which is neither a slash nor a star
	EventChar
)

// CCommentEvent is the general event we use in this example to drive
// out state machine.
type CCommentEvent interface {
	hsm.Event
	Char() byte
}

// CCommentSlashEvent wraps everything for a event of EventSlash type.
// For the simplicity of this example, it is left empty.
// An event struct may contain any datas and methods needed in real project.
type CCommentSlashEvent struct {
	*hsm.StdEvent
	// empty struct
}

func NewCCommentSlashEvent() *CCommentSlashEvent {
	return &CCommentSlashEvent{
		hsm.NewStdEvent(EventSlash),
	}
}

// Char() is part of CCommentEvent interface.
func (_ *CCommentSlashEvent) Char() byte {
	return '/'
}

// CCommentStarEvent wraps everything for a event of EventStar type.
// For the simplicity of this example, it is left empty.
type CCommentStarEvent struct {
	*hsm.StdEvent
	// empty struct
}

func NewCCommentStarEvent() *CCommentStarEvent {
	return &CCommentStarEvent{
		hsm.NewStdEvent(EventStar),
	}
}

// Char() is part of CCommentEvent interface.
func (_ *CCommentStarEvent) Char() byte {
	return '*'
}

// CCommentCharEvent wraps everything for a event of EventChar type.
// For the simplicity of this example, it contains only a single character.
type CCommentCharEvent struct {
	*hsm.StdEvent
	c byte
}

func NewCCommentCharEvent(c byte) *CCommentCharEvent {
	return &CCommentCharEvent{
		hsm.NewStdEvent(EventChar),
		c,
	}
}

// Char() is part of CCommentEvent interface.
func (self *CCommentCharEvent) Char() byte {
	return self.c
}
