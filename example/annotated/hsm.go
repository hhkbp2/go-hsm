package annotated

import hsm "github.com/hhkbp2/go-hsm"

const (
	HSMTypeAnnotated hsm.HSMType = hsm.HSMTypeUser
)

type AnnotatedHSM struct {
	*hsm.StdHSM

	foo bool
}

func NewAnnotatedHSM(top, initial hsm.State) *AnnotatedHSM {
	return &AnnotatedHSM{
		StdHSM: hsm.NewStdHSM(HSMTypeAnnotated, top, initial),
	}
}

func (self *AnnotatedHSM) Init() {
	self.StdHSM.Init2(self, hsm.StdEvents[hsm.EventInit])
}

func (self *AnnotatedHSM) Dispatch(event hsm.Event) {
	self.StdHSM.Dispatch2(self, event)
}

func (self *AnnotatedHSM) QTran(targetStateID string) {
	target := self.StdHSM.LookupState(targetStateID)
	self.StdHSM.QTranHSM(self, target)
}

func (self *AnnotatedHSM) CurrentStateID() string {
	return self.StdHSM.State.ID()
}

func (self *AnnotatedHSM) GetFoo() bool {
	return self.foo
}

func (self *AnnotatedHSM) SetFoo(newFoo bool) {
	self.foo = newFoo
}
