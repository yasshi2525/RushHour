package entities

import (
	"fmt"
)

// RailEdge connects from RailNode to RailNode.
// It's directional.
type RailEdge struct {
	Base
	Persistence
	Shape

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
		Base:        m.NewBase(RAILEDGE, f.O),
		Persistence: NewPersistence(),
		Shape:       NewShapeEdge(&f.Point, &t.Point),
	}
	re.Init(m)
	re.Resolve(f.O, f, t)
	re.Marshal()
	m.Add(re)
	re.O.ReRouting = true
	return re
}

// B returns base information of this elements.
func (re *RailEdge) B() *Base {
	return &re.Base
}

// P returns time information for database.
func (re *RailEdge) P() *Persistence {
	return &re.Persistence
}

// S returns entities' position.
func (re *RailEdge) S() *Shape {
	return &re.Shape
}

// Init do nothing
func (re *RailEdge) Init(m *Model) {
	re.Base.Init(RAILEDGE, m)
	re.LineTasks = make(map[uint]*LineTask)
	re.Trains = make(map[uint]*Train)
}

// From represents start point
func (re *RailEdge) From() Entity {
	return re.FromNode
}

// To represents end point
func (re *RailEdge) To() Entity {
	return re.ToNode
}

// Cost represents distance
func (re *RailEdge) Cost() float64 {
	return re.Shape.Dist() / Const.Train.Speed
}

// CheckDelete check remain relation.
func (re *RailEdge) CheckDelete() error {
	for _, obj := range []*RailEdge{re, re.Reverse} {
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
	delete(re.FromNode.OutEdges, re.ID)
	delete(re.ToNode.InEdges, re.ID)
	re.FromNode.Shape.UnRefer(re.S())
	re.ToNode.Shape.UnRefer(re.S())
	re.O.UnResolve(re)
}

// Delete removes this entity with related ones.
func (re *RailEdge) Delete() {
	eachLineTask(re.Reverse.LineTasks, func(lt *LineTask) {
		lt.Shave(re.Reverse)
	})
	eachLineTask(re.LineTasks, func(lt *LineTask) {
		lt.Shave(re)
	})
	re.O.ReRouting = true
	re.M.Delete(re.Reverse)
	re.M.Delete(re)
}

// Resolve set reference
func (re *RailEdge) Resolve(args ...Entity) {
	var doneFrom bool
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			re.O = obj
			obj.Resolve(re)
		case *RailNode:
			if !doneFrom {
				re.O, re.FromNode, re.Shape.P1 = obj.O, obj, &obj.Point
				doneFrom = true
				obj.OutEdges[re.ID] = re
				obj.Shape.Refer(re.S())
			} else {
				re.ToNode, re.Shape.P2 = obj, &obj.Point
				obj.InEdges[re.ID] = re
				obj.Shape.Refer(re.S())
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

// UnResolve unregisters specified refernce.
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
	if re.O != nil {
		re.OwnerID = re.O.ID
	}
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

// UnMarshal set reference from id.
func (re *RailEdge) UnMarshal() {
	re.Resolve(
		re.M.Find(PLAYER, re.OwnerID),
		re.M.Find(RAILNODE, re.FromID),
		re.M.Find(RAILNODE, re.ToID),
		re.M.Find(RAILEDGE, re.ReverseID))
}

// String represents status
func (re *RailEdge) String() string {
	re.Marshal()
	ostr := ""
	if re.O != nil {
		ostr = fmt.Sprintf(":%s", re.O.Short())
	}
	posstr := ""
	if re.Pos() != nil {
		posstr = fmt.Sprintf(":%s", re.Pos())
	}
	return fmt.Sprintf("%s(%d):f=%d,t=%d,r=%d%s%s", re.Type().Short(),
		re.ID, re.FromID, re.ToID, re.ReverseID, posstr, ostr)
}
