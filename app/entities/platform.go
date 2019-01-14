package entities

import (
	"fmt"
	"time"
)

// Platform is the base Human wait for Train.
// Platform can enter only through Gate.
type Platform struct {
	Base
	Owner

	M          *Model             `gorm:"-" json:"-"`
	OnRailNode *RailNode          `gorm:"-" json:"-"`
	InStation  *Station           `gorm:"-" json:"-"`
	WithGate   *Gate              `gorm:"-" json:"-"`
	RailLines  map[uint]*RailLine `gorm:"-" json:"-"`
	InTasks    map[uint]*LineTask `gorm:"-" json:"-"`
	StayTasks  map[uint]*LineTask `gorm:"-" json:"-"`
	OutTasks   map[uint]*LineTask `gorm:"-" json:"-"`
	Trains     map[uint]*Train    `gorm:"-" json:"-"`
	Passengers map[uint]*Human    `gorm:"-" json:"-"`
	// key is id of Platform
	Transports map[uint]*Transport `gorm:"-" json:"-"`
	outSteps   map[uint]*Step
	inSteps    map[uint]*Step

	Capacity int `gorm:"not null" json:"cap"`
	Occupied int `gorm:"-"        json:"used"`

	StationID  uint `gorm:"not null" json:"stid"`
	GateID     uint `gorm:"-"        json:"gid"`
	RailNodeID uint `gorm:"not null" json:"rnid"`
}

// NewPlatform creates instance
func (m *Model) NewPlatform(rn *RailNode, g *Gate) *Platform {
	p := &Platform{
		Base:     NewBase(m.GenID(PLATFORM)),
		Capacity: Const.Platform.Capacity,
	}
	p.Init(m)
	p.Resolve(rn.Own, rn, g.InStation, g)
	p.Marshal()
	m.Add(p)

	// find LineTask such as dest to new platform point
	eachLineTask(rn.InTasks, func(lt *LineTask) {
		lt.InsertDestination(p)
	})
	// find LineTask such as dept from new platform point
	eachLineTask(rn.OutTasks, func(lt *LineTask) {
		lt.InsertDeparture(p)
	})

	p.GenOutSteps()
	p.GenInSteps()
	return p
}

func (p *Platform) GenOutSteps() {
	// P -> G
	p.M.NewStep(p, p.WithGate)
}

func (p *Platform) GenInSteps() {
	// G -> P
	p.M.NewStep(p.WithGate, p)
}

// Idx returns unique id field.
func (p *Platform) Idx() uint {
	return p.ID
}

// Type returns type of entitiy
func (p *Platform) Type() ModelType {
	return PLATFORM
}

// Init creates map.
func (p *Platform) Init(m *Model) {
	p.M = m
	p.RailLines = make(map[uint]*RailLine)
	p.InTasks = make(map[uint]*LineTask)
	p.StayTasks = make(map[uint]*LineTask)
	p.OutTasks = make(map[uint]*LineTask)
	p.Trains = make(map[uint]*Train)
	p.Passengers = make(map[uint]*Human)
	p.Transports = make(map[uint]*Transport)
	p.outSteps = make(map[uint]*Step)
	p.inSteps = make(map[uint]*Step)
}

// Pos returns location
func (p *Platform) Pos() *Point {
	if p.OnRailNode == nil {
		return nil
	}
	return p.OnRailNode.Pos()
}

// IsIn returns it should be view or not.
func (p *Platform) IsIn(x float64, y float64, scale float64) bool {
	return p.Pos().IsIn(x, y, scale)
}

// OutSteps returns where it can go to
func (p *Platform) OutSteps() map[uint]*Step {
	return p.outSteps
}

// InSteps returns where it comes from
func (p *Platform) InSteps() map[uint]*Step {
	return p.inSteps
}

// Resolve set reference
func (p *Platform) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			p.Owner = NewOwner(obj)
			obj.Resolve(p)
		case *RailNode:
			p.OnRailNode = obj
			obj.Resolve(p)
		case *Station:
			p.InStation = obj
			obj.Resolve(p)
		case *Gate:
			p.WithGate = obj
			obj.Resolve(p)
		case *RailLine:
			p.RailLines[obj.ID] = obj
		case *LineTask:
			switch obj.TaskType {
			case OnDeparture:
				p.StayTasks[obj.ID] = obj
			default:
				if obj.Dept == p {
					p.OutTasks[obj.ID] = obj
				} else if obj.Dest == p {
					p.InTasks[obj.ID] = obj
				} else {
					panic(fmt.Errorf("unrelated LineTask %v -> %v", obj, p))
				}
			}
		case *Train:
			p.Trains[obj.ID] = obj
		case *Human:
			p.Passengers[obj.ID] = obj
			p.Occupied++
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	p.Marshal()
}

func (p *Platform) UnResolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *LineTask:
			switch obj.TaskType {
			case OnDeparture:
				delete(p.StayTasks, obj.ID)
				delete(p.OnRailNode.InTasks, obj.ID)
				delete(p.OnRailNode.OutTasks, obj.ID)
			default:
				if obj.Dept == p {
					delete(p.OutTasks, obj.ID)
				} else if obj.Dest == p {
					delete(p.InTasks, obj.ID)
				} else {
					panic(fmt.Errorf("unrelated LineTask %v -> %v", obj, p))
				}
			}
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
}

// Marshal set id from reference
func (p *Platform) Marshal() {
	if p.OnRailNode != nil {
		p.RailNodeID = p.OnRailNode.ID
	}
	if p.WithGate != nil {
		p.GateID = p.WithGate.ID
	}
	if p.InStation != nil {
		p.StationID = p.InStation.ID
	}
}

func (p *Platform) UnMarshal() {
	st := p.M.Find(STATION, p.StationID).(*Station)
	p.Resolve(
		p.M.Find(PLAYER, p.OwnerID),
		p.M.Find(RAILNODE, p.RailNodeID),
		st, st.Gate)
}

// CheckDelete checks related reference
func (p *Platform) CheckDelete() error {
	if len(p.Trains) > 0 {
		return fmt.Errorf("blocked by Train of %v", p.Trains)
	}
	return nil
}

// BeforeDelete delete related reference.
func (p *Platform) BeforeDelete() {
	for _, h := range p.Passengers {
		h.UnResolve(p)
	}
	for _, t := range p.Trains {
		t.UnResolve(p)
	}
	eachLineTask(p.StayTasks, func(lt *LineTask) {
		lt.Shrink(p)
	})
	p.OnRailNode.UnResolve(p)
	p.Own.UnResolve(p)
}

func (p *Platform) Delete() {
	for _, s := range p.outSteps {
		p.M.Delete(s)
	}
	for _, s := range p.inSteps {
		p.M.Delete(s)
	}
	p.M.Delete(p)
}

// Permits represents Player is permitted to control
func (p *Platform) Permits(o *Player) bool {
	return p.Owner.Permits(o)
}

func (p *Platform) IsNew() bool {
	return p.Base.IsNew()
}

// IsChanged returns true when it is changed after Backup()
func (p *Platform) IsChanged(after ...time.Time) bool {
	return p.Base.IsChanged(after...)
}

// Reset set status as not changed
func (p *Platform) Reset() {
	p.Base.Reset()
}

// String represents status
func (p *Platform) String() string {
	p.Marshal()
	ostr := ""
	if p.Own != nil {
		ostr = fmt.Sprintf(":%s", p.Own.Short())
	}
	posstr := ""
	if p.Pos() != nil {
		posstr = fmt.Sprintf(":%s", p.Pos())
	}
	nmstr := ""
	if p.InStation != nil {
		nmstr = fmt.Sprintf(":%s", p.InStation.Name)
	}
	return fmt.Sprintf("%s(%d):st=%d,g=%d,rn=%d,i=%d,o=%d,h=%d/%d%s%s%s",
		p.Type().Short(),
		p.ID, p.StationID, p.GateID, p.RailNodeID,
		len(p.inSteps), len(p.outSteps), len(p.Passengers), p.Capacity,
		posstr, ostr, nmstr)
}
