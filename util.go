package hsm

import "container/list"
import "github.com/hhkbp2/go-hsm/assert"

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
