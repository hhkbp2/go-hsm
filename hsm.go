package hsm

import "container/List"
import "reflect"
import "fmt"

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
    type_ uint32
}

func (stdEvent *StdEvent) Type() uint32 {
    return stdEvent.type_
}

type State interface {
    ID() string

    GetSuper() *State
    SetSuper(super *State) void
    GetChildren() []*State
    AddChild(hsm *HSM, child *State) void
    GetState() *State
    setState(state *State) void

    Entry(hsm *HSM, event *Event) *State
    Exit(hsm *HSM, event *Event) *State
    Handle(hsm *HSM, event *Event) *State
}

type StateHead struct {
    super    *State
    children list.List
}

func NewStateHead(super *State) (*StateHead, error) {
    children = list.List.New()
    children.Init()
    head = &StateHead{
        super:    super,
        children: children,
    }
    return &head, nil
}

func (head *StateHead) AddChild(hsm *HSM, child *State) {
    head.children.PushBack(child)
    hsm.NodeMap[child.state.ID()] = child
}

type HSM struct {
    SourceState *State
    State       *State
    NodeMap     map[string]*State
}

func NewHSM(top *State) (*HSM, error) {
    hsm = &HSM{
        SourceState: nil,
        State:       top,
        NodeMap:     make(map[string]*State)}
    return hsm, nil
}

func (hsm *HSM) Init(event *Event) {

}

func (hsm *HSM) Dispatch(event *Event) {

    newState = hsm.State(event)
    HSM.State = newState
}

func TRIGGER_ENTRY(hsm *HSM, state *State, event *Event) *State {
    return node.Entry(hsm, event)
}

func TRIGGER_EXIT(hsm *HSM, state *State, event *Event) *State {
    return state.Exit(hsm, event)
}

func TRIGGER(hsm *HSM, state *State, event *Event) *State {
    return state.Handle(hsm, event)
}

func (hsm *HSM) TRAN(targetID string, event *Event) {
    target = hsm.NodeMap[targetID]
    for s := hsm.State; s != hsm.SourceState; {
        t := TRIGGER_EXIT(hsm, s, event)
        if t != nil {
            s = t
        } else {
            s = TRIGGER(hsm, s, &StdEvent{EventEmpty})
        }
    }
}

type S1 struct {
}

func NewS1() (*S1, error) {
    return &S1{}, nil
}

func (*S1) Name() string {
    return "S1"
}

func (*S1) Entry(event *Event) {
    fmt.Println("S1:Entry event=", event)
}

func (*S1) Exit(event *Event) {
    fmt.Println("S1:Exit event=", event)
}

type S11 struct {
    super S1
}

func NewS11() (*S11, error) {
    return &S11{}, nil
}

func (*S11) Name() string {
    return "S11"
}

func (*S11) Entry(event *Event) {
    fmt.Println("S11:Entry event=", event)
}

func (*S11) Exit(event *Event) {
    fmt.Println("S12:Exit event=", event)
}

type S12 struct {
    super S1
}

func NewS12() (*S12, error) {
    return &S12{}, nil
}

func (*S12) Name() string {
    return "S12"
}

func (*S12) Entry(event *Event) {
    fmt.Printnl("S12:Entry event=", event)
}

func (*S12) Exit(event *Event) {
    fmt.Println("S12:Exit event=", event)
}

type TopState struct {
}

func NewTopState() (*TopState, error) {
    return &TopState{}
}

func (*TopState) Name() string {
    return "TOP"
}

func (*TopState) Entry(event *Event) {
    TRAN("S1")
    return
}

func (*TopState) Exit(event *Event) {
    return
}

func main() {
    top, err = NewTopState()
    rootNode = NewState(nil, top)
    s1, err = NewS1()
    n1 = NewTreeNode(rootNode, s1)
    rootNode.AddChild(n1)
    s11, err = NewS11()
    n11 = NewTreeNode(n1, s11)
    n1.AddChild(n11)
    s12, err = NewS12()
    n12 = NewTreeNode(n1, s12)
    n1.AddChild(n12)

}
