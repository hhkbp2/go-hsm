package hsm

import "container/list"
import "fmt"
import "github.com/hhkbp2/go-hsm/assert"

const (
    EventEmpty = iota
    EventInit
    EventEntry
    EventExit
    EventUser
)

const (
    Event1 = EventUser + iota
    Event2
    Event3
)

const TopStateID = "TOP"

type Event interface {
    Type() uint32
}

type StdEvent struct {
    Type_ uint32
}

func NewStdEvent(type_ uint32) (*StdEvent, error) {
    return &StdEvent{type_}, nil
}

func (stdEvent *StdEvent) Type() uint32 {
    return stdEvent.Type_
}

type State interface {
    ID() string

    Super() (super State)

    Init(hsm *HSM, event Event) (state State)
    Entry(hsm *HSM, event Event) (state State)
    Exit(hsm *HSM, event Event) (state State)
    Handle(hsm *HSM, event Event) (state State)
}

type StateHead struct {
    super State
}

func MakeStateHead(super State) StateHead {
    return StateHead{
        super: super,
    }
}

func (head *StateHead) Super() State {
    return head.super
}

func (head *StateHead) Init(hsm *HSM, event Event) (state State) {
    return head.Super()
}

func (head *StateHead) Entry(hsm *HSM, event Event) (state State) {
    return head.Super()
}

func (head *StateHead) Exit(hsm *HSM, event Event) (state State) {
    return head.Super()
}

type HSM struct {
    SourceState State
    State       State
    StateTable  map[string]State
}

// initial must set top as parent
func NewHSM(top, initial State) (*HSM, error) {
    hsm := &HSM{
        SourceState: initial,
        State:       top,
        StateTable:  make(map[string]State),
    }
    hsm.StateTable[top.ID()] = top
    return hsm, nil
}

func (hsm *HSM) PrintS() {
    fmt.Println("hsm State: ", hsm.State.ID())
    fmt.Println("hsm SourceState: ", hsm.SourceState.ID())
}

func (hsm *HSM) Init() {
    hsm.init(&StdEvent{EventInit})
}

func (hsm *HSM) init(event Event) {
    // health check
    assert.NotEqual(nil, hsm.State)
    assert.NotEqual(nil, hsm.SourceState)
    // check top state initialized. hsm.State.ID() should be "TOP"
    assert.Equal(hsm.StateTable[TopStateID], hsm.State) // HSM not executed yet
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

func (hsm *HSM) IsIn(stateID string) bool {
    state := hsm.StateTable[stateID]
    return hsm.isIn(state)
}

func (hsm *HSM) isIn(state State) bool {
    // nagivate from current state up to all super state and
    // try to find specified `state'
    for s := hsm.State; s != nil; s = Trigger(hsm, hsm.State, &StdEvent{EventEmpty}) {
        if s == state {
            // a match is found
            return true
        }
    }
    // no match found
    return false
}

func (hsm *HSM) Dispatch(event Event) {
    // Use `SourceState' to record the state which handle the event indeed(which
    // could be super, super-super, ... state).
    // `State' would stay unchange pointing at the current(most concrete) state.
    for hsm.SourceState = hsm.State; hsm.SourceState != nil; {
        hsm.SourceState = Trigger(hsm, hsm.SourceState, event)
    }
}

func Trigger(hsm *HSM, state State, event Event) State {
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

func (hsm *HSM) GetState() State {
    return hsm.State
}

func (hsm *HSM) QInit(targetID string) {
    assert.NotEqual(TopStateID, targetID)
    target := hsm.StateTable[targetID]
    hsm.qinit(target)
}

func (hsm *HSM) qinit(state State) {
    hsm.State = state
}

func (hsm *HSM) QTran(targetID string, event Event) {
    assert.NotEqual(TopStateID, targetID)
    target := hsm.StateTable[targetID]
    hsm.qtran(target, event)
}

func (hsm *HSM) qtran(target State, event Event) {
    var p, q, s State
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

    stateChain := list.New()
    stateChain.Init()

    stateChain.PushBack(nil)
    stateChain.PushBack(target)

    // (a) check `SourceState' == `target' (transition to self)
    if hsm.SourceState == target {
        Trigger(hsm, hsm.SourceState, &StdEvent{EventExit})
        goto inLCA
    }
    // (b) check `SourceState' == `target.Super()'
    p = Trigger(hsm, target, &StdEvent{EventEmpty})
    if hsm.SourceState == p {
        goto inLCA
    }
    // (c) check `SourceState.Super()' == `target.Super()' (most common)
    q = Trigger(hsm, hsm.SourceState, &StdEvent{EventEmpty})
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
        if q == lca.Value {
            // do not enter the LCA
            stateChain.Remove(stateChain.Back())
            goto inLCA
        }
    }
    // (g) check each `SourceState.Super().Super()...' for target...
    for s = q; s != nil; s = Trigger(hsm, s, &StdEvent{EventEmpty}) {
        for lca := stateChain.Back(); lca != nil; lca = lca.Prev() {
            if s == lca.Value {
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
    for e := stateChain.Back(); e != nil && e.Value != nil; {
        s, ok := e.Value.(State)
        assert.Equal(true, ok)
        Trigger(hsm, s, &StdEvent{EventEntry}) // enter `s' state
        stateChain.Remove(stateChain.Back())
        e = stateChain.Back()
    }
    // update current state
    hsm.State = target
    for Trigger(hsm, target, &StdEvent{EventInit}) == nil {
        // initial transition must go *one* level deep
        assert.Equal(s, Trigger(hsm, hsm.State, &StdEvent{EventEmpty}))
        target = hsm.State
        Trigger(hsm, target, &StdEvent{EventEntry})
    }
}

func TruncateList(l *list.List, e *list.Element) *list.List {
    assert.NotEqual(nil, l)
    assert.NotEqual(nil, e)
    // remove `e' and all element after `e' from `l'
    var next *list.Element
    for ; e != nil; e = next {
        next = e.Next()
        l.Remove(e)
    }
    return l
}

type Top struct {
    StateHead
}

func NewTop() (*Top, error) {
    return &Top{MakeStateHead(nil)}, nil
}

func (top *Top) ID() string {
    return TopStateID
}

func (top *Top) Init(hsm *HSM, event Event) (state State) {
    return nil
}

func (top *Top) Entry(hsm *HSM, event Event) (state State) {
    return nil
}

func (top *Top) Exit(hsm *HSM, event Event) (state State) {
    return nil
}

func (top *Top) Handle(hsm *HSM, event Event) (state State) {
    return nil
}

type Initial struct {
    StateHead
}

func NewInitial(super State) (*Initial, error) {
    return &Initial{MakeStateHead(super)}, nil
}

func (*Initial) ID() string {
    return "Initial"
}

func (self *Initial) Init(hsm *HSM, event Event) (state State) {
    fmt.Println("GLOBAL INIT in", self.ID(), ":Init()")
    hsm.QInit("S1")
    return nil
}

func (self *Initial) Handle(hsm *HSM, event Event) (state State) {
    // should never be called
    assert.True(false)
    return self.Super()
}

type S1 struct {
    StateHead
}

func NewS1(super State) (*S1, error) {
    return &S1{MakeStateHead(super)}, nil
}

func (*S1) ID() string {
    return "S1"
}

func (self *S1) Init(hsm *HSM, event Event) (state State) {
    fmt.Println(self.ID(), ":Init event=", event)
    hsm.QInit("S11")
    return nil
}

func (self *S1) Entry(hsm *HSM, event Event) (state State) {
    fmt.Println(self.ID(), ":Entry event=", event)
    return nil
}

func (self *S1) Exit(hsm *HSM, event Event) (state State) {
    fmt.Println(self.ID(), ":Exit event=", event)
    return nil
}

func (self *S1) Handle(hsm *HSM, event Event) (state State) {
    fmt.Println(self.ID(), ":Handle event=", event)
    return self.Super()
}

type S11 struct {
    StateHead
}

func NewS11(super State) (*S11, error) {
    return &S11{MakeStateHead(super)}, nil
}

func (*S11) ID() string {
    return "S11"
}

func (self *S11) Entry(hsm *HSM, event Event) (state State) {
    fmt.Println(self.ID(), ":Entry event=", event)
    return nil
}

func (self *S11) Exit(hsm *HSM, event Event) (state State) {
    fmt.Println(self.ID(), ":Exit event=", event)
    return nil
}

func (self *S11) Handle(hsm *HSM, event Event) (state State) {
    fmt.Println(self.ID(), ":Handle event=", event)
    switch event.Type() {
    case Event1:
        hsm.QTran("S11", event)
        return nil
    case Event2:
        hsm.QTran("S12", event)
        return nil
    }
    return self.Super()
}

type S12 struct {
    StateHead
}

func NewS12(super State) (*S12, error) {
    return &S12{MakeStateHead(super)}, nil
}

func (*S12) ID() string {
    return "S12"
}

func (self *S12) Entry(hsm *HSM, event Event) (state State) {
    fmt.Println(self.ID(), ":Entry event=", event)
    return nil
}

func (self *S12) Exit(hsm *HSM, event Event) (state State) {
    fmt.Println(self.ID(), ":Exit event=", event)
    return nil
}

func (self *S12) Handle(hsm *HSM, event Event) (state State) {
    fmt.Println(self.ID(), ":Handle event=", event)
    switch event.Type() {
    case Event2:
        hsm.QTran("S11", event)
        return nil
    }
    return self.Super()
}

func NewWorld() (*HSM, error) {
    top, _ := NewTop()
    initial, _ := NewInitial(top)
    s1, _ := NewS1(top)
    s11, _ := NewS11(s1)
    s12, _ := NewS12(s1)
    hsm, _ := NewHSM(top, initial)
    hsm.StateTable[initial.ID()] = initial
    hsm.StateTable[s1.ID()] = s1
    hsm.StateTable[s11.ID()] = s11
    hsm.StateTable[s12.ID()] = s12
    hsm.Init()
    return hsm, nil
}
