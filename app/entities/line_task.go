package entities

import (
	"fmt"
	"math"
	"time"
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
	Owner

	TaskType LineTaskType `gorm:"not null" json:"type"`

	M         *Model    `gorm:"-" json:"-"`
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
		Base:     NewBase(m.GenID(LINETASK)),
		TaskType: OnDeparture,
	}
	lt.Init(m)
	lt.Resolve(l.Own, l, p)
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
		Base: NewBase(m.GenID(LINETASK)),
	}
	lt.Init(m)
	lt.Resolve(l.Own, l, re)
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

func (lt *LineTask) Depart(force ...bool) *LineTask {
	if !(len(force) > 0 && force[0]) && lt.next != nil {
		panic(fmt.Errorf("Tried to depart from Connectted LineTask: %v", lt))
	}
	if lt.TaskType != OnStopping {
		panic(fmt.Errorf("Tried to depart from invald TaskType : %v", lt))
	}
	return lt.M.NewLineTaskDept(lt.RailLine, lt.Dest, lt)
}

func (lt *LineTask) DepartIf(force ...bool) *LineTask {
	if !(len(force) > 0 && force[0]) && lt.next != nil {
		panic(fmt.Errorf("Tried to depart from Connectted LineTask: %v", lt))
	}
	if lt.TaskType == OnStopping {
		return lt.Depart(force...)
	}
	return lt
}

func (lt *LineTask) Stretch(re *RailEdge, force ...bool) *LineTask {
	if !(len(force) > 0 && force[0]) && lt.next != nil {
		panic(fmt.Errorf("Tried to add RailEdge to Connectted LineTask: %v -> %v", re, lt))
	}
	if lt.ToNode() != re.FromNode {
		panic(fmt.Errorf("Tried to add far RailEdge to LineTask: %v -> %v", re, lt))
	}

	// when task is "stop", append task "departure"
	tail := lt.DepartIf(force...)
	return lt.M.NewLineTask(lt.RailLine, re, tail)
}

// InsertRailEdge corrects RailLine for new RailEdge
// RailEdge.From must be original RailNode.
// RailEdge.To   must be      new RailPoint.
//
// Before (a) ---------------> (b) -> (c)
// After  (a) -> (X) -> (a) -> (b) -> (c)
//
// RailEdge : (a) -> (X)
func (lt *LineTask) InsertRailEdge(re *RailEdge) {
	if lt.ToNode() != re.FromNode {
		panic(fmt.Errorf("Tried to insert far RailEdge to LineTask: %v -> %v", re, lt))
	}
	next := lt.Next()                     // = (b) -> (c)
	tail := lt.Stretch(re, true)          // = (a) -> (X)
	tail = tail.Stretch(re.Reverse, true) // = (X) -> (a)

	// when (X) is station and is stopped, append "dept" task after it
	if next != nil && next.TaskType != OnDeparture {
		tail = tail.DepartIf()
	}
	tail.SetNext(next) // (a) -> (b) -> (c)
}

func (lt *LineTask) InsertDestination(p *Platform) {
	if lt.TaskType == OnDeparture {
		panic(fmt.Errorf("try to insert destination to dept LineTask: %v -> %v", p, lt))
	}
	lt.Dest = p
	lt.DestID = p.ID
	if lt.RailLine.AutoPass {
		// change move -> pass
		lt.TaskType = OnPassing
		lt.RailLine.ReRouting = true
	} else {
		// change move -> stop
		lt.TaskType = OnStopping
		next := lt.next
		lt.Depart(true).SetNext(next)
	}
}

func (lt *LineTask) InsertDeparture(p *Platform) {
	lt.SetDept(p)
}

func (lt *LineTask) Shrink(p *Platform) {
	if lt.Stay != p {
		panic(fmt.Errorf("try to shrink far platform: %v -> %v", p, lt))
	}
	if lt.before != nil {
		lt.before.SetDest(nil)
		lt.before.TaskType = OnMoving
		lt.before.SetNext(lt.next)
		lt.before = nil
	}
	if lt.next != nil {
		lt.next.SetDept(nil)
		lt.next = nil
	}
	lt.Delete()
}

func (lt *LineTask) Shave(re *RailEdge) {
	if lt.Moving != re {
		panic(fmt.Errorf("try to shave far edge: %v -> %v", re, lt))
	}
	if lt.next != nil {
		if lt.next.Moving != re.Reverse {
			panic(fmt.Errorf("try to shave linear RailLine: %v -> %v", re.Reverse, lt.next))
		}
		if lt.before != nil {
			// skip redundant Departure
			if lt.before.TaskType == OnDeparture && lt.next.next != nil && lt.next.next.TaskType == OnDeparture {
				lt.before.SetNext(lt.next.next.next)
				lt.next.next.Delete()
			} else {
				lt.before.SetNext(lt.next.next)
			}
		}
		lt.next.Delete()
	}
	lt.Delete()
}

// Idx returns unique id field.
func (lt *LineTask) Idx() uint {
	return lt.ID
}

// Type returns type of entitiy
func (lt *LineTask) Type() ModelType {
	return LINETASK
}

// Init do nothing
func (lt *LineTask) Init(m *Model) {
	lt.M = m
	lt.Trains = make(map[uint]*Train)
	lt.OverSteps = make(map[uint]*Step)
}

// Pos returns entities' position
func (lt *LineTask) Pos() *Point {
	switch lt.TaskType {
	case OnDeparture:
		if lt.Stay == nil {
			return nil
		}
		return lt.Stay.Pos()
	default:
		if lt.Moving == nil {
			return nil
		}
		return lt.Moving.Pos()
	}
}

// IsIn returns it should be view or not.
func (lt *LineTask) IsIn(x float64, y float64, scale float64) bool {
	switch lt.TaskType {
	case OnDeparture:
		return lt.Stay.IsIn(x, y, scale)
	default:
		return lt.Moving.IsIn(x, y, scale)
	}
}

func (lt *LineTask) Step(prog *float64, sec *float64) {
	canDist := *sec * Const.Train.Speed
	remainDist := (1.0 - *prog) * lt.Cost()
	if remainDist < canDist {
		*sec += remainDist / Const.Train.Speed
		*prog = 1.0
	} else {
		*prog += *sec * Const.Train.Speed / lt.Cost()
		*sec = 0
	}
}

func (lt *LineTask) Loc(prog float64) *Point {
	if prog < 0.5 && lt.before.TaskType == OnDeparture {
		return lt.Moving.Div(2 * prog * prog)
	} else if prog > 0.5 && lt.TaskType == OnDeparture {
		return lt.Moving.Div(-2*prog*prog + 4*prog - 1)
	}
	return lt.Moving.Div(prog)
}

// Resolve set reference
func (lt *LineTask) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			lt.Owner = NewOwner(obj)
			obj.Resolve(lt)
		case *Platform:
			lt.Stay = obj
			lt.Dept = obj
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

func (lt *LineTask) UnResolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Train:
			delete(lt.Trains, obj.ID)
			lt.RailLine.UnResolve(obj)
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
		lt.Stay.UnResolve(lt)
	}
	if lt.Dept != nil {
		lt.Dept.UnResolve(lt)
	}
	if lt.Moving != nil {
		lt.Moving.UnResolve(lt)
	}
	if lt.Dest != nil {
		lt.Dest.UnResolve(lt)
	}
	if lt.before != nil && lt.before.next == lt {
		lt.before.SetNext(nil)
	}
	if lt.next != nil && lt.next.before == lt {
		lt.next.SetBefore(nil)
	}
	lt.RailLine.UnResolve(lt)
	lt.Own.UnResolve(lt)
}

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

// Permits represents Player is permitted to control
func (lt *LineTask) Permits(o *Player) bool {
	return lt.Owner.Permits(o)
}

// From represents start point
func (lt *LineTask) From() Indexable {
	return lt.FromNode()
}

// To represents end point
func (lt *LineTask) To() Indexable {
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

// Cost represents distance
func (lt *LineTask) Cost() float64 {
	switch lt.TaskType {
	case OnDeparture:
		return 0
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
	}
	lt.next = v
	if v != nil {
		lt.NextID = v.ID
		v.SetBefore(lt)
	} else {
		lt.NextID = ZERO
	}
	lt.RailLine.ReRouting = true
	lt.Change()
}

func (lt *LineTask) SetDept(p *Platform) {
	lt.Dept = p
	if p != nil {
		lt.DeptID = p.ID
	} else {
		lt.DeptID = ZERO
	}
}

func (lt *LineTask) SetDest(p *Platform) {
	lt.Dest = p
	if p != nil {
		lt.DestID = p.ID
	} else {
		lt.DestID = ZERO
	}
}

func (lt *LineTask) SetBefore(v *LineTask) {
	lt.before = v
	if v != nil {
		lt.BeforeID = v.ID
	} else {
		lt.BeforeID = ZERO
	}
}

func (lt *LineTask) IsNew() bool {
	return lt.Base.IsNew()
}

// IsChanged returns true when it is changed after Backup()
func (lt *LineTask) IsChanged(after ...time.Time) bool {
	return lt.Base.IsChanged(after...)
}

// Reset set status as not changed
func (lt *LineTask) Reset() {
	lt.Base.Reset()
}

// String represents status
func (lt *LineTask) String() string {
	lt.Marshal()
	ostr := ""
	if lt.Own != nil {
		ostr = fmt.Sprintf(":%s", lt.Own.Short())
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
	posstr := ""
	if lt.Pos() != nil {
		posstr = fmt.Sprintf(":%s", lt.Pos())
	}
	nmstr := ""
	if lt.RailLine != nil {
		nmstr = fmt.Sprintf(":%s", lt.RailLine.Name)
	}
	return fmt.Sprintf("%s(%d):%v,l=%d%s%s%s%s%s%s%s%s%s", lt.Type().Short(),
		lt.ID, lt.TaskType, lt.RailLineID, before, next, stay, dept, moving, dest,
		posstr, ostr, nmstr)
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
