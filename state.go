package hsm

import "container/list"

const (
    TopStateID      = "TOP"
    InitialStateID  = "Initial"
    TerminalStateID = "Terminal"
)

type State interface {
    ID() string

    Super() (super State)
    Children() []State
    AddChild(child State)

    Init(hsm HSM, event Event) (state State)
    Entry(hsm HSM, event Event) (state State)
    Exit(hsm HSM, event Event) (state State)
    Handle(hsm HSM, event Event) (state State)
}

type StateHead struct {
    super    State
    children *list.List
}

func NewStateHead(super State) *StateHead {
    children := list.New()
    children.Init()
    return &StateHead{
        super:    super,
        children: children,
    }
}

func (self *StateHead) Super() State {
    return self.super
}

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

func (self *StateHead) AddChild(child State) {
    AssertFalse(ListIn(self.children, child))
    self.children.PushBack(child)
}

func (self *StateHead) Init(hsm HSM, event Event) (state State) {
    return self.Super()
}

func (self *StateHead) Entry(hsm HSM, event Event) (state State) {
    return self.Super()
}

func (self *StateHead) Exit(hsm HSM, event Event) (state State) {
    return self.Super()
}

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
    // should never be called
    return self.Super()
}

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
    // should never be called
    return self.Super()
}
