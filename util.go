package hsm

import "container/list"
import "errors"

func Trigger(hsm HSM, state State, event Event) State {
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

func ListTruncate(l *list.List, e *list.Element) *list.List {
    AssertNotEqual(nil, l)
    AssertNotEqual(nil, e)
    // remove `e' and all element after `e' from `l'
    var next *list.Element
    for ; e != nil; e = next {
        next = e.Next()
        l.Remove(e)
    }
    return l
}

func ListFind(l *list.List, value interface{}) (*list.Element, error) {
    predicate := func(v interface{}) bool {
        return ObjectAreEqual(value, v)
    }
    return ListFindIf(l, predicate)
}

func ListFindIf(l *list.List, predicate func(value interface{}) bool) (*list.Element, error) {
    for e := l.Front(); e != nil; e = e.Next() {
        if predicate(e.Value) {
            return e, nil
        }
    }
    return nil, errors.New("find no match")
}

func ListIn(l *list.List, value interface{}) bool {
    e, err := ListFind(l, value)
    if e == nil && err != nil {
        return false
    }
    return true
}
