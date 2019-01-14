package entities

import (
	"fmt"
	"time"
)

// Gate represents ticket gate in Station.
// Human must pass Gate to enter/leave Platform.
type Gate struct {
	Base
	Owner

	// Num represents how many Human can pass at the same time
	Num uint `gorm:"not null" json:"num"`
	// Mobility represents time one Human pass Gate.
	Mobility float64 `gorm:"not null" json:"mobility"`
	// Occupied represents how many Gate are used by Human.
	Occupied uint `gorm:"not null" json:"occupied"`

	M            *Model    `gorm:"-" json:"-"`
	InStation    *Station  `gorm:"-" json:"-"`
	WithPlatform *Platform `gorm:"-" json:"-"`
	out          map[uint]*Step
	in           map[uint]*Step

	StationID  uint `gorm:"not null" json:"stid"`
	PlatformID uint `gorm:"-"        json:"pid"`
}

// NewGate creates instance
func (m *Model) NewGate(st *Station) *Gate {
	g := &Gate{
		Base: NewBase(m.GenID(GATE)),
		Num:  Const.Gate.Num,
	}
	g.Init(m)
	g.Resolve(st.Own, st)
	g.Marshal()
	m.Add(g)

	g.GenOutSteps()
	g.GenInSteps()

	return g
}

func (g *Gate) GenOutSteps() {
	// skip G -> P
	// G -> C
	for _, c := range g.M.Companies {
		g.M.NewStep(g, c)
	}
}

func (g *Gate) GenInSteps() {
	// skip P -> G
	// R -> G
	for _, r := range g.M.Residences {
		g.M.NewStep(r, g)
	}
}

// Idx returns unique id field.
func (g *Gate) Idx() uint {
	return g.ID
}

// Type returns type of entitiy
func (g *Gate) Type() ModelType {
	return GATE
}

// Init creates map.
func (g *Gate) Init(m *Model) {
	g.M = m
	g.out = make(map[uint]*Step)
	g.in = make(map[uint]*Step)
}

// Pos returns location
func (g *Gate) Pos() *Point {
	if g.WithPlatform == nil {
		return nil
	}
	return g.WithPlatform.Pos()
}

// IsIn returns it should be view or not.
func (g *Gate) IsIn(x float64, y float64, scale float64) bool {
	return g.Pos().IsIn(x, y, scale)
}

// OutSteps returns where it can go to
func (g *Gate) OutSteps() map[uint]*Step {
	return g.out
}

// InSteps returns where it comes from
func (g *Gate) InSteps() map[uint]*Step {
	return g.in
}

// Resolve set reference
func (g *Gate) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			g.Owner = NewOwner(obj)
			obj.Resolve(g)
		case *Station:
			g.InStation = obj
			obj.Resolve(g)
		case *Platform:
			g.WithPlatform = obj
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	g.Marshal()
}

// Marshal set id from reference
func (g *Gate) Marshal() {
	if g.InStation != nil {
		g.StationID = g.InStation.ID
	}
	if g.WithPlatform != nil {
		g.PlatformID = g.WithPlatform.ID
	}
}

func (g *Gate) UnMarshal() {
	g.Resolve(
		g.M.Find(PLAYER, g.OwnerID), 
		g.M.Find(STATION, g.StationID))
}

// Permits represents Player is permitted to control
func (g *Gate) Permits(o *Player) bool {
	return g.Owner.Permits(o)
}

// CheckDelete check remain relation.
func (g *Gate) CheckDelete() error {
	return nil
}

// UnRef deletes related reference
func (g *Gate) UnRef() {

}

func (g *Gate) Delete() {
	for _, s := range g.out {
		g.M.Delete(s)
	}
	for _, s := range g.in {
		g.M.Delete(s)
	}
	g.M.Delete(g)
}

func (g *Gate) IsNew() bool {
	return g.Base.IsNew()
}

// IsChanged returns true when it is changed after Backup()
func (g *Gate) IsChanged(after ...time.Time) bool {
	return g.Base.IsChanged(after...)
}

// Reset set status as not changed
func (g *Gate) Reset() {
	g.Base.Reset()
}

// String represents status
func (g *Gate) String() string {
	g.Marshal()
	ostr := ""
	if g.Own != nil {
		ostr = fmt.Sprintf(":%s", g.Own.Short())
	}
	ststr := ""
	if g.InStation != nil {
		ststr = fmt.Sprintf(":%s", g.InStation.Name)
	}
	posstr := ""
	if g.Pos() != nil {
		posstr = fmt.Sprintf(":%s", g.Pos())
	}
	return fmt.Sprintf("%s(%d):st=%d,p=%d,i=%d,o=%d%s%s%s", g.Type().Short(),
		g.ID, g.StationID, g.PlatformID, len(g.in), len(g.out), posstr, ostr, ststr)
}
