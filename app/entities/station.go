package entities

import (
	"fmt"
)

// Station composes on Platform and Gate
type Station struct {
	Base
	Owner

	Platform *Platform `gorm:"-" json:"-"`
	Gate     *Gate     `gorm:"-" json:"-"`

	PlatformID uint `gorm:"-" json:"pid"`
	GateID     uint `gorm:"-" json:"gid"`

	Name string `json:"name"`
}

// NewStation create new instance.
func NewStation(stid uint, gid uint, pid uint, rn *RailNode) (*Station, *Gate, *Platform) {
	p := &Platform{
		Base:       NewBase(pid),
		Owner:      rn.Owner,
		in:         make(map[uint]*Step),
		out:        make(map[uint]*Step),
		Trains:     make(map[uint]*Train),
		Passengers: make(map[uint]*Human),
		LineTasks:  make(map[uint]*LineTask),
		OnRailNode: rn,
	}
	rn.OverPlatform = p

	g := &Gate{
		Base:  NewBase(gid),
		Owner: rn.Owner,
		in:    make(map[uint]*Step),
		out:   make(map[uint]*Step),
	}

	st := &Station{
		Base:     NewBase(stid),
		Owner:    rn.Owner,
		Platform: p,
		Gate:     g,
	}

	p.InStation = st
	g.InStation = st

	p.ResolveRef()
	g.ResolveRef()
	st.ResolveRef()

	return st, g, p
}

// Idx returns unique id field.
func (st *Station) Idx() uint {
	return st.ID
}

// Type returns type of entitiy
func (st *Station) Type() ModelType {
	return STATION
}

// Init creates map.
func (st *Station) Init() {
}

// Pos returns location
func (st *Station) Pos() *Point {
	return st.Platform.Pos()
}

// IsIn returns it should be view or not.
func (st *Station) IsIn(center *Point, scale float64) bool {
	return st.Pos().IsIn(center, scale)
}

// Resolve set reference from id.
func (st *Station) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			st.Own = obj
		case *Gate:
			st.Gate = obj
		case *Platform:
			st.Platform = obj
			st.Gate.Resolve(obj)
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	st.ResolveRef()
}

// ResolveRef resolve Owner reference
func (st *Station) ResolveRef() {
	st.Owner.ResolveRef()
	if st.Platform != nil {
		st.PlatformID = st.Platform.ID
	}
	if st.Gate != nil {
		st.GateID = st.Gate.ID
	}
}

// CheckRemove checks related reference
func (st *Station) CheckRemove() error {
	if err := st.Gate.CheckRemove(); err != nil {
		return err
	}
	if err := st.Platform.CheckRemove(); err != nil {
		return err
	}
	return nil
}

// UnRef delete related reference
func (st *Station) UnRef() {
	st.Platform.UnRef()
	st.Gate.UnRef()
}

// Permits represents Player is permitted to control
func (st *Station) Permits(o *Player) bool {
	return st.Owner.Permits(o)
}

// String represents status
func (st *Station) String() string {
	return fmt.Sprintf("%s(%d):g=%d,p=%d:%v:%s", Meta.Attr[st.Type()].Short,
		st.ID, st.Platform.ID, st.Gate.ID, st.Pos(), st.Name)
}
