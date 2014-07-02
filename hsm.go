package hsm

import "container/list"

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

func (hsm *HSM) Init() {
    hsm.init(&StdEvent{EventInit})
}

func (hsm *HSM) init(event Event) {
    // health check
    AssertNotEqual(nil, hsm.State)
    AssertNotEqual(nil, hsm.SourceState)
    // check top state initialized. hsm.State.ID() should be "TOP"
    AssertEqual(hsm.StateTable[TopStateID], hsm.State) // HSM not executed yet
    // save State in a temporary
    s := hsm.State
    // top-most initial transition
    Trigger(hsm, hsm.SourceState, event)
    // initial transition must go *one* level deep
    AssertEqual(s, Trigger(hsm, hsm.State, &StdEvent{EventEmpty}))
    // update the termporary
    s = hsm.State
    // enter the state
    Trigger(hsm, s, &StdEvent{EventEntry})
    for Trigger(hsm, s, &StdEvent{EventInit}) == nil { // init handled?
        // initial transition must go *one* level deep
        AssertEqual(s, Trigger(hsm, hsm.State, &StdEvent{EventEmpty}))
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

func (hsm *HSM) GetState() State {
    return hsm.State
}

func (hsm *HSM) QInit(targetID string) {
    AssertNotEqual(TopStateID, targetID)
    target := hsm.StateTable[targetID]
    hsm.qinit(target)
}

func (hsm *HSM) qinit(state State) {
    hsm.State = state
}

func (hsm *HSM) QTran(targetID string, event Event) {
    AssertNotEqual(TopStateID, targetID)
    target := hsm.StateTable[targetID]
    hsm.qtran(target, event)
}

func (hsm *HSM) qtran(target State, event Event) {
    var p, q, s State
    for s := hsm.State; s != hsm.SourceState; {
        // we are about to dereference `s'
        AssertNotEqual(nil, s)
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
    AssertTrue(false)
inLCA: // now we are in the LCA of `SourceState' and `target'
    // retrace the entry path in reverse order
    for e := stateChain.Back(); e != nil && e.Value != nil; {
        s, ok := e.Value.(State)
        AssertEqual(true, ok)
        Trigger(hsm, s, &StdEvent{EventEntry}) // enter `s' state
        stateChain.Remove(stateChain.Back())
        e = stateChain.Back()
    }
    // update current state
    hsm.State = target
    for Trigger(hsm, target, &StdEvent{EventInit}) == nil {
        // initial transition must go *one* level deep
        AssertEqual(s, Trigger(hsm, hsm.State, &StdEvent{EventEmpty}))
        target = hsm.State
        Trigger(hsm, target, &StdEvent{EventEntry})
    }
}
