package hsm

import "github.com/hhkbp2/go-hsm/assert"

const (
    TopStateID     = "TOP"
    InitialStateID = "Initial"
)

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
    return InitialStateID
}

func (self *Initial) Init(hsm *HSM, event Event) (state State) {
    hsm.QInit("S1")
    return nil
}

func (self *Initial) Handle(hsm *HSM, event Event) (state State) {
    // should never be called
    assert.True(false)
    return self.Super()
}
