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
	from *RailNode
	to   *RailNode

	Trains    map[uint]*Train    `gorm:"-" json:"-"`
	LineTasks map[uint]*LineTask `gorm:"-" json:"-"`

	FromID uint `gorm:"not null" json:"from"`
	ToID   uint `gorm:"not null" json:"to"`
}

// NewRailEdge create new instance and relates RailNode
func NewRailEdge(id uint, f *RailNode, t *RailNode) *RailEdge {
	re := &RailEdge{
		Base:  NewBase(id),
		Owner: f.Owner,
		from:  f,
		to:    t,
	}
	re.Init()
	re.ResolveRef()

	f.OutEdge[re.ID] = re
	t.InEdge[re.ID] = re

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
	re.Trains = make(map[uint]*Train)
	re.LineTasks = make(map[uint]*LineTask)
}

// Pos returns location
func (re *RailEdge) Pos() *Point {
	if re.from == nil || re.to == nil {
		return nil
	}
	return re.from.Pos().Center(re.to.Pos())
}

// IsIn return true when from, to, center is in,
func (re *RailEdge) IsIn(center *Point, scale float64) bool {
	return re.from.Pos().IsInLine(re.to.Pos(), center, scale)
}

// From represents start point
func (re *RailEdge) From() Locationable {
	return re.from
}

// To represents end point
func (re *RailEdge) To() Locationable {
	return re.to
}

// Cost represents distance
func (re *RailEdge) Cost() float64 {
	return re.from.Pos().Dist(re.to.Pos())
}

// Unrelate delete relations to RailNode
func (re *RailEdge) Unrelate() {
	delete(re.from.OutEdge, re.ID)
	delete(re.to.InEdge, re.ID)
}

// Resolve set reference
func (re *RailEdge) Resolve(args ...interface{}) {
	var doneFrom bool
	for _, raw := range args {
		switch obj := raw.(type) {
		case *RailNode:
			if !doneFrom {
				re.Owner, re.from = obj.Owner, obj
				doneFrom = true
				obj.OutEdge[re.ID] = re
			} else {
				re.to = obj
				obj.InEdge[re.ID] = re
			}
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
	if re.from != nil {
		re.FromID = re.from.ID
	}
	if re.to != nil {
		re.ToID = re.to.ID
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

// IsChanged returns true when it is changed after Backup()
func (re *RailEdge) IsChanged(after ...time.Time) bool {
	return re.Base.IsChanged(after)
}

// Reset set status as not changed
func (re *RailEdge) Reset() {
	re.Base.Reset()
}

// String represents status
func (re *RailEdge) String() string {
	ostr := ""
	if re.Own != nil {
		ostr = fmt.Sprintf(":%s", re.Own.Short())
	}
	posstr := ""
	if re.Pos() != nil {
		posstr = fmt.Sprintf(":%s", re.Pos())
	}
	return fmt.Sprintf("%s(%d):f=%d,t=%d%s%s", Meta.Attr[re.Type()].Short,
		re.ID, re.FromID, re.ToID, posstr, ostr)
}
