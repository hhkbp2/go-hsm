package hsm

import "container/list"

const (
    HSMTypeStd = iota
)

type HSM interface {
    Type() uint32

    Init()
    Dispatch(event Event)

    GetState() State
    IsIn(stateID string) bool

    QInit(targetStateID string)
    QTran(targetStateID string)
}

type StdHSM struct {
    HSMType     uint32
    SourceState State
    State       State
    StateTable  map[string]State
}

// initial must set top as parent
func NewStdHSM(top, initial State) *StdHSM {
    AssertEqual(TopStateID, top.ID())
    AssertEqual(InitialStateID, initial.ID())
    hsm := &StdHSM{
        HSMType:     HSMTypeStd,
        SourceState: initial,
        State:       top,
        StateTable:  make(map[string]State),
    }
    hsm.StateTable[top.ID()] = top
    // setup state table
    hsm.setupStateTable()
    return hsm
}

func (self *StdHSM) Type() uint32 {
    return self.HSMType
}

func (self *StdHSM) setupStateTable() {
    for traverse_queue := self.State.Children(); len(traverse_queue) != 0; {
        state := traverse_queue[0]
        traverse_queue = traverse_queue[1:]
        _, ok := self.StateTable[state.ID()]
        AssertFalse(ok)
        self.StateTable[state.ID()] = state
        children := state.Children()
        for _, state := range children {
            traverse_queue = append(traverse_queue, state)
        }
    }
}

func (self *StdHSM) Init() {
    self.init(&StdEvent{EventInit})
}

func (self *StdHSM) init(event Event) {
    // health check
    AssertNotEqual(nil, self.State)
    AssertNotEqual(nil, self.SourceState)
    // check HSM is not executed yet
    AssertEqual(self.StateTable[TopStateID], self.State)
    AssertEqual(self.StateTable[InitialStateID], self.SourceState)
    // save State in a temporary
    s := self.State
    // top-most initial transition
    Trigger(self, self.SourceState, event)
    // initial transition must go *one* level deep
    AssertEqual(s, Trigger(self, self.State, &StdEvent{EventEmpty}))
    // update the termporary
    s = self.State
    // enter the state
    Trigger(self, s, &StdEvent{EventEntry})
    for Trigger(self, s, &StdEvent{EventInit}) == nil { // init handled?
        // initial transition must go *one* level deep
        AssertEqual(s, Trigger(self, self.State, &StdEvent{EventEmpty}))
        s = self.State
        // enter the substate
        Trigger(self, s, &StdEvent{EventEntry})
    }
    // we are in well-initialized state now
}

func (self *StdHSM) Dispatch(event Event) {
    // Use `SourceState' to record the state which handle the event indeed(which
    // could be super, super-super, ... state).
    // `State' would stay unchange pointing at the current(most concrete) state.
    for self.SourceState = self.State; self.SourceState != nil; {
        self.SourceState = Trigger(self, self.SourceState, event)
    }
}

func (self *StdHSM) GetState() State {
    return self.State
}
func (self *StdHSM) IsIn(stateID string) bool {
    state := self.StateTable[stateID]
    return self.isIn(state)
}

func (self *StdHSM) isIn(state State) bool {
    // nagivate from current state up to all super state and
    // try to find specified `state'
    for s := self.State; s != nil; s = Trigger(self, self.State, &StdEvent{EventEmpty}) {
        if s == state {
            // a match is found
            return true
        }
    }
    // no match found
    return false
}

func (self *StdHSM) QInit(targetStateID string) {
    AssertNotEqual(TopStateID, targetStateID)
    target, ok := self.StateTable[targetStateID]
    AssertTrue(ok)
    self.qinit(target)
}

func (self *StdHSM) qinit(state State) {
    self.State = state
}

func (self *StdHSM) QTran(targetStateID string) {
    AssertNotEqual(TopStateID, targetStateID)
    target, ok := self.StateTable[targetStateID]
    AssertTrue(ok)
    self.qtran(target)
}

func (self *StdHSM) qtran(target State) {
    var p, q, s State
    for s := self.State; s != self.SourceState; {
        // we are about to dereference `s'
        AssertNotEqual(nil, s)
        t := Trigger(self, s, &StdEvent{EventExit})
        if t != nil {
            s = t
        } else {
            s = Trigger(self, s, &StdEvent{EventEmpty})
        }
    }

    stateChain := list.New()
    stateChain.Init()

    stateChain.PushBack(nil)
    stateChain.PushBack(target)

    // (a) check `SourceState' == `target' (transition to self)
    if self.SourceState == target {
        Trigger(self, self.SourceState, &StdEvent{EventExit})
        goto inLCA
    }
    // (b) check `SourceState' == `target.Super()'
    p = Trigger(self, target, &StdEvent{EventEmpty})
    if self.SourceState == p {
        goto inLCA
    }
    // (c) check `SourceState.Super()' == `target.Super()' (most common)
    q = Trigger(self, self.SourceState, &StdEvent{EventEmpty})
    if q == p {
        Trigger(self, self.SourceState, &StdEvent{EventExit})
        goto inLCA
    }
    // (d) check `SourceState.Super()' == `target'
    if q == target {
        Trigger(self, self.SourceState, &StdEvent{EventExit})
        stateChain.Remove(stateChain.Back())
        goto inLCA
    }
    // (e) check rest of `SourceState' == `target.Super().Super()...'  hierarchy
    stateChain.PushBack(p)
    s = Trigger(self, p, &StdEvent{EventEmpty})
    for s != nil {
        if self.SourceState == s {
            goto inLCA
        }
        stateChain.PushBack(s)
        s = Trigger(self, s, &StdEvent{EventEmpty})
    }
    // exit source state
    Trigger(self, self.SourceState, &StdEvent{EventExit})
    // (f) check rest of `SourceState.Super()' == `target.Super().Super()...'
    for lca := stateChain.Back(); lca != nil; lca = lca.Prev() {
        if q == lca.Value {
            // do not enter the LCA
            stateChain.Remove(stateChain.Back())
            goto inLCA
        }
    }
    // (g) check each `SourceState.Super().Super()...' for target...
    for s = q; s != nil; s = Trigger(self, s, &StdEvent{EventEmpty}) {
        for lca := stateChain.Back(); lca != nil; lca = lca.Prev() {
            if s == lca.Value {
                stateChain = ListTruncate(stateChain, lca)
                goto inLCA
            }
        }
        Trigger(self, s, &StdEvent{EventExit})
    }
    // malformed HSM
    AssertTrue(false)
inLCA: // now we are in the LCA of `SourceState' and `target'
    // retrace the entry path in reverse order
    for e := stateChain.Back(); e != nil && e.Value != nil; {
        s, ok := e.Value.(State)
        AssertEqual(true, ok)
        Trigger(self, s, &StdEvent{EventEntry}) // enter `s' state
        stateChain.Remove(stateChain.Back())
        e = stateChain.Back()
    }
    // update current state
    self.State = target
    for Trigger(self, target, &StdEvent{EventInit}) == nil {
        // initial transition must go *one* level deep
        AssertEqual(s, Trigger(self, self.State, &StdEvent{EventEmpty}))
        target = self.State
        Trigger(self, target, &StdEvent{EventEntry})
    }
}
