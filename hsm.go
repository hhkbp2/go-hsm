package hsm

import "container/List"
import "reflect"
import "fmt"
import "assert"

const (
    EventEmpty = iota
    EventInit
    EventEntry
    EventExit
    EventUser
)

const TopStateID = "TOP"

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

    Super() *State
    SetSuper(super *State) void
    Children() []*State
    AddChild(hsm *HSM, child *State) void
    State() *State
    setState(state *State) void

    Init(hsm *HSM, event *Event) *State
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

// initial must set top as parent
func NewHSM(top, initial *State) (*HSM, error) {
    hsm = &HSM{
        SourceState: initial,
        State:       top,
        NodeMap:     make(map[string]*State)}
    hsm.NodeMap[top.ID()] = top
    return hsm, nil
}

func (hsm *HSM) Init() {
    hsm.init(&StdEvent{EventEmpty})
}

func (hsm *HSM) init(event *Event) {
    // health check
    assert.NotEqual(nil, hsm.State)
    assert.NotEqual(nil, hsm.SourceState)
    // check top state initialized. hsm.State.ID() should be "TOP"
    assert.Equal(hsm.NodeMap[hsm.State.ID()], hsm.State) // HSM not executed yet
    // save State in a temporary
    s := hsm.State
    // top-most initial transition
    Trigger(hsm, hsm.SourceState, event)
    // initial transition must go *one* level deep
    assert.Equal(s, Trigger(hsm, hsm.State, &StdEvent{EventEmpty}))
    // update the termporary
    s = hsm.State
    // enter the state
    Trigger(hsm, s, &StdEvent{EventEntry})
    for Trigger(hsm, s, &StdEvent{EventInit}) == nil { // init handled?
        // initial transition must go *one* level deep
        assert.Equal(s, Trigger(hsm, hsm.State, &StdEvent{EventEmpty}))
        s = hsm.State
        // enter the substate
        Trigger(hsm, s, &StdEvent{EventEntry})
    }
    // we are in well-initialized state now
}

func (hsm *HSM) IsInState(state *State) bool {
    // nagivate from current state up to all super state and
    // try to find specified `state'
    for s := hsm.State; s != nil; s = Trigger(hsm.State, &StdEvent{EventEmpty}) {
        if s == state {
            // a match is found
            return true
        }
    }
    // no match found
    return false
}

func (hsm *HSM) Dispatch(event *Event) {
    // Use `SourceState' to record the state which handle the event indeed(which
    // could be super, super-super, ... state).
    // `State' would stay unchange pointing at the current(most concrete) state.
    for hsm.SourceState = hsm.State; hsm.SourceState != nil; {
        hsm.SourceState = Trigger(hsm, hsm.SourceState, event)
    }
}

func Trigger(hsm *HSM, state *State, event *Event) *State {
    // dispatch the specified `event' to the corresponding method
    switch event.Type() {
    case EventEmpty:
        return state.Super()
    case EventInit:
        return state.Init(hsm, event)
    case EventEntry:
        return state.Entry(hsm, event)
    case EventExit:
        return state.Exit(hsm, event)
    default:
        return state.Handle(hsm, event)
    }
}

func (hsm *HSM) Tran(targetID string, event *Event) {
    assert.NotEqual(TopStateID, targetID)
    target = hsm.NodeMap[targetID]
    for s := hsm.State; s != hsm.SourceState; {
        // we are about to dereference `s'
        assert.NotEqual(nil, s)
        t := Trigger(hsm, s, &StdEvent{EventExit})
        if t != nil {
            s = t
        } else {
            s = Trigger(hsm, s, &StdEvent{EventEmpty})
        }
    }

    stateChain := list.List.New()
    stateChain.Init()

    stateChain.PushBack(nil)
    stateChain.PushBack(target)

    // (a) check `SourceState' == `target' (transition to self)
    if hsm.SourceState == target {
        Trigger(hsm, hsm.SourceState, &StdEvent{EventExit})
        goto inLCA
    }
    // (b) check `SourceState' == `target.Super()'
    p := Trigger(hsm, target, &StdEvent{EventEmpty})
    if hsm.SourceState == p {
        goto inLCA
    }
    // (c) check `SourceState.Super()' == `target.Super()' (most common)
    q := Trigger(hsm, hsm.SourceState, &StdEvent{EventEmpty})
    if q == p {
        Trigger(hsm, hsm.SourceState, &StdEvent{EventExit})
        goto inLCA
    }
    // (d) check `SourceState.Super()' == `target'
    if q == target {
        Trigger(hsm, hsm.SourceState, &StdEvent{EventExit})
        stateChain.Remove(stateChain.Back())
        goto inLCA
    }
    // (e) check rest of `SourceState' == `target.Super().Super()...'  hierarchy
    stateChain.PushBack(p)
    s = Trigger(hsm, p, &StdEvent{EventEmpty})
    for s != nil {
        if hsm.SourceState == s {
            goto inLCA
        }
        stateChain.PushBack(s)
        s = Trigger(hsm, s, &StdEvent{EventEmpty})
    }
    // exit source state
    Trigger(hsm, hsm.SourceState, &StdEvent{EventExit})
    // (f) check rest of `SourceState.Super()' == `target.Super().Super()...'
    for lca := stateChain.Back(); lca != nil; lca = lca.Prev() {
        if q == lca {
            // do not enter the LCA
            stateChain.Remove(stateChain.Back())
            goto inLCA
        }
    }
    // (g) check each `SourceState.Super().Super()...' for target...
    for s = q; s != nil; s = Trigger(s, &StdEvent{EventEmpty}) {
        for lca := stateChain.Back(); lca != nil; lca = lca.Prev() {
            if s == lca {
                stateChain = TruncateList(stateChain, lca)
                goto inLCA
            }
        }
        Trigger(hsm, s, &StdEvent{EventExit})
    }
    // malformed HSM
    assert.True(false)
inLCA: // now we are in the LCA of `SourceState' and `target'
    // retrace the entry path in reverse order
    for s = stateChain.Back(); s != nil && s.Value != nil; {
        Trigger(hsm, s, &StdEvent{EventEntry}) // enter `s' state
        stateChain.Remove(stateChain.Back())
        s = stateChain.Back()
    }
    // update current state
    hsm.State = target
    for Trigger(hsm, target, &StdEvent{EventInit}) == nil {
        // initial transition must go *one* level deep
        assertEqual(s, Trigger(hsm, hsm.State, &StdEvent{EventEmpty}))
        target = hsm.State
        Trigger(hsm, target, &StdEvent{EventEntry})
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
    hsm.Tran("S1")
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
