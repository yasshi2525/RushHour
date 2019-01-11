package entities

import (
	"fmt"
	"time"
)

// RailEdge connects from RailNode to RailNode.
// It's directional.
type RailEdge struct {
	Base
	Owner

	FromNode  *RailNode
	ToNode    *RailNode
	Reverse   *RailEdge          `gorm:"-" json:"-"`
	RailLines map[uint]*RailLine `gorm:"-" json:"-"`
	LineTasks map[uint]*LineTask `gorm:"-" json:"-"`
	Trains    map[uint]*Train    `gorm:"-" json:"-"`

	FromID    uint `gorm:"not null" json:"from"`
	ToID      uint `gorm:"not null" json:"to"`
	ReverseID uint `gorm:"not null" json:"eid"`
}

// NewRailEdge create new instance and relates RailNode
func (m *Model) NewRailEdge(f *RailNode, t *RailNode) *RailEdge {
	re := &RailEdge{
		Base:     NewBase(m.GenID(RAILEDGE)),
		Owner:    f.Owner,
		FromNode: f,
		ToNode:   t,
	}
	re.Init()
	re.ResolveRef()

	f.Own.Resolve(re)
	f.OutEdge[re.ID] = re
	t.InEdge[re.ID] = re

	m.Add(re)
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
func (re *RailEdge) Init() {
	re.RailLines = make(map[uint]*RailLine)
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
	return re.FromNode.Pos().Dist(re.ToNode)
}

// Unrelate delete relations to RailNode
func (re *RailEdge) UnRef() {
	delete(re.FromNode.OutEdge, re.ID)
	delete(re.ToNode.InEdge, re.ID)
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
				obj.OutEdge[re.ID] = re
			} else {
				re.ToNode = obj
				obj.InEdge[re.ID] = re
			}
		case *RailEdge:
			re.Reverse = obj
			obj.Reverse = re
		case *RailLine:
			re.RailLines[obj.ID] = obj
		case *LineTask:
			re.LineTasks[obj.ID] = obj
		case *Train:
			re.Trains[obj.ID] = obj
			obj.Resolve(re)
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	re.ResolveRef()
}

// ResolveRef set id from reference
func (re *RailEdge) ResolveRef() {
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

// CheckRemove check remain relation.
func (re *RailEdge) CheckRemove() error {
	for _, lt := range re.LineTasks {
		if err := lt.CheckRemove(); err != nil {
			return fmt.Errorf("blocked by LineTask of %v; %v", lt, err)
		}
	}
	if len(re.Trains) > 0 {
		return fmt.Errorf("blocked by Train of %v", re.Trains)
	}
	return nil
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
	re.ResolveRef()
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
