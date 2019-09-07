package entities

import (
	"fmt"
	"math"
)

// LineTaskType represents the state what Train should do now.
type LineTaskType uint

const (
	// OnDeparture represents the state that Train waits for departure in Station.
	OnDeparture LineTaskType = iota + 1
	// OnMoving represents the state that Train runs to next RailNode.
	OnMoving
	// OnStopping represents the state that Train stops to next Station.
	OnStopping
	// OnPassing represents the state that Train passes to next Station.
	OnPassing
)

// LineTask is the element of Line.
// The chain of LineTask represents Line structure.
type LineTask struct {
	Base
	Persistence

	TaskType LineTaskType `gorm:"not null" json:"type"`

	Moving    *RailEdge `gorm:"-" json:"-"`
	Stay      *Platform `gorm:"-" json:"-"`
	Dept      *Platform `gorm:"-" json:"-"`
	Dest      *Platform `gorm:"-" json:"-"`
	RailLine  *RailLine `gorm:"-"        json:"-"`
	before    *LineTask
	next      *LineTask
	Trains    map[uint]*Train `gorm:"-" json:"-"`
	OverSteps map[uint]*Step  `gorm:"-" json:"-"`

	RailLineID uint `gorm:"not null" json:"lid"`
	BeforeID   uint `gorm:"-"        json:"before,omitempty"`
	NextID     uint `                json:"next,omitempty"`
	StayID     uint `                json:"psid,omitempty"`
	DeptID     uint `gorm:"-"        json:"p1id,omitempty"`
	MovingID   uint `                json:"reid,omitempty"`
	DestID     uint `gorm:"-"        json:"p2id,omitempty"`
}

// NewLineTaskDept create "dept"
func (m *Model) NewLineTaskDept(l *RailLine, p *Platform, tail ...*LineTask) *LineTask {
	lt := &LineTask{
		Base:        m.NewBase(LINETASK, l.O),
		Persistence: NewPersistence(),
		TaskType:    OnDeparture,
	}
	lt.Init(m)
	lt.Resolve(l.O, l, p)
	lt.Marshal()
	if len(tail) > 0 && tail[0] != nil {
		tail[0].SetNext(lt)
	}
	m.Add(lt)
	lt.RailLine.ReRouting = true
	return lt
}

// NewLineTask creates "stop" or "pass" or "moving"
func (m *Model) NewLineTask(l *RailLine, re *RailEdge, tail ...*LineTask) *LineTask {
	lt := &LineTask{
		Base:        m.NewBase(LINETASK, l.O),
		Persistence: NewPersistence(),
	}
	lt.Init(m)
	lt.Resolve(l.O, l, re)
	lt.Marshal()

	if re.ToNode.OverPlatform == nil {
		lt.TaskType = OnMoving
	} else {
		if lt.RailLine.AutoPass {
			lt.TaskType = OnPassing
		} else {
			lt.TaskType = OnStopping
		}
	}

	if len(tail) > 0 && tail[0] != nil {
		tail[0].SetNext(lt)
	}
	m.Add(lt)
	lt.RailLine.ReRouting = true
	return lt
}

// B returns base information of this elements.
func (lt *LineTask) B() *Base {
	return &lt.Base
}

// P returns time information for database.
func (lt *LineTask) P() *Persistence {
	return &lt.Persistence
}

// Init initializes map
func (lt *LineTask) Init(m *Model) {
	lt.Base.Init(LINETASK, m)
	lt.Trains = make(map[uint]*Train)
	lt.OverSteps = make(map[uint]*Step)
}

// Step procceed progress to certain time.
func (lt *LineTask) Step(prog *float64, sec *float64) {
	canDist := *sec * Const.Train.Speed
	remainDist := (1.0 - *prog) * lt.Cost()
	if remainDist < canDist {
		*sec -= remainDist / Const.Train.Speed
		*prog = 1.0
	} else {
		*prog += *sec * Const.Train.Speed / lt.Cost()
		*sec = 0
	}
}

// Loc returns Point which devides progress ratio to it.
func (lt *LineTask) Loc(prog float64) *Point {
	if lt.TaskType == OnDeparture {
		return lt.Stay.Pos()
	}
	if prog < 0.5 && lt.before.TaskType == OnDeparture {
		return lt.Moving.Div(2 * prog * prog)
	} else if prog > 0.5 && lt.TaskType == OnStopping {
		return lt.Moving.Div(-2*prog*prog + 4*prog - 1)
	}
	return lt.Moving.Div(prog)
}

// Resolve set reference
func (lt *LineTask) Resolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			lt.O = obj
			obj.Resolve(lt)
		case *Platform:
			lt.Stay = obj
			lt.Dept = obj
			lt.Dest = obj
			lt.RailLine.Resolve(obj)
			obj.Resolve(lt)
			obj.OnRailNode.OutTasks[lt.ID] = lt
			obj.OnRailNode.InTasks[lt.ID] = lt
		case *RailEdge:
			lt.Moving = obj
			lt.RailLine.Resolve(obj)
			obj.Resolve(lt)
			if p := obj.FromNode.OverPlatform; p != nil {
				lt.Dept = p
				p.Resolve(lt)
			}
			if p := obj.ToNode.OverPlatform; p != nil {
				lt.Dest = p
				p.Resolve(lt)
			}
			obj.FromNode.OutTasks[lt.ID] = lt
			obj.ToNode.InTasks[lt.ID] = lt
		case *RailLine:
			lt.RailLine = obj
			obj.Resolve(lt)
		case *LineTask:
			lt.next = obj
			if obj != nil {
				obj.SetBefore(lt)
			}
		case *Train:
			lt.Trains[obj.ID] = obj
			lt.RailLine.Resolve(obj)
			switch lt.TaskType {
			case OnDeparture:
				lt.Stay.Resolve(obj)
			default:
				lt.Moving.Resolve(obj)
			}
		case *Step:
			lt.OverSteps[obj.ID] = obj
			lt.RailLine.Resolve(obj)
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	lt.Marshal()
}

// UnResolve unregisters specified refernce.
func (lt *LineTask) UnResolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Train:
			delete(lt.Trains, obj.ID)
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
}

// Marshal set id from reference
func (lt *LineTask) Marshal() {
	if lt.RailLine != nil {
		lt.RailLineID = lt.RailLine.ID
	}
	if lt.before != nil {
		lt.BeforeID = lt.before.ID
	}
	if lt.next != nil {
		lt.NextID = lt.next.ID
	}
	if lt.Moving != nil {
		lt.MovingID = lt.Moving.ID
	}
	if lt.Stay != nil {
		lt.StayID = lt.Stay.ID
	}
	if lt.Dept != nil {
		lt.DeptID = lt.Dept.ID
	} else {
		lt.DeptID = ZERO
	}
	if lt.Dest != nil {
		lt.DestID = lt.Dest.ID
	} else {
		lt.DestID = ZERO
	}
}

// UnMarshal set reference from id.
func (lt *LineTask) UnMarshal() {
	lt.Resolve(
		lt.M.Find(PLAYER, lt.OwnerID),
		lt.M.Find(RAILLINE, lt.RailLineID))
	// nullable fields
	if lt.NextID != ZERO {
		lt.Resolve(lt.M.Find(LINETASK, lt.NextID))
	}
	if lt.StayID != ZERO {
		lt.Resolve(lt.M.Find(PLATFORM, lt.StayID))
	}
	if lt.MovingID != ZERO {
		lt.Resolve(lt.M.Find(RAILEDGE, lt.MovingID))
	}
}

// BeforeDelete remove related refernce
func (lt *LineTask) BeforeDelete() {
	if lt.Stay != nil {
		delete(lt.Stay.OnRailNode.OutTasks, lt.ID)
		delete(lt.Stay.OnRailNode.InTasks, lt.ID)
		lt.Stay.UnResolve(lt)
	}
	if lt.Moving != nil {
		delete(lt.Moving.FromNode.OutTasks, lt.ID)
		delete(lt.Moving.ToNode.InTasks, lt.ID)
		lt.Moving.UnResolve(lt)
	}
	if lt.before != nil && lt.before.next == lt {
		lt.before.SetNext(nil)
	}
	if lt.next != nil && lt.next.before == lt {
		lt.next.SetBefore(nil)
	}
	lt.RailLine.UnResolve(lt)
	lt.O.UnResolve(lt)
}

// Delete removes this entity with related ones.
func (lt *LineTask) Delete() {
	for _, t := range lt.Trains {
		t.SetTask(lt.next)
	}
	lt.RailLine.ReRouting = true
	lt.M.Delete(lt)
}

// CheckDelete check remain relation.
func (lt *LineTask) CheckDelete() error {
	return nil
}

// From represents start point
func (lt *LineTask) From() Entity {
	return lt.FromNode()
}

// To represents end point
func (lt *LineTask) To() Entity {
	return lt.ToNode()
}

// FromNode represents start point
func (lt *LineTask) FromNode() *RailNode {
	switch lt.TaskType {
	case OnDeparture:
		return lt.Stay.OnRailNode
	default:
		return lt.Moving.FromNode
	}
}

// ToNode represents end point
func (lt *LineTask) ToNode() *RailNode {
	switch lt.TaskType {
	case OnDeparture:
		return lt.Stay.OnRailNode
	default:
		return lt.Moving.ToNode
	}
}

// Cost represents how many seconds it takes.
func (lt *LineTask) Cost() float64 {
	switch lt.TaskType {
	case OnDeparture:
		var sum float64
		for _, oth := range lt.RailLine.Tasks {
			if oth.TaskType != OnDeparture {
				sum += oth.Cost()
			}
		}
		if length := len(lt.RailLine.Trains); length != 0 {
			return sum / Const.Train.Speed / float64(length)
		}
		return math.MaxFloat64
	default:
		cost := lt.Moving.Cost()
		if lt.before.TaskType == OnDeparture {
			cost += 0.5 * lt.Moving.Cost() * Const.Train.Slowness
		} else {
			cost += 0.5 * lt.Moving.Cost() * Const.Train.Slowness * lt.before.Moving.Angle(lt.Moving) / math.Pi
		}
		if lt.TaskType == OnStopping {
			cost += 0.5 * lt.Moving.Cost() * Const.Train.Slowness
		} else {
			cost += 0.5 * lt.Moving.Cost() * Const.Train.Slowness * lt.Moving.Angle(lt.next.Moving) / math.Pi
		}
		return cost
	}
}

// Before return before field
func (lt *LineTask) Before() *LineTask {
	return lt.before
}

// Next return next field
func (lt *LineTask) Next() *LineTask {
	return lt.next
}

// SetNext changes self changed status for backup
func (lt *LineTask) SetNext(v *LineTask) {
	if lt == v {
		panic(fmt.Errorf("try to self loop: %v", lt))
	}
	if lt.next != nil {
		lt.next.SetBefore(nil)
		lt.next = nil
		lt.NextID = ZERO
	}
	lt.next = v
	if v != nil {
		if lt.TaskType == OnPassing && v.TaskType == OnDeparture {
			panic(fmt.Errorf("try to set Dept to Pass : %v -> %v", v, lt))
		}
		if lt.ToNode() != v.FromNode() {
			panic(fmt.Errorf("try to set far task : %v -> %v", v, lt))
		}
		lt.NextID = v.ID
		v.SetBefore(lt)
	} else {
		lt.NextID = ZERO
	}
	lt.RailLine.ReRouting = true
	lt.Change()
}

// SetDept set Platform to departure.
func (lt *LineTask) SetDept(p *Platform) {
	lt.Dept = p
	if p != nil {
		lt.DeptID = p.ID
	} else {
		lt.DeptID = ZERO
	}
}

// SetDest set Platform to destination.
func (lt *LineTask) SetDest(p *Platform) {
	lt.Dest = p
	if p != nil {
		lt.DestID = p.ID
	} else {
		lt.DestID = ZERO
	}
}

// SetBefore set LineTask to before.
func (lt *LineTask) SetBefore(v *LineTask) {
	lt.before = v
	if v != nil {
		lt.BeforeID = v.ID
	} else {
		lt.BeforeID = ZERO
	}
}

// String represents status
func (lt *LineTask) String() string {
	lt.Marshal()
	ostr := ""
	if lt.O != nil {
		ostr = fmt.Sprintf(":%s", lt.O.Short())
	}
	before, next, stay, dept, moving, dest := "", "", "", "", "", ""
	if lt.before != nil {
		before = fmt.Sprintf(",before=%d", lt.before.ID)
	}
	if lt.next != nil {
		next = fmt.Sprintf(",next=%d", lt.next.ID)
	}
	if lt.Stay != nil {
		stay = fmt.Sprintf(",stay=%d", lt.Stay.ID)
	}
	if lt.Dept != nil {
		dept = fmt.Sprintf(",dept=%d", lt.Dept.ID)
	}
	if lt.Moving != nil {
		moving = fmt.Sprintf(",re=%d", lt.Moving.ID)
	}
	if lt.Dest != nil {
		dest = fmt.Sprintf(",dest=%d", lt.Dest.ID)
	}
	nmstr := ""
	if lt.RailLine != nil {
		nmstr = fmt.Sprintf(":%s", lt.RailLine.Name)
	}
	return fmt.Sprintf("%s(%d):%v,l=%d%s%s%s%s%s%s%s%s", lt.Type().Short(),
		lt.ID, lt.TaskType, lt.RailLineID, before, next, stay, dept, moving, dest, ostr, nmstr)
}

func (ltt LineTaskType) String() string {
	switch ltt {
	case OnDeparture:
		return "dept"
	case OnMoving:
		return "move"
	case OnStopping:
		return "stop"
	case OnPassing:
		return "pass"
	}
	return "????"
}

// eachLineTasks skips LineTask which was added in inner loop
func eachLineTask(lts map[uint]*LineTask, callback func(*LineTask)) {
	copies := make([]*LineTask, len(lts))
	i := 0
	for _, lt := range lts {
		copies[i] = lt
		i++
	}
	for _, lt := range copies {
		callback(lt)
	}
}
