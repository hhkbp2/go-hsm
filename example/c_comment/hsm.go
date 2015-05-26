package c_comment

import hsm "github.com/hhkbp2/go-hsm"
import "container/list"

const (
	HSMTypeCComment = hsm.HSMTypeStd + 1 + iota
)

// CCommentHSM represent the Hierarchical State Machine for C Comment Parsing.
// This example demonstrate how to parse C code, recording the code chars
// and ignore those chars in /* */ comments, working the similar way as
// a C code Parser reads in the source code and does token scanning.
type CCommentHSM struct {
	*hsm.StdHSM
	// codeCharList is a simple data structure to record all chars for code
	codeCharList *list.List
}

func NewCCommentHSM(top, initial hsm.State) *CCommentHSM {
	return &CCommentHSM{
		StdHSM:       hsm.NewStdHSM(HSMTypeCComment, top, initial),
		codeCharList: list.New(),
	}
}

// Init() is part of HSM interface.
func (self *CCommentHSM) Init() {
	self.StdHSM.Init2(self, hsm.StdEvents[hsm.EventInit])
}

// Dispatch() is part of HSM interface.
func (self *CCommentHSM) Dispatch(event hsm.Event) {
	self.StdHSM.Dispatch2(self, event)
}

func (self *CCommentHSM) QTran(targetStateID string) {
	target := self.StdHSM.LookupState(targetStateID)
	self.StdHSM.QTranHSM(self, target)
}

func (self *CCommentHSM) QTranOnEvent(targetStateID string, event hsm.Event) {
	target := self.StdHSM.LookupState(targetStateID)
	self.StdHSM.QTranHSMOnEvents(self, target, event, event, event)
}

// CurrentStateID() is a abstract leaking method to get current state.
// Just for example demonstration.
func (self *CCommentHSM) CurrentStateID() string {
	return self.StdHSM.State.ID()
}

// ProcessCodeChar() is to record all chars of the code.
func (self *CCommentHSM) ProcessCodeChar(c byte) {
	self.codeCharList.PushBack(c)
}

// TraverseCode() is the interface of traversing all chars of the code.
func (self *CCommentHSM) TraverseCode(f func(value interface{}) interface{}) {
	MapOnList(f, self.codeCharList)
}
