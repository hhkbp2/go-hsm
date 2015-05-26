package c_comment

import hsm "github.com/hhkbp2/go-hsm"

// Four states are defined according to the state chart.
// Define their IDs first since go-hsm use ID(of type string in Golang)
// as the target of state transfer.
const (
	StateCodeID    = "code"
	StateSlashID   = "slash"
	StateStarID    = "star"
	StateCommentID = "comment"
)

// CodeState represent the code state in state chart.
type CodeState struct {
	// the embedded StateHead implements most methods of the interface State(
	// without the method ID and Handle).
	// We use a feature in Golang called 'anonymous field' to achieve the
	// goal of providing interfaces and implementations heritance, just like
	// heritance from non-pure abstract class in C++ or
	// heritance from abstract parent class(not interface) in Java.
	*hsm.StateHead

	// entryCount and initCount are helper variables to demonstrate Init() of state
	// would always be triggered after its Entry().
	entryCount int
	initCount  int
}

func NewCodeState(super hsm.State) *CodeState {
	// initialized CodeState
	object := &CodeState{
		StateHead:  hsm.NewStateHead(super),
		entryCount: 0,
		initCount:  0,
	}
	// hook up to super
	super.AddChild(object)
	return object
}

// ID() is part of the interface State.
func (_ *CodeState) ID() string {
	return StateCodeID
}

// Entry() is part of the interface State.
// It would be called on state entry.
func (self *CodeState) Entry(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Entry")
	hsm.AssertEqual(event.Type(), hsm.EventEntry)
	hsm.AssertTrue(self.entryCount >= self.initCount)
	self.entryCount++
	return nil
}

// Init() is part of the interface State.
// It would be called on state initialization.
func (self *CodeState) Init(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Init")
	hsm.AssertEqual(event.Type(), hsm.EventInit)
	hsm.AssertTrue(self.entryCount >= self.initCount)
	self.initCount++
	// Return super state in Init() to tell state machine that
	// this state doesn't has state needed to be initialized.
	// If QInit() is called(there is a state/sub-state needed to be
	// initialized), then return nil.
	// QTran() is not suitable to call in Init(), only QInit() could
	// be used.
	return self.Super()
}

// Exit() is part of the interface State.
// It would be called on state exit.
func (self *CodeState) Exit(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Exit")
	hsm.AssertNotEqual(event.Type(), hsm.EventExit)
	hsm.AssertEqual(event.Type(), EventSlash)
	return nil
}

// Handle() is part of the interface State.
// It would be called with input events.
// Please refer to the state chart for all state transfers and
// actions round them.
func (self *CodeState) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle")
	switch event.Type() {
	case EventSlash:
		// Trigger a state transfer. The target state is slash state.
		sm.QTranOnEvent(StateSlashID, event)
		// Return nil to tell the state machine that this event is completely
		// handled. No need to propagrate it to our super state.
		return nil
	case EventChar:
		e, ok := event.(*CCommentCharEvent)
		hsm.AssertTrue(ok)
		commentHSM, ok := sm.(*CCommentHSM)
		hsm.AssertTrue(ok)
		commentHSM.ProcessCodeChar(e.Char())
		return nil
	}
	// Return super state to tell the state machine that this event is
	// not completely handled. Please propagrate it to our super state.
	return self.Super()
}

// SlashState represents the slash state in state chart.
type SlashState struct {
	*hsm.StateHead
}

func NewSlashState(super hsm.State) *SlashState {
	object := &SlashState{
		StateHead: hsm.NewStateHead(super),
	}
	super.AddChild(object)
	return object
}

func (self *SlashState) ID() string {
	return StateSlashID
}

func (self *SlashState) Entry(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Entry")
	return nil
}

func (self *SlashState) Exit(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Exit")
	return nil
}

func (self *SlashState) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle")
	switch event.Type() {
	case EventChar:
		fallthrough
	case EventSlash:
		e, ok := event.(CCommentEvent)
		hsm.AssertTrue(ok)
		commentHSM, ok := sm.(*CCommentHSM)
		hsm.AssertTrue(ok)
		// Record a slash char and e.Char() before calling QTran() to
		// trigger a state transfer, since both Entry and Exit of this state
		// are not good places for these codes.
		commentHSM.ProcessCodeChar('/')
		commentHSM.ProcessCodeChar(e.Char())
		sm.QTran(StateCodeID)
		return nil
	case EventStar:
		sm.QTran(StateCommentID)
		return nil
	}
	return self.Super()
}

// StarState represents the star state in state chart.
type StarState struct {
	*hsm.StateHead
}

func NewStarState(super hsm.State) *StarState {
	object := &StarState{
		StateHead: hsm.NewStateHead(super),
	}
	super.AddChild(object)
	return object
}

func (self *StarState) ID() string {
	return StateStarID
}

func (self *StarState) Entry(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Entry")
	return nil
}

func (self *StarState) Exit(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Exit")
	return nil
}

func (self *StarState) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle")
	switch event.Type() {
	case EventStar:
		sm.QTran(StateStarID)
		return nil
	case EventChar:
		sm.QTran(StateCommentID)
		return nil
	case EventSlash:
		sm.QTran(StateCodeID)
		return nil
	}
	return self.Super()
}

// CommentState represents the comment state in state chart.
type CommentState struct {
	*hsm.StateHead
}

func NewCommentState(super hsm.State) *CommentState {
	object := &CommentState{
		StateHead: hsm.NewStateHead(super),
	}
	super.AddChild(object)
	return object
}

func (self *CommentState) ID() string {
	return StateCommentID
}

func (self *CommentState) Entry(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Entry")
	return nil
}

func (self *CommentState) Exit(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Exit")
	return nil
}

func (self *CommentState) Handle(sm hsm.HSM, event hsm.Event) hsm.State {
	Logln(self.ID(), "- Handle")
	switch event.Type() {
	case EventChar:
		fallthrough
	case EventSlash:
		sm.QTran(StateCommentID)
		return nil
	case EventStar:
		sm.QTran(StateStarID)
		return nil
	}
	return self.Super()
}
