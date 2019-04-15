package entities

import (
	"fmt"
)

// Platform is the base Human wait for Train.
// Platform can enter only through Gate.
type Platform struct {
	Base
	Persistence
	Shape

	Capacity int `gorm:"not null" json:"cap"`
	Occupied int `gorm:"-"        json:"used"`

	OnRailNode *RailNode          `gorm:"-" json:"-"`
	InStation  *Station           `gorm:"-" json:"-"`
	WithGate   *Gate              `gorm:"-" json:"-"`
	InTasks    map[uint]*LineTask `gorm:"-" json:"-"`
	StayTasks  map[uint]*LineTask `gorm:"-" json:"-"`
	OutTasks   map[uint]*LineTask `gorm:"-" json:"-"`
	Trains     map[uint]*Train    `gorm:"-" json:"-"`
	Passengers map[uint]*Human    `gorm:"-" json:"-"`
	// key is id of Platform
	Transports map[uint]*Transport `gorm:"-" json:"-"`
	outSteps   map[uint]*Step
	inSteps    map[uint]*Step

	StationID  uint `json:"stid"`
	RailNodeID uint `json:"rnid"`
	GateID     uint `gorm:"-" json:"gid"`
}

// NewPlatform creates instance
func (m *Model) NewPlatform(rn *RailNode, g *Gate) *Platform {
	p := &Platform{
		Base:        m.NewBase(PLATFORM, rn.O),
		Persistence: NewPersistence(),
		Shape:       rn.Shape,
		Capacity:    Const.Platform.Capacity,
	}
	p.Init(m)
	p.Resolve(rn.O, rn, g.InStation, g)
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

// B returns base information of this elements.
func (p *Platform) B() *Base {
	return &p.Base
}

// P returns time information for database.
func (p *Platform) P() *Persistence {
	return &p.Persistence
}

// S returns entities' position.
func (p *Platform) S() *Shape {
	return &p.Shape
}

// GenOutSteps generates Steps from this Platform.
func (p *Platform) GenOutSteps() {
	// P -> G
	p.M.NewStep(p, p.WithGate)
}

// GenInSteps generates Steps to this Platform.
func (p *Platform) GenInSteps() {
	// G -> P
	p.M.NewStep(p.WithGate, p)
}

// Init creates map.
func (p *Platform) Init(m *Model) {
	p.Base.Init(PLATFORM, m)
	p.InTasks = make(map[uint]*LineTask)
	p.StayTasks = make(map[uint]*LineTask)
	p.OutTasks = make(map[uint]*LineTask)
	p.Trains = make(map[uint]*Train)
	p.Passengers = make(map[uint]*Human)
	p.Transports = make(map[uint]*Transport)
	p.outSteps = make(map[uint]*Step)
	p.inSteps = make(map[uint]*Step)
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
func (p *Platform) Resolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			p.O = obj
			obj.Resolve(p)
		case *RailNode:
			p.OnRailNode = obj
			p.Shape = obj.Shape
			obj.Resolve(p)
		case *Station:
			p.InStation = obj
			obj.Resolve(p)
		case *Gate:
			p.WithGate = obj
			obj.Resolve(p)
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

// UnResolve unregisters specified refernce.
func (p *Platform) UnResolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *LineTask:
			switch obj.TaskType {
			case OnDeparture:
				delete(p.StayTasks, obj.ID)
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
	if p.O != nil {
		p.OwnerID = p.O.ID
	}
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

// UnMarshal set reference from id.
func (p *Platform) UnMarshal() {
	st := p.M.Find(STATION, p.StationID).(*Station)
	p.Resolve(
		p.M.Find(PLAYER, p.OwnerID),
		p.M.Find(RAILNODE, p.RailNodeID),
		st, st.Gate)
}

// CheckDelete checks related reference
func (p *Platform) CheckDelete() error {
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
	eachLineTask(p.InTasks, func(lt *LineTask) {
		if lt.TaskType == OnDeparture {
			panic(fmt.Errorf("lt should not be OnDeparture: %v", lt))
		}
		lt.TaskType = OnMoving
		lt.SetDest(nil)
	})
	p.OnRailNode.UnResolve(p)
	p.O.UnResolve(p)
}

// Delete removes this entity with related ones.
func (p *Platform) Delete() {
	for _, s := range p.outSteps {
		p.M.Delete(s)
	}
	for _, s := range p.inSteps {
		p.M.Delete(s)
	}
	p.M.Delete(p)
}

// String represents status
func (p *Platform) String() string {
	p.Marshal()
	ostr := ""
	if p.O != nil {
		ostr = fmt.Sprintf(":%s", p.O.Short())
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
