package main

import "fmt"
import . "github.com/hhkbp2/go-hsm"

const (
    Event1 = EventUser + iota
    Event2
    Event3
)

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

func main() {
    m, _ := NewWorld()
    events := []Event{
        &StdEvent{Event2},
        &StdEvent{Event1},
        &StdEvent{Event1},
        &StdEvent{Event2},
    }
    for _, e := range events {
        fmt.Println("> dispatch event:", e)
        m.Dispatch(e)
    }
}
