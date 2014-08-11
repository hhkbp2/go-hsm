package hsm

import "container/list"

// state IDs for all the default states
const (
    TopStateID      = "TOP"
    InitialStateID  = "Initial"
    TerminalStateID = "Terminal"
)

// State represents the interface that every state class should implements.
type State interface {
    // Returns the ID of this state
    ID() string

    // Returns parent state of this state
    Super() (super State)
    // Returns all children states of this state
    Children() []State
    // It add state child as a child of this state.
    AddChild(child State)

    // Called when this state is targeted in state initialization.
    Init(hsm HSM, event Event) (state State)
    // Called when entering this state
    Entry(hsm HSM, event Event) (state State)
    // Called when exiting this state
    Exit(hsm HSM, event Event) (state State)
    // Called when event dispatched to this state
    Handle(hsm HSM, event Event) (state State)
}

// StateHead is the head of every state to maintain the
// parent/child relationship so that all states compose
// the whole state hierarchy of state machine.
// It provides the default implementations of Super(), Children(), AddChild()
// for states.
type StateHead struct {
    // pointer to parent state
    super State
    // links to all children state
    children *list.List
}

// NewStateHead() is the constructor for StateHead.
func NewStateHead(super State) *StateHead {
    children := list.New()
    return &StateHead{
        super:    super,
        children: children,
    }
}

// Super() is part of interface State.
func (self *StateHead) Super() State {
    return self.super
}

// Children() is part of interface State.
func (self *StateHead) Children() []State {
    length := self.children.Len()
    children := make([]State, 0, length)
    for e := self.children.Front(); e != nil; e = e.Next() {
        state, ok := e.Value.(State)
        AssertTrue(ok)
        children = append(children, state)
    }
    return children
}

// AddChild() is part of interface State.
func (self *StateHead) AddChild(child State) {
    AssertFalse(ListIn(self.children, child))
    self.children.PushBack(child)
}

// Init() is part of interface State.
func (self *StateHead) Init(hsm HSM, event Event) (state State) {
    return self.Super()
}

// Entry() is part of interface State.
func (self *StateHead) Entry(hsm HSM, event Event) (state State) {
    return self.Super()
}

// Exit() is part of interface State.
func (self *StateHead) Exit(hsm HSM, event Event) (state State) {
    return self.Super()
}

// The default top state for state machines.
// It provides dummy implementations for interface State and presents the
// default hehaviors for every state.
type Top struct {
    *StateHead
}

func NewTop() *Top {
    return &Top{NewStateHead(nil)}
}

func (self *Top) ID() string {
    return TopStateID
}

func (self *Top) Init(hsm HSM, event Event) (state State) {
    return nil
}

func (self *Top) Entry(hsm HSM, event Event) (state State) {
    return nil
}

func (self *Top) Exit(hsm HSM, event Event) (state State) {
    return nil
}

func (self *Top) Handle(hsm HSM, event Event) (state State) {
    return nil
}

// The default initial state for state machines. It's the start point
// of state machine.
type Initial struct {
    *StateHead
    InitStateID string
}

func NewInitial(super State, initStateID string) *Initial {
    object := &Initial{NewStateHead(super), initStateID}
    super.AddChild(object)
    return object
}

func (*Initial) ID() string {
    return InitialStateID
}

func (self *Initial) Init(hsm HSM, event Event) (state State) {
    hsm.QInit(self.InitStateID)
    return nil
}

func (self *Initial) Handle(hsm HSM, event Event) (state State) {
    // Events are only dispatch after state initialization is done.
    // It should never run here.
    AssertTrue(false)
    return self.Super()
}

// the default terminal state for state machines. It's the end point
// of state machine.
type Terminal struct {
    *StateHead
}

func NewTerminal(super State) *Terminal {
    object := &Terminal{NewStateHead(super)}
    super.AddChild(object)
    return object
}

func (*Terminal) ID() string {
    return TerminalStateID
}

func (self *Terminal) Handle(hsm HSM, event Event) (state State) {
    // Events dispatched to terminal state are not handled by default.
    return self.Super()
}
