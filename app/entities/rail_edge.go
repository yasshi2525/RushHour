package entities

import (
	"fmt"
	"math"
	"time"
)

// RailEdge connects from RailNode to RailNode.
// It's directional.
type RailEdge struct {
	Base
	Owner

	M         *Model             `gorm:"-" json:"-"`
	FromNode  *RailNode          `gorm:"-" json:"-"`
	ToNode    *RailNode          `gorm:"-" json:"-"`
	Reverse   *RailEdge          `gorm:"-" json:"-"`
	LineTasks map[uint]*LineTask `gorm:"-" json:"-"`
	Trains    map[uint]*Train    `gorm:"-" json:"-"`

	FromID    uint `gorm:"not null" json:"from"`
	ToID      uint `gorm:"not null" json:"to"`
	ReverseID uint `gorm:"not null" json:"eid"`
}

// NewRailEdge create new instance and relates RailNode
func (m *Model) NewRailEdge(f *RailNode, t *RailNode) *RailEdge {
	re := &RailEdge{
		Base: NewBase(m.GenID(RAILEDGE)),
	}
	re.Init(m)
	re.Resolve(f.Own, f, t)
	re.Marshal()
	m.Add(re)
	re.Own.ReRouting = true
	return re
}

// Idx returns unique id field.
func (re *RailEdge) Idx() uint {
	return re.ID
}

// Type returns type of entitiy
func (re *RailEdge) Type() ModelType {
	return RAILEDGE
}

// Init do nothing
func (re *RailEdge) Init(m *Model) {
	re.M = m
	re.LineTasks = make(map[uint]*LineTask)
	re.Trains = make(map[uint]*Train)
}

// Pos returns location
func (re *RailEdge) Pos() *Point {
	if re.FromNode == nil || re.ToNode == nil {
		return nil
	}
	return re.FromNode.Pos().Center(re.ToNode)
}

// IsIn return true when from, to, center is in,
func (re *RailEdge) IsIn(x float64, y float64, scale float64) bool {
	return re.FromNode.Pos().IsInLine(re.ToNode, x, y, scale)
}

// From represents start point
func (re *RailEdge) From() Indexable {
	return re.FromNode
}

// To represents end point
func (re *RailEdge) To() Indexable {
	return re.ToNode
}

// Cost represents distance
func (re *RailEdge) Cost() float64 {
	return re.FromNode.Pos().Dist(re.ToNode) / Const.Train.Speed
}

func (re *RailEdge) Angle(to *RailEdge) float64 {
	v := &Point{re.ToNode.X - re.FromNode.X, re.ToNode.Y - re.FromNode.Y}
	lenV := math.Sqrt(v.X*v.X + v.Y*v.Y)
	u := &Point{to.ToNode.X - to.FromNode.X, to.ToNode.Y - to.FromNode.Y}
	lenU := math.Sqrt(u.X*u.X + u.Y*u.Y)
	inner := (v.X*u.X + v.Y*u.Y) / (lenV * lenU)
	theta := math.Acos(inner)
	if math.IsNaN(theta) {
		return math.Pi
	}
	return theta
}

func (re *RailEdge) Div(prog float64) *Point {
	return re.FromNode.Pos().Div(re.ToNode, prog)
}

// CheckDelete check remain relation.
func (re *RailEdge) CheckDelete() error {
	for _, obj := range []*RailEdge{re, re.Reverse} {
		if len(obj.Trains) > 0 {
			return fmt.Errorf("blocked by Train of %v", re.Trains)
		}
		for _, lt := range obj.LineTasks {
			// if RailLine is not sharp, forbit remove
			if lt.next != nil && lt.next.Moving != re.Reverse {
				return fmt.Errorf("blocked by LineTask of %v", lt)
			}
		}
	}
	return nil
}

// BeforeDelete delete relations to RailNode
func (re *RailEdge) BeforeDelete() {
	for _, t := range re.Trains {
		t.SetTask(t.task.next)
	}
	delete(re.FromNode.OutEdges, re.ID)
	delete(re.ToNode.InEdges, re.ID)
	re.Own.UnResolve(re)
}

func (re *RailEdge) Delete() {
	eachLineTask(re.Reverse.LineTasks, func(lt *LineTask) {
		lt.Shave(re.Reverse)
	})
	eachLineTask(re.LineTasks, func(lt *LineTask) {
		lt.Shave(re)
	})
	re.Own.ReRouting = true
	re.M.Delete(re.Reverse)
	re.M.Delete(re)
}

// Resolve set reference
func (re *RailEdge) Resolve(args ...interface{}) {
	var doneFrom bool
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			re.Owner = NewOwner(obj)
			obj.Resolve(re)
		case *RailNode:
			if !doneFrom {
				re.Owner, re.FromNode = obj.Owner, obj
				doneFrom = true
				obj.OutEdges[re.ID] = re
			} else {
				re.ToNode = obj
				obj.InEdges[re.ID] = re
			}
		case *RailEdge:
			re.Reverse = obj
			obj.Reverse = re
		case *LineTask:
			re.LineTasks[obj.ID] = obj
		case *Train:
			re.Trains[obj.ID] = obj
			obj.Resolve(re)
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	re.Marshal()
}

func (re *RailEdge) UnResolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *LineTask:
			delete(re.LineTasks, obj.ID)
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
}

// Marshal set id from reference
func (re *RailEdge) Marshal() {
	if re.FromNode != nil {
		re.FromID = re.FromNode.ID
	}
	if re.ToNode != nil {
		re.ToID = re.ToNode.ID
	}
	if re.Reverse != nil {
		re.ReverseID = re.Reverse.ID
	}
}

func (re *RailEdge) UnMarshal() {
	re.Resolve(
		re.M.Find(PLAYER, re.OwnerID),
		re.M.Find(RAILNODE, re.FromID),
		re.M.Find(RAILNODE, re.ToID),
		re.M.Find(RAILEDGE, re.ReverseID))
}

// Permits represents Player is permitted to control
func (re *RailEdge) Permits(o *Player) bool {
	return re.Owner.Permits(o)
}

func (re *RailEdge) IsNew() bool {
	return re.Base.IsNew()
}

// IsChanged returns true when it is changed after Backup()
func (re *RailEdge) IsChanged(after ...time.Time) bool {
	return re.Base.IsChanged(after...)
}

// Reset set status as not changed
func (re *RailEdge) Reset() {
	re.Base.Reset()
}

// String represents status
func (re *RailEdge) String() string {
	re.Marshal()
	ostr := ""
	if re.Own != nil {
		ostr = fmt.Sprintf(":%s", re.Own.Short())
	}
	posstr := ""
	if re.Pos() != nil {
		posstr = fmt.Sprintf(":%s", re.Pos())
	}
	return fmt.Sprintf("%s(%d):f=%d,t=%d,r=%d%s%s", re.Type().Short(),
		re.ID, re.FromID, re.ToID, re.ReverseID, posstr, ostr)
}
