package entities

import (
	"fmt"
)

// Station composes on Platform and Gate
type Station struct {
	Base
	Persistence
	Shape

	Name string `gorm:"not null" json:"name"`

	Platform *Platform `gorm:"-" json:"-"`
	Gate     *Gate     `gorm:"-" json:"-"`

	PlatformID uint `gorm:"-" json:"pid"`
	GateID     uint `gorm:"-" json:"gid"`
}

// NewStation create new instance.
func (m *Model) NewStation(o *Player) *Station {
	st := &Station{
		Base:        m.NewBase(STATION, o),
		Persistence: NewPersistence(),
	}
	st.Init(m)
	st.Resolve(o)
	st.Marshal()
	m.Add(st)
	return st
}

// B returns base information of this elements.
func (st *Station) B() *Base {
	return &st.Base
}

// P returns time information for database.
func (st *Station) P() *Persistence {
	return &st.Persistence
}

// S returns entities' position.
func (st *Station) S() *Shape {
	return &st.Shape
}

// Init creates map.
func (st *Station) Init(m *Model) {
	st.Base.Init(STATION, m)
}

// Resolve set reference from id.
func (st *Station) Resolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			st.O = obj
			obj.Resolve(st)
		case *Gate:
			st.Gate = obj
		case *Platform:
			st.Platform = obj
			st.Shape = obj.Shape
			st.M.RootCluster.Add(st)
			obj.Resolve(st.Gate)
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	st.Marshal()
}

// Marshal resolve Owner reference
func (st *Station) Marshal() {
	if st.O != nil {
		st.OwnerID = st.O.ID
	}
	if st.Platform != nil {
		st.PlatformID = st.Platform.ID
	}
	if st.Gate != nil {
		st.GateID = st.Gate.ID
	}
}

// UnMarshal set reference from id.
func (st *Station) UnMarshal() {
	st.Resolve(st.M.Find(PLAYER, st.OwnerID))
}

// CheckDelete checks related reference
func (st *Station) CheckDelete() error {
	if err := st.Platform.CheckDelete(); err != nil {
		return err
	}
	if err := st.Gate.CheckDelete(); err != nil {
		return err
	}
	return nil
}

// BeforeDelete delete related reference
func (st *Station) BeforeDelete() {
	st.O.UnResolve(st)
}

// Delete removes this entity with related ones.
func (st *Station) Delete() {
	if st.Gate != nil {
		st.M.Delete(st.Gate)
	}
	if st.Platform != nil {
		st.M.Delete(st.Platform)
	}
	st.M.Delete(st)
}

// String represents status
func (st *Station) String() string {
	st.Marshal()
	ostr := ""
	if st.O != nil {
		ostr = fmt.Sprintf(":%s", st.O.Short())
	}
	posstr := ""
	if st.Pos() != nil {
		posstr = fmt.Sprintf(":%s", st.Pos())
	}
	return fmt.Sprintf("%s(%d):g=%d,p=%d%s%s:%s", st.Type().Short(),
		st.ID, st.PlatformID, st.GateID, posstr, ostr, st.Name)
}
