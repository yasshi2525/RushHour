package entities

import (
	"fmt"
)

// RailLine represents how Train should run.
type RailLine struct {
	Base
	Persistence

	Name      string `json:"name"`
	AutoExt   bool   `json:"auto_ext"`
	AutoPass  bool   `json:"auto_pass"`
	ReRouting bool   `gorm:"-" json:"-"`

	RailEdges map[uint]*RailEdge `gorm:"-" json:"-"`
	Stops     map[uint]*Platform `gorm:"-" json:"-"`
	Tasks     map[uint]*LineTask `gorm:"-" json:"-"`
	Trains    map[uint]*Train    `gorm:"-" json:"-"`
	Steps     map[uint]*Step     `gorm:"-" json:"-"`
}

// NewRailLine create instance
func (m *Model) NewRailLine(o *Player) *RailLine {
	l := &RailLine{
		Base:        m.NewBase(RAILLINE, o),
		Persistence: NewPersistence(),
	}
	l.Init(m)
	l.Resolve(o)
	l.Marshal()
	o.Resolve(l)
	m.Add(l)
	return l
}

// B returns base information of this elements.
func (l *RailLine) B() *Base {
	return &l.Base
}

// P returns time information for database.
func (l *RailLine) P() *Persistence {
	return &l.Persistence
}

// StartPlatform creates LineTask which depart from specified Platform.
func (l *RailLine) StartPlatform(p *Platform) (*LineTask, *LineTask) {
	if len(l.Tasks) > 0 {
		panic(fmt.Errorf("try to start from already built RailLine %v", l))
	}
	if l.AutoExt {
		rn := p.OnRailNode
		if tk := rn.Tracks[rn.ID]; tk != nil {
			return l.StartEdge(tk.Via)
		}
	}
	if l.AutoPass {
		return nil, nil
	}
	tail := l.M.NewLineTaskDept(l, p)
	return tail, tail
}

// StartEdge creates LineTask which depart from specified RailEdge.
func (l *RailLine) StartEdge(re *RailEdge) (*LineTask, *LineTask) {
	if len(l.Tasks) > 0 {
		panic(fmt.Errorf("try to start from built RailLine: %v", l))
	}
	var head, tail *LineTask
	if p := re.FromNode.OverPlatform; !l.AutoPass && p != nil {
		head = l.M.NewLineTaskDept(l, p)
		tail = l.M.NewLineTask(l, re, head)
	} else {
		head = l.M.NewLineTask(l, re)
		tail = head
	}
	if l.AutoExt {
		tail.Stretch(re.Reverse).SetNext(head)
		return head, nil
	}
	return head, tail
}

// Complement connects head and tail with minimum distance route.
func (l *RailLine) Complement() {
	if len(l.Tasks) == 0 || l.IsRing() {
		panic(fmt.Errorf("try to complement empty or ring RailLine: %v", l))
	}
	head, tail := l.Borders()
	tk := tail.ToNode().Tracks[head.FromNode().ID]
	for tk != nil && !l.CanRing() {
		tail = tail.Stretch(tk.Via)
		tk = tail.ToNode().Tracks[head.FromNode().ID]
	}
}

// RingIf connects head and tail if can.
func (l *RailLine) RingIf() bool {
	if l.CanRing() {
		head, tail := l.Borders()
		tail.SetNext(head)
		return true
	}
	return false
}

// ClearTransports eraces Transport information.
func (l *RailLine) ClearTransports() {
	for _, p := range l.Stops {
		p.Transports = make(map[uint]*Transport)
	}
}

// Init makes map
func (l *RailLine) Init(m *Model) {
	l.Base.Init(RAILLINE, m)
	l.RailEdges = make(map[uint]*RailEdge)
	l.Stops = make(map[uint]*Platform)
	l.Tasks = make(map[uint]*LineTask)
	l.Trains = make(map[uint]*Train)
	l.Steps = make(map[uint]*Step)
}

// Resolve set reference
func (l *RailLine) Resolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			l.O = obj
			obj.Resolve(l)
		case *RailEdge:
			l.RailEdges[obj.ID] = obj
		case *Platform:
			l.Stops[obj.ID] = obj
		case *LineTask:
			l.Tasks[obj.ID] = obj
		case *Train:
			l.Trains[obj.ID] = obj
		case *Step:
			l.Steps[obj.ID] = obj
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	l.Marshal()
}

// Marshal set if from reference
func (l *RailLine) Marshal() {
}

// UnMarshal set reference from id.
func (l *RailLine) UnMarshal() {
	l.Resolve(l.M.Find(PLAYER, l.OwnerID))
}

// UnResolve unregisters specified refernce.
func (l *RailLine) UnResolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *RailEdge:
			delete(l.RailEdges, obj.ID)
		case *LineTask:
			delete(l.Tasks, obj.ID)
		case *Train:
			delete(l.Trains, obj.ID)
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
}

// CheckDelete check remain relation.
func (l *RailLine) CheckDelete() error {
	return nil
}

// BeforeDelete remove reference of related entity
func (l *RailLine) BeforeDelete() {
	l.O.UnResolve(l)
}

// Delete removes this entity with related ones.
func (l *RailLine) Delete() {
	for _, t := range l.Trains {
		t.SetTask(nil)
	}
	for _, lt := range l.Tasks {
		lt.Delete()
	}
	l.M.Delete(l)
}

// Borders returns head and tail of LineTask.
// Head and tail are nil when LineTask loops
// Tail is undirecting LineTask, that is LineTask.Next is nil
// Head is undirected  LineTask because head of chain is what any other doesn't target
func (l *RailLine) Borders() (*LineTask, *LineTask) {
	var head, tail *LineTask
	for _, lt := range l.Tasks {
		if lt.Before() == nil {
			head = lt
		}
		if lt.Next() == nil {
			tail = lt
		}
		if head != nil && tail != nil {
			return head, tail
		}
	}
	// looped
	return nil, nil
}

// IsRing returns whether LineTask is looping or not
func (l *RailLine) IsRing() bool {
	if len(l.Tasks) <= 1 {
		return false
	}
	h, t := l.Borders()
	return h == nil && t == nil
}

// CanRing returns head and tail is adjustable.
func (l *RailLine) CanRing() bool {
	head, tail := l.Borders()

	// ringed loop can't loop any more
	if head == nil && tail == nil {
		return false
	}

	switch tail.TaskType {
	case OnDeparture:
		// allow: dept -> not dept
		return head.TaskType != OnDeparture && tail.Stay == head.Dept
	case OnMoving:
		fallthrough
	case OnPassing:
		// allow: move/pass -> not dept
		return head.TaskType != OnDeparture && tail.Moving.ToNode == head.Moving.FromNode
	case OnStopping:
		// allow: stop -> dept
		return head.TaskType == OnDeparture && tail.Dest == head.Stay
	default:
		panic(fmt.Errorf("invalid type = %v", tail.TaskType))
	}
}

// String represents status
func (l *RailLine) String() string {
	l.Marshal()
	ostr := ""
	if l.O != nil {
		ostr = fmt.Sprintf(":%s", l.O.Short())
	}
	return fmt.Sprintf("%s(%d):lt=%d%s:%s", l.Type().Short(),
		l.ID, len(l.Tasks), ostr, l.Name)
}
