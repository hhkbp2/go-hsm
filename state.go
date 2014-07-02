package hsm

const (
    TopStateID     = "TOP"
    InitialStateID = "Initial"
)

type State interface {
    ID() string

    Super() (super State)

    Init(hsm HSM, event Event) (state State)
    Entry(hsm HSM, event Event) (state State)
    Exit(hsm HSM, event Event) (state State)
    Handle(hsm HSM, event Event) (state State)
}

type StateHead struct {
    super State
}

func MakeStateHead(super State) StateHead {
    return StateHead{
        super: super,
    }
}

func (self *StateHead) Super() State {
    return self.super
}

func (self *StateHead) Init(hsm HSM, event Event) (state State) {
    return self.Super()
}

func (self *StateHead) Entry(hsm HSM, event Event) (state State) {
    return self.Super()
}

func (self *StateHead) Exit(hsm HSM, event Event) (state State) {
    return self.Super()
}

type Top struct {
    StateHead
}

func NewTop() (*Top, error) {
    return &Top{MakeStateHead(nil)}, nil
}

func (self *Top) ID() string {
    return TopStateID
}

func (self *Top) Init(hsm HSM, event Event) (state State) {
    return nil
}

func (self *Top) Entry(hsm HSM, event Event) (state State) {
    return nil
}

func (self *Top) Exit(hsm HSM, event Event) (state State) {
    return nil
}

func (self *Top) Handle(hsm HSM, event Event) (state State) {
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

func (self *Initial) Init(hsm HSM, event Event) (state State) {
    hsm.QInit("S1")
    return nil
}

func (self *Initial) Handle(hsm HSM, event Event) (state State) {
    // should never be called
    return self.Super()
}
