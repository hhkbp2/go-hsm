package c_comment

import "log"
import "container/list"

func Logln(v ...interface{}) {
	log.Println(v...)
}

// Call func f on every element of list lst.
// Like a map() style function on lisp, func f must get
// the type of element right.
// To achieve generalization, interface{}(void pointer) is used in Golang.
func MapOnList(
	f func(value interface{}) interface{},
	lst *list.List) *list.List {

	result := list.New()
	for e := lst.Front(); e != nil; e = e.Next() {
		result.PushBack(f(e.Value))
	}
	return result
}
