package hsm

import "container/list"
import "errors"

// Trigger() is a helper function to dispatch event of different types to
// the corresponding method.
func Trigger(hsm HSM, state State, event Event) State {
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

func TriggerInit(hsm HSM, state State, event Event) State {
	return state.Init(hsm, event)
}

func TriggerEntry(hsm HSM, state State, event Event) State {
	return state.Entry(hsm, event)
}

func TriggerExit(hsm HSM, state State, event Event) State {
	return state.Exit(hsm, event)
}

// ListTruncate() removes elements from `e' to the last element in list `l'.
// The range to be removed is [e, l.Back()]. It returns list `l'.
func ListTruncate(l *list.List, e *list.Element) *list.List {
	AssertNotEqual(nil, l)
	AssertNotEqual(nil, e)
	// remove `e' and all elements after `e'
	var next *list.Element
	for ; e != nil; e = next {
		next = e.Next()
		l.Remove(e)
	}
	return l
}

// ListFind() searchs the first element which has the same value of `value' in
// list `l'. It uses object comparation for equality check.
func ListFind(l *list.List, value interface{}) (*list.Element, error) {
	predicate := func(v interface{}) bool {
		return ObjectAreEqual(value, v)
	}
	return ListFindIf(l, predicate)
}

// ListFindIf() searchs for and element of the list `l' that
// satifies the predicate function `predicate'.
func ListFindIf(l *list.List, predicate func(value interface{}) bool) (*list.Element, error) {
	for e := l.Front(); e != nil; e = e.Next() {
		if predicate(e.Value) {
			return e, nil
		}
	}
	return nil, errors.New("find no match")
}

// ListIn() tests whether `value' is in list `l'.
func ListIn(l *list.List, value interface{}) bool {
	e, err := ListFind(l, value)
	if e == nil && err != nil {
		return false
	}
	return true
}
