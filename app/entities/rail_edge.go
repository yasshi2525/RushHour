package entities

import (
	"fmt"
)

// RailEdge connects from RailNode to RailNode.
// It's directional.
type RailEdge struct {
	Model
	Owner
	from *RailNode
	to   *RailNode

	Trains map[uint]*Train `gorm:"-" json:"-"`

	FromID uint `gorm:"not null" json:"from"`
	ToID   uint `gorm:"not null" json:"to"`
}

// NewRailEdge create new instance and relates RailNode
func NewRailEdge(id uint, f *RailNode, t *RailNode) *RailEdge {
	re := &RailEdge{
		Model:  NewModel(id),
		Owner:  f.Owner,
		from:   f,
		to:     t,
		Trains: make(map[uint]*Train),
	}
	re.ResolveRef()

	f.OutEdge[re.ID] = re
	t.InEdge[re.ID] = re
	return re
}

// Idx returns unique id field.
func (re *RailEdge) Idx() uint {
	return re.ID
}

// Init do nothing
func (re *RailEdge) Init() {
	re.Model.Init()
	re.Owner.Init()
	re.Trains = make(map[uint]*Train)
}

// Pos returns location
func (re *RailEdge) Pos() *Point {
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
			} else {
				re.to = obj
			}
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
	re.Owner.ResolveRef()
	if re.from != nil {
		re.FromID = re.from.ID
	}
	if re.to != nil {
		re.ToID = re.to.ID
	}
}

// Permits represents Player is permitted to control
func (re *RailEdge) Permits(o *Player) bool {
	return re.Owner.Permits(o)
}

// String represents status
func (re *RailEdge) String() string {
	return fmt.Sprintf("%s(%d):f=%d,t=%d:%v", Meta.Static[RAILEDGE],
		re.ID, re.from.ID, re.to.ID, re.Pos())
}
