package entities

import (
	"fmt"
	"time"
)

// Station composes on Platform and Gate
type Station struct {
	Base
	Owner

	M        *Model    `gorm:"-" json:"-"`
	Platform *Platform `gorm:"-" json:"-"`
	Gate     *Gate     `gorm:"-" json:"-"`

	PlatformID uint `gorm:"-" json:"pid"`
	GateID     uint `gorm:"-" json:"gid"`

	Name string `json:"name"`
}

// NewStation create new instance.
func (m *Model) NewStation(o *Player) *Station {
	st := &Station{
		Base: NewBase(m.GenID(STATION)),
	}
	st.Init(m)
	st.Resolve(o)
	st.Marshal()
	m.Add(st)
	return st
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
func (st *Station) Init(m *Model) {
	st.M = m
}

// Pos returns location
func (st *Station) Pos() *Point {
	if st.Platform == nil {
		return nil
	}
	return st.Platform.Pos()
}

// IsIn returns it should be view or not.
func (st *Station) IsIn(x float64, y float64, scale float64) bool {
	return st.Pos().IsIn(x, y, scale)
}

// Resolve set reference from id.
func (st *Station) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			st.Owner = NewOwner(obj)
			obj.Resolve(st)
		case *Gate:
			st.Gate = obj
		case *Platform:
			st.Platform = obj
			obj.Resolve(st.Gate)
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	st.Marshal()
}

// Marshal resolve Owner reference
func (st *Station) Marshal() {
	if st.Platform != nil {
		st.PlatformID = st.Platform.ID
	}
	if st.Gate != nil {
		st.GateID = st.Gate.ID
	}
}

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
}

func (st *Station) Delete() {
	st.M.Delete(st.Gate)
	st.M.Delete(st.Platform)
	st.M.Delete(st)
}

// Permits represents Player is permitted to control
func (st *Station) Permits(o *Player) bool {
	return st.Owner.Permits(o)
}

func (st *Station) IsNew() bool {
	return st.Base.IsNew()
}

// IsChanged returns true when it is changed after Backup()
func (st *Station) IsChanged(after ...time.Time) bool {
	return st.Base.IsChanged(after...)
}

// Reset set status as not changed
func (st *Station) Reset() {
	st.Base.Reset()
}

// String represents status
func (st *Station) String() string {
	st.Marshal()
	ostr := ""
	if st.Own != nil {
		ostr = fmt.Sprintf(":%s", st.Own.Short())
	}
	posstr := ""
	if st.Pos() != nil {
		posstr = fmt.Sprintf(":%s", st.Pos())
	}
	return fmt.Sprintf("%s(%d):g=%d,p=%d%s%s:%s", st.Type().Short(),
		st.ID, st.PlatformID, st.GateID, posstr, ostr, st.Name)
}
